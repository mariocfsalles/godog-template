package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"orders-tests/helpers"

	"github.com/segmentio/kafka-go"
)

type Consumed struct {
	Evt map[string]any
	Raw []byte
	Msg kafka.Message
}

type KafkaCtx struct {
	Brokers  []string
	Topic    string
	Reader   *kafka.Reader
	Cancel   context.CancelFunc
	Events   chan Consumed
	GroupID  string
	ClientID string
	Print    bool
	Filter   string
	StartAt  string
}

func formatHeaders(hdrs []kafka.Header) string {
	if len(hdrs) == 0 {
		return "∅"
	}
	var b strings.Builder
	for i, h := range hdrs {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(h.Key)
		b.WriteString("=")
		if len(h.Value) == 0 {
			b.WriteString("∅")
		} else {
			b.WriteString(string(h.Value))
		}
	}
	return b.String()
}

func isKafkaVerbose() bool {
	return strings.EqualFold(os.Getenv("KAFKA_VERBOSE"), "true")
}

func (k *KafkaCtx) clientIDOrDefault() string {
	if k.ClientID != "" {
		return k.ClientID
	}
	return "bdd-tests"
}

func (k *KafkaCtx) ResolveOffset() int64 {
	if k.StartAt == "beginning" {
		return kafka.FirstOffset
	}
	return kafka.LastOffset
}

func (k *KafkaCtx) ShouldPrint(msg kafka.Message, raw []byte) bool {
	if !(k.Print || isKafkaVerbose()) {
		return false
	}
	if k.Filter == "" {
		return true
	}
	// check filter in payload and headers
	rawStr := string(raw)
	if strings.Contains(rawStr, k.Filter) {
		return true
	}
	for _, h := range msg.Headers {
		if strings.Contains(strings.ToLower(h.Key), strings.ToLower(k.Filter)) ||
			strings.Contains(strings.ToLower(string(h.Value)), strings.ToLower(k.Filter)) {
			return true
		}
	}
	return false
}

func (k *KafkaCtx) PrintMessage(m kafka.Message, raw []byte) {
	group := k.GroupID
	if group == "" {
		group = "∅"
	}
	client := k.ClientID
	if client == "" {
		client = "∅"
	}

	fmt.Printf(
		"\n◆ Kafka consumer\n"+
			"  ├─ group     → %s\n"+
			"  ├─ client    → %s\n"+
			"  └─ topic     → %s\n"+
			"◆ Kafka message\n"+
			"  ├─ partition  → %d\n"+
			"  ├─ offset     → %d\n"+
			"  ├─ time       → %s\n"+
			"  ├─ key        → %s\n"+
			"  ├─ headers    → %s\n"+
			"  └─ value      →\n%s\n",
		group,
		client,
		m.Topic,
		m.Partition,
		m.Offset,
		m.Time.Format(time.RFC3339),
		func() string {
			if m.Key == nil {
				return "∅"
			}
			return string(m.Key)
		}(),
		formatHeaders(m.Headers),
		helpers.PrettyJSON(raw),
	)
}

func (k *KafkaCtx) EnsureTopic(topic string, partitions int, replication int) error {
	// connect to controller and try to create
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  k.clientIDOrDefault() + "-admin",
	}
	// get controller via metadata
	conn, err := dialer.DialContext(context.Background(), "tcp", k.Brokers[0])
	if err != nil {
		return fmt.Errorf("dial broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("controller: %w", err)
	}

	ctrlAddr := net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port))
	ctrlConn, err := dialer.DialContext(context.Background(), "tcp", ctrlAddr)
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer ctrlConn.Close()

	// CreateTopics is idempotent: if exists, ignore
	err = ctrlConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	})
	if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
		// some versions return generic error; do a small wait + metadata to confirm
		time.Sleep(500 * time.Millisecond)
	}
	// wait for topic to appear in metadata
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		tps, _ := conn.ReadPartitions()
		for _, p := range tps {
			if p.Topic == topic {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("topic %q not visible after create", topic)
}

func (k *KafkaCtx) StartFrom(topic, at string) error {
	if at != "beginning" && at != "end" {
		return fmt.Errorf("invalid start mode %q (use \"beginning\" or \"end\")", at)
	}
	k.StartAt = at
	return k.Start(topic)
}

func (k *KafkaCtx) Start(topic string) error {
	k.Topic = topic
	if err := k.EnsureTopic(topic, 1, 1); err != nil {
		return err
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  k.clientIDOrDefault(),
	}
	k.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.Brokers,
		Topic:          topic,
		GroupID:        k.GroupID,
		StartOffset:    k.ResolveOffset(),
		GroupBalancers: []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}},
		MinBytes:       1,
		MaxBytes:       10 << 20,
		MaxWait:        500 * time.Millisecond,
		Dialer:         dialer,
	})

	ctx, cancel := context.WithCancel(context.Background())
	k.Cancel = cancel

	go func() {
		defer close(k.Events)
		backoff := 200 * time.Millisecond
		for {
			m, err := k.Reader.ReadMessage(ctx)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				time.Sleep(backoff)
				if backoff < 2*time.Second {
					backoff *= 2
				}
				continue
			}
			backoff = 200 * time.Millisecond

			var obj map[string]any
			_ = json.Unmarshal(m.Value, &obj)

			if k.ShouldPrint(m, m.Value) {
				k.PrintMessage(m, m.Value)
			}

			select {
			case k.Events <- Consumed{Evt: obj, Raw: m.Value, Msg: m}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (k *KafkaCtx) Stop() {
	if k.Cancel != nil {
		k.Cancel()
	}
	if k.Reader != nil {
		_ = k.Reader.Close()
	}
}

func (k *KafkaCtx) Drain() {
	for {
		select {
		case <-k.Events:
		default:
			return
		}
	}
}

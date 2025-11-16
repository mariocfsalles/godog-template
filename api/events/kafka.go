package events

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic, clientID string) *Publisher {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},    // ordenação por chave
		RequiredAcks: kafka.RequireAll, // acks=-1
		Async:        false,
		Transport:    &kafka.Transport{ClientID: clientID},
	}
	return &Publisher{writer: w}
}

func (p *Publisher) Close() error {
	if p == nil || p.writer == nil {
		return nil
	}
	return p.writer.Close()
}

// PublishWithDigest publica o evento no Kafka e adiciona o header x-sha256.
func (p *Publisher) PublishWithDigest(
	ctx context.Context,
	key string,
	evt any,
	headers map[string]string,
) (string, error) {
	b, _ := json.Marshal(evt)
	sum := sha256.Sum256(b)
	digest := hex.EncodeToString(sum[:])

	var hs []kafka.Header
	for k, v := range headers {
		hs = append(hs, kafka.Header{Key: k, Value: []byte(v)})
	}
	hs = append(hs, kafka.Header{Key: "x-sha256", Value: []byte(digest)})

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:     []byte(key),
		Value:   b,
		Time:    time.Now(),
		Headers: hs,
	})
	return digest, err
}

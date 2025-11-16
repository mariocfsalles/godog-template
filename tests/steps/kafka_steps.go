package steps

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"orders-tests/helpers"
	"time"
)

// Consumer control
func (t *TestData) stepStartTopic(topic string) error {
	return t.kafka.Start(topic)
}

func (t *TestData) stepStartTopicFrom(topic, where string) error {
	return t.kafka.StartFrom(topic, where)
}

func (t *TestData) stepKafkaDrain() error {
	t.kafka.Drain()
	return nil
}

func (t *TestData) stepKafkaPrintOn() error {
	t.kafka.Print = true
	t.kafka.Filter = ""
	return nil
}

func (t *TestData) stepKafkaPrintOnFilter(filter string) error {
	t.kafka.Print = true
	t.kafka.Filter = filter
	return nil
}

func (t *TestData) stepKafkaPrintOff() error {
	t.kafka.Print = false
	t.kafka.Filter = ""
	return nil
}

// Event expectation
func (t *TestData) stepExpectEvent(topic, eventType, varName string, seconds int) error {
	// Garante que o consumer está no tópico certo
	if topic != t.kafka.Topic {
		t.kafka.Stop()
		t.kafka.Topic = ""
		if err := t.kafka.Start(topic); err != nil {
			return err
		}
	}

	// valor que veio da resposta HTTP (capturado em stepCaptureID)
	wantID := t.api.Vars[varName]
	timeout := time.After(time.Duration(seconds) * time.Second)

	for {
		select {
		case c, ok := <-t.kafka.Events:
			if !ok {
				return fmt.Errorf("event stream closed")
			}
			evt := c.Evt

			// Filtra pelo tipo
			if et, ok := evt["type"].(string); !ok || et != eventType {
				continue
			}

			// Compara id (antes era orderId; agora o evento tem só "id")
			if !helpers.MatchID(evt["id"], wantID) {
				continue
			}

			// Valida digest com header HTTP opcional
			if hdr := t.api.LastHdr.Get("X-Event-Sha256"); hdr != "" {
				sum := sha256.Sum256(c.Raw)
				got := hex.EncodeToString(sum[:])
				if got != hdr {
					return fmt.Errorf("payload digest mismatch: got=%s want=%s", got, hdr)
				}
			}
			return nil

		case <-timeout:
			return fmt.Errorf("event %q for id=%s not received in %ds",
				eventType, wantID, seconds)
		}
	}
}

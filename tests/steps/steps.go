package steps

import (
	"context"
	"fmt"
	"net/http"
	"orders-tests/domain"
	"orders-tests/types"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

type TestData struct {
	api           *domain.ApiCtx
	kafka         *domain.KafkaCtx
	lastOrderReq  types.OrderRequest
	lastOrderResp types.OrderResponse
}

func newAPI() *domain.ApiCtx {
	base := os.Getenv("API_BASE")
	if base == "" {
		base = "http://localhost:3000"
	}
	base = strings.TrimRight(base, "/")
	debug := strings.EqualFold(os.Getenv("HTTP_DEBUG"), "true")
	return &domain.ApiCtx{
		BaseURL: base,
		Vars:    map[string]string{},
		Debug:   debug,
		ReqHdr:  http.Header{},
	}
}

func newKafkaCtx() *domain.KafkaCtx {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9094"
	}
	start := os.Getenv("KAFKA_START")
	if start != "beginning" && start != "end" {
		start = "end" // padrão seguro
	}
	return &domain.KafkaCtx{
		Brokers: strings.Split(brokers, ","),
		Events:  make(chan domain.Consumed, 256),
		GroupID: fmt.Sprintf("bdd-%d", time.Now().UnixNano()),
		StartAt: start,
	}
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {}

func InitializeScenario(s *godog.ScenarioContext) {
	t := &TestData{
		api:   newAPI(),
		kafka: newKafkaCtx(),
	}

	s.Step(`^I send POST ([^ ]+) with JSON:$`, t.stepPostJSON)
	s.Step(`^I send PUT ([^ ]+) with JSON:$`, t.stepPutJSON)
	s.Step(`^I send GET ([^ ]+)$`, t.stepGet)
	s.Step(`^I set headers:$`, t.stepSetHeaders)
	s.Step(`^the HTTP status should be (\d+)$`, t.stepAssertStatus)
	s.Step(`^the response body should be:$`, t.stepResponseBodyShouldBe)
	s.Step(`^I store the "([^"]+)" from the response body into "([^"]+)"$`, t.stepCaptureID)
	s.Step(`^I have an order created via API:$`, t.stepHaveOrderViaAPI)

	s.Step(`^the topic "([^"]+)" is accessible$`, t.stepStartTopic)
	s.Step(`^the topic "([^"]+)" is accessible from the (beginning|end)$`, t.stepStartTopicFrom)
	s.Step(`^there must be an event on topic "([^"]+)" of type "([^"]+)" for "([^"]+)" within (\d+)s$`, t.stepExpectEvent)
	s.Step(`^I start printing Kafka events$`, t.stepKafkaPrintOn)
	s.Step(`^I start printing Kafka events matching "([^"]+)"$`, t.stepKafkaPrintOnFilter)
	s.Step(`^I stop printing Kafka events$`, t.stepKafkaPrintOff)
	s.Step(`^I clear any pending Kafka events$`, t.stepKafkaDrain)

	s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		t.kafka.Stop()
		return ctx, nil
	})
}

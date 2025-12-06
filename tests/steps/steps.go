package steps

import (
	"net/http"
	envconfig "orchestrator"
	"orchestrator/domain"
	"orchestrator/types"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

type TestData struct {
	api                 *domain.ApiCtx
	lastStoreSearchReq  types.StoreSearchRequest
	lastStoreSearchResp types.StoreSearchResponse
	lastStoreResp       types.StoreResponse
	lastProductsResp    types.ProductResponse
	lastLabelsResp      types.LabelResponse
}

func init() {
	envconfig.Load()
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

func InitializeTestSuite(sc *godog.TestSuiteContext) {}

func InitializeScenario(s *godog.ScenarioContext) {
	t := &TestData{
		api: newAPI(),
	}

	s.Step(`^I send POST ([^ ]+) with JSON:$`, t.stepPostJSON)
	s.Step(`^I send GET ([^ ]+)$`, t.stepGet)
	s.Step(`^I set headers:$`, t.stepSetHeaders)
	s.Step(`^I set query params:$`, t.stepISetQueryParams)
	s.Step(`^I set vars:$`, t.stepSetVars) // ⟵ novo
	s.Step(`^the HTTP status should be (\d+)$`, t.stepAssertStatus)
	s.Step(`^the response body should be:$`, t.stepResponseBodyShouldBe)

	s.Step(`^I store the "([^"]+)" from the response body into "([^"]+)"$`, t.stepCaptureID)
	s.Step(`^the "([^"]+)" response should match the snapshot$`, t.stepResponseShouldMatchSnapshot)
}

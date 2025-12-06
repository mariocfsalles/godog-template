package steps

import (
	"encoding/json"
	"fmt"
	"net/url"
	"orchestrator/helpers"
	"orchestrator/snapshot"
	"orchestrator/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func (t *TestData) stepPostJSON(path string, body *godog.DocString) error {

	var req types.StoreSearchRequest
	if err := json.Unmarshal([]byte(body.Content), &req); err != nil {
		return fmt.Errorf("invalid JSON for Request: %w", err)
	}
	t.lastStoreSearchReq = req

	var resp types.StoreSearchResponse
	if err := t.api.Post(path, req, &resp); err != nil {
		return err
	}
	t.lastStoreSearchResp = resp
	return nil
}

func (t *TestData) stepGet(path string) error {
	// 1) endpoint de product labelling - products
	if strings.HasPrefix(path, "/stores/") &&
		strings.Contains(path, "/productLabelling/products") {

		var resp types.ProductResponse
		if err := t.api.Get(path, &resp); err != nil {
			return err
		}
		t.lastProductsResp = resp
		return nil
	}

	// 2) endpoint de product labelling - labels
	if strings.HasPrefix(path, "/stores/") &&
		strings.Contains(path, "/productLabelling/labels") {

		var resp types.LabelResponse
		if err := t.api.Get(path, &resp); err != nil {
			return err
		}
		t.lastLabelsResp = resp
		return nil
	}

	// 3) demais GET /stores/:store_id → StoreResponse
	if strings.HasPrefix(path, "/stores/") {
		var resp types.StoreResponse
		if err := t.api.Get(path, &resp); err != nil {
			return err
		}
		t.lastStoreResp = resp
		return nil
	}

	return t.api.Get(path, nil)
}


func (t *TestData) stepSetHeaders(table *godog.Table) error {
	for i, row := range table.Rows {
		if len(row.Cells) < 2 {
			return fmt.Errorf("row %d must have 2 columns (Header, Value)", i)
		}
		key := row.Cells[0].Value
		raw := row.Cells[1].Value
		val := raw

		for k, v := range t.api.Vars {
			val = strings.ReplaceAll(val, "$"+k, v)
		}

		if strings.HasPrefix(val, "$") {
			envName := strings.TrimPrefix(val, "$")
			if envVal := os.Getenv(envName); envVal != "" {
				val = envVal
			}
		}

		t.api.ReqHdr.Set(key, val)
	}
	return nil
}

func (t *TestData) stepISetQueryParams(table *godog.Table) error {
	if table == nil || len(table.Rows) == 0 {
		return fmt.Errorf("query params table is empty")
	}

	if t.api.Query == nil {
		t.api.Query = url.Values{}
	}

	for _, row := range table.Rows {
		if len(row.Cells) != 2 {
			return fmt.Errorf("expected 2 columns per row (key | value), got %d", len(row.Cells))
		}

		key := strings.TrimSpace(row.Cells[0].Value)
		val := strings.TrimSpace(row.Cells[1].Value)

		if key == "" {
			continue
		}

		t.api.Query.Set(key, val)
	}

	return nil
}

func (t *TestData) stepSetVars(table *godog.Table) error {
	if t.api.Vars == nil {
		t.api.Vars = map[string]string{}
	}

	for i, row := range table.Rows {
		if len(row.Cells) < 2 {
			return fmt.Errorf("row %d must have 2 columns (Key, Value)", i)
		}
		key := row.Cells[0].Value
		val := row.Cells[1].Value

		if strings.HasPrefix(val, "$") {
			envName := strings.TrimPrefix(val, "$")
			if envVal := os.Getenv(envName); envVal != "" {
				val = envVal
			}
		}

		t.api.Vars[key] = val
	}

	return nil
}

func (t *TestData) stepAssertStatus(code int) error {
	if t.api.LastResp == nil {
		return fmt.Errorf("nenhuma resposta HTTP recebida")
	}
	if t.api.LastResp.StatusCode != code {
		return fmt.Errorf("esperado %d, obtido %d. body: %s",
			code, t.api.LastResp.StatusCode, string(t.api.LastBody))
	}
	return nil
}

func (t *TestData) stepResponseShouldMatchSnapshot(snapshotName string) error {
	cfg, ok := snapshot.Registry[snapshotName]
	if !ok {
		return fmt.Errorf("no snapshot configuration registered for %q", snapshotName)
	}

	fixturePath := filepath.Join("..", "fixtures", cfg.File)

	if err := cfg.CheckFunc(fixturePath, t.api.LastBody); err != nil {
		return fmt.Errorf("snapshot mismatch for %q: %w", snapshotName, err)
	}

	return nil
}

func (t *TestData) stepResponseBodyShouldBe(expectedDoc *godog.DocString) error {
	if t.api.LastBody == nil {
		return fmt.Errorf("no HTTP response body available")
	}

	var expected any
	if err := json.Unmarshal([]byte(expectedDoc.Content), &expected); err != nil {
		return fmt.Errorf("invalid expected JSON: %w", err)
	}

	var actual any
	if err := json.Unmarshal(t.api.LastBody, &actual); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if err := helpers.MatchWithPlaceholders(expected, actual, ""); err != nil {
		return fmt.Errorf(
			"%v\n\nexpected:\n%s\n\ngot:\n%s\n",
			err,
			helpers.PrettyJSON([]byte(expectedDoc.Content)),
			helpers.PrettyJSON(t.api.LastBody),
		)
	}

	return nil
}

func (t *TestData) stepCaptureID(fieldExpr, varName string) error {
	if t.api.Vars == nil {
		t.api.Vars = map[string]string{}
	}

	source := "body"
	field := fieldExpr

	if i := strings.Index(fieldExpr, "."); i != -1 {
		source = fieldExpr[:i]
		field = fieldExpr[i+1:]
	}

	switch source {
	case "product":
		if len(t.lastProductsResp.Values) == 0 {
			return fmt.Errorf("no product values captured to read %q", fieldExpr)
		}
		v := t.lastProductsResp.Values[0]

		switch field {
		case "id":
			if v.ID == "" {
				return fmt.Errorf("product.id is empty")
			}
			t.api.Vars[varName] = v.ID
			return nil
		case "itemId":
			if v.ItemID == "" {
				return fmt.Errorf("product.itemId is empty")
			}
			t.api.Vars[varName] = v.ItemID
			return nil
		default:
			return fmt.Errorf("unsupported product field %q in %q", field, fieldExpr)
		}

	case "label":
		if len(t.lastLabelsResp.Values) == 0 {
			return fmt.Errorf("no label values captured to read %q", fieldExpr)
		}
		v := t.lastLabelsResp.Values[0]

		switch field {
		case "labelId":
			if v.LabelID == "" {
				return fmt.Errorf("label.labelId is empty")
			}
			t.api.Vars[varName] = v.LabelID
			return nil
		default:
			return fmt.Errorf("unsupported label field %q in %q", field, fieldExpr)
		}

	case "body", "":
		var body map[string]any
		if err := json.Unmarshal(t.api.LastBody, &body); err != nil {
			return fmt.Errorf("unable to unmarshal LastBody: %w", err)
		}
		raw, ok := body[field]
		if !ok {
			return fmt.Errorf("field %q not found in body", field)
		}
		s, ok := raw.(string)
		if !ok {
			return fmt.Errorf("field %q is not a string (got %T)", field, raw)
		}
		t.api.Vars[varName] = s
		return nil

	default:
		return fmt.Errorf("unknown capture source %q in %q", source, fieldExpr)
	}
}

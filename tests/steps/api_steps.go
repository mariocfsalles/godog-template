package steps

import (
	"encoding/json"
	"fmt"
	"orders-tests/helpers"
	"orders-tests/types"
	"strings"

	"github.com/cucumber/godog"
)

func (t *TestData) stepPostJSON(path string, body *godog.DocString) error {
	// parse do DocString para a struct de request
	var req types.OrderRequest
	if err := json.Unmarshal([]byte(body.Content), &req); err != nil {
		return fmt.Errorf("invalid JSON for Request: %w", err)
	}
	t.lastOrderReq = req

	var resp types.OrderResponse
	if err := t.api.Post(path, req, &resp); err != nil {
		return err
	}
	t.lastOrderResp = resp

	// guarda o id para outros steps
	if t.api.Vars == nil {
		t.api.Vars = map[string]string{}
	}
	t.api.Vars["order_id"] = resp.ID

	return nil
}

func (t *TestData) stepPutJSON(path string, body *godog.DocString) error {
	if strings.Contains(path, "/orders/") && strings.HasSuffix(path, "/status") {
		var req types.OrderRequest
		if err := json.Unmarshal([]byte(body.Content), &req); err != nil {
			return fmt.Errorf("invalid JSON for StatusRequest: %w", err)
		}

		var resp types.OrderResponse
		if err := t.api.Put(path, req, &resp); err != nil {
			return err
		}
		return nil
	}
	return t.api.Put(path, body.Content, nil)
}

func (t *TestData) stepGet(path string) error {
	// Se for GET /orders/{id}, usamos a struct de Response
	if strings.HasPrefix(path, "/orders/") {
		var resp types.OrderResponse
		if err := t.api.Get(path, &resp); err != nil {
			return err
		}
		t.lastOrderResp = resp
		return nil
	}

	// Demais GETs: só chama e mantém o corpo bruto em LastBody
	return t.api.Get(path, nil)
}

func (t *TestData) stepSetHeaders(table *godog.Table) error {
	for i, row := range table.Rows {
		// pula cabeçalho se quiser, mas aqui assumo que não tem
		if len(row.Cells) < 2 {
			return fmt.Errorf("row %d must have 2 columns (Header, Value)", i)
		}
		key := row.Cells[0].Value
		val := row.Cells[1].Value

		// Expande variáveis do contexto, se usar $var
		for k, v := range t.api.Vars {
			val = strings.ReplaceAll(val, "$"+k, v)
		}

		t.api.ReqHdr.Set(key, val)
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

func (t *TestData) stepCaptureID(field, varName string) error {
	if t.api.Vars == nil {
		t.api.Vars = map[string]string{}
	}

	if field == "id" && t.lastOrderResp.ID != "" {
		t.api.Vars[varName] = t.lastOrderResp.ID
		return nil
	}

	var body map[string]any
	if err := json.Unmarshal(t.api.LastBody, &body); err != nil {
		return err
	}
	raw, ok := body[field]
	if !ok {
		return fmt.Errorf("field %q not found", field)
	}
	s, ok := raw.(string)
	if !ok {
		return fmt.Errorf("field %q is not a string (got %T)", field, raw)
	}
	t.api.Vars[varName] = s
	return nil
}

func (t *TestData) stepHaveOrderViaAPI(doc *godog.DocString) error {
	t.api.ReqHdr.Set("Content-Type", "application/json")
	if err := t.stepPostJSON("/orders", doc); err != nil {
		return err
	}
	if err := t.stepAssertStatus(201); err != nil {
		return err
	}
	return t.stepCaptureID("id", "order_id")
}

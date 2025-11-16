package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orders-tests/helpers"
	"strings"
	"time"
)

type ApiCtx struct {
	BaseURL  string
	LastHdr  http.Header
	LastResp *http.Response

	ReqHdr http.Header

	Debug    bool
	LastBody []byte
	Vars     map[string]string
}

func (a *ApiCtx) ResolvePath(p string) string {
	for k, v := range a.Vars {
		p = strings.ReplaceAll(p, "{"+k+"}", v)
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func (a *ApiCtx) LogReq(method, url string, body []byte, hdr http.Header) {
	if !a.Debug {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	fmt.Printf("\n━━━ %s ✦ HTTP request\n", ts)
	fmt.Printf("◆ %s %s\n", method, url)
	for k, vals := range hdr {
		for _, v := range vals {
			fmt.Printf("↦ %s: %s\n", k, v)
		}
	}
	if len(body) > 0 {
		fmt.Printf("→ body:\n%s\n", helpers.PrettyJSON(body))
	}
}

func (a *ApiCtx) LogResp(resp *http.Response, body []byte) {
	if !a.Debug {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	fmt.Printf("━━━ %s ✦ HTTP response\n", ts)
	fmt.Printf("◆ %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	for k, vals := range resp.Header {
		for _, v := range vals {
			fmt.Printf("↤ %s: %s\n", k, v)
		}
	}
	if len(body) > 0 {
		fmt.Printf("← body:\n%s\n\n", helpers.PrettyJSON(body))
	}
}

func (a *ApiCtx) Post(path string, reqBody any, respDest any) error {
	client := &http.Client{Timeout: 15 * time.Second}
	path = a.ResolvePath(path)
	url := a.BaseURL + path

	// 1) prepara o corpo
	var bodyBytes []byte
	switch v := reqBody.(type) {
	case nil:
		bodyBytes = nil
	case []byte:
		bodyBytes = v
	case string:
		bodyBytes = []byte(v)
	default:
		// qualquer struct/map vira JSON
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyBytes = b

		// se ninguém setou Content-Type, assume JSON
		if a.ReqHdr.Get("Content-Type") == "" {
			a.ReqHdr.Set("Content-Type", "application/json")
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	// 2) aplica headers configurados via step
	for k, vals := range a.ReqHdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	a.LogReq(http.MethodPost, url, bodyBytes, req.Header)

	// 3) faz a chamada
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.LastResp = resp
	a.LastHdr = resp.Header.Clone()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.LastBody = b

	a.LogResp(resp, b)

	// 4) desserializa na struct de resposta, se pedirem
	if respDest != nil && len(a.LastBody) > 0 {
		if err := json.Unmarshal(a.LastBody, respDest); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (a *ApiCtx) Put(path string, reqBody any, respDest any) error {
	client := &http.Client{Timeout: 15 * time.Second}
	path = a.ResolvePath(path)
	url := a.BaseURL + path

	// 1) prepara o corpo
	var bodyBytes []byte
	switch v := reqBody.(type) {
	case nil:
		bodyBytes = nil
	case []byte:
		bodyBytes = v
	case string:
		bodyBytes = []byte(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyBytes = b
		if a.ReqHdr.Get("Content-Type") == "" {
			a.ReqHdr.Set("Content-Type", "application/json")
		}
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	// 2) aplica headers configurados via step
	for k, vals := range a.ReqHdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	a.LogReq(http.MethodPut, url, bodyBytes, req.Header)

	// 3) faz a chamada
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.LastResp = resp
	a.LastHdr = resp.Header.Clone()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.LastBody = b

	a.LogResp(resp, b)

	// 4) desserializa na struct de resposta, se pedirem
	if respDest != nil && len(a.LastBody) > 0 {
		if err := json.Unmarshal(a.LastBody, respDest); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (a *ApiCtx) Get(path string, respDest any) error {
	client := &http.Client{Timeout: 15 * time.Second}
	path = a.ResolvePath(path)
	url := a.BaseURL + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for k, vals := range a.ReqHdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	a.LogReq(http.MethodGet, url, nil, req.Header)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.LastResp = resp
	a.LastHdr = resp.Header.Clone()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.LastBody = body

	a.LogResp(resp, body)

	if respDest != nil && len(a.LastBody) > 0 {
		if err := json.Unmarshal(a.LastBody, respDest); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

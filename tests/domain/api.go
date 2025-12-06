package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"orchestrator/helpers"
	"os"
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
	Query    url.Values
}

func (a *ApiCtx) ResolvePath(p string) string {
	for k, v := range a.Vars {
		p = strings.ReplaceAll(p, "{"+k+"}", v)
		p = strings.ReplaceAll(p, ":"+k, v)
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

	url, err := a.buildURL(path)
	if err != nil {
		return err
	}
	defer func() { a.Query = nil }()

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

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	for k, vals := range a.ReqHdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	a.LogReq(http.MethodPost, url, bodyBytes, req.Header)

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

	if respDest != nil && len(a.LastBody) > 0 {
		if err := json.Unmarshal(a.LastBody, respDest); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (a *ApiCtx) Get(path string, respDest any) error {
	client := &http.Client{Timeout: 15 * time.Second}

	url, err := a.buildURL(path)
	if err != nil {
		return err
	}
	defer func() { a.Query = nil }()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for k, vals := range a.ReqHdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	if respDest != nil && req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
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

func (a *ApiCtx) EffectiveBaseURL() string {
	// 1) Escape hatch explícito via vars
	if raw, ok := a.Vars["base_url"]; ok && strings.TrimSpace(raw) != "" {
		return helpers.TrimBaseURL(raw)
	}

	// 2) Header de subscrição decide qual base usar
	subKey := strings.TrimSpace(a.ReqHdr.Get("Ocp-Apim-Subscription-Key"))
	if subKey != "" {
		// VUSION
		if vusionKey := os.Getenv("VUSION_PRO_SUBSCRIPTION_KEY"); vusionKey != "" && subKey == vusionKey {
			if base := os.Getenv("API_BASE_VUSION"); base != "" {
				return helpers.TrimBaseURL(base)
			}
		}

		// VLINK
		if vlinkKey := os.Getenv("VLINK_PRO_SUBSCRIPTION_KEY"); vlinkKey != "" && subKey == vlinkKey {
			if base := os.Getenv("API_BASE_VLINK"); base != "" {
				return helpers.TrimBaseURL(base)
			}
		}
	}

	// 3) Fallback: a BaseURL padrão
	return a.BaseURL
}

func (a *ApiCtx) buildURL(path string) (string, error) {
	var fullURL string

	// Se já vier absoluta, não mexe na base
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		fullURL = path
	} else {
		path = a.ResolvePath(path)
		fullURL = a.EffectiveBaseURL() + path
	}

	// Anexa query params, se houver
	if a.Query != nil && len(a.Query) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return "", fmt.Errorf("parse URL: %w", err)
		}

		q := u.Query()
		for k, vals := range a.Query {
			for _, v := range vals {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	return fullURL, nil
}

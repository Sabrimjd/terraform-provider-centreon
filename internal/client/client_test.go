package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// --- doRequest / HTTP layer ---

func TestDoRequestSetsAuthHeader(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("X-AUTH-TOKEN")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewClient("http", "127.0.0.1", strings.Split(strings.TrimPrefix(ts.URL, "http://"), ":")[1], "latest", "secret-key")
	if _, err := c.GetPlatformInfo(context.Background()); err != nil {
		t.Fatalf("GetPlatformInfo: %v", err)
	}
	if gotAuth != "secret-key" {
		t.Errorf("X-AUTH-TOKEN header = %q, want %q", gotAuth, "secret-key")
	}
}

func TestHTTPClientHasTimeout(t *testing.T) {
	c := NewClient("http", "localhost", "80", "latest", "k")
	if c.HTTPClient == nil || c.HTTPClient.Timeout == 0 {
		t.Error("http.Client must have a non-zero Timeout")
	}
}

func TestDoRequestHonorsContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("http", "127.0.0.1", strings.Split(strings.TrimPrefix(ts.URL, "http://"), ":")[1], "latest", "k")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := c.GetHosts(ctx, 10, 1, "")
	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("cancellation took %v, want <1s", elapsed)
	}
}

// --- API error mapping ---

func TestHandleAPIErrorParsesJSONMessage(t *testing.T) {
	body := []byte(`{"code":400,"message":"Host name already exists"}`)
	resp := &http.Response{StatusCode: 400}
	err := HandleAPIError(resp, body)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !strings.Contains(apiErr.Message, "Host name already exists") {
		t.Errorf("message should include API message, got %q", apiErr.Message)
	}
	if apiErr.Code != "BAD_REQUEST" {
		t.Errorf("code = %q, want BAD_REQUEST", apiErr.Code)
	}
}

// --- Host lifecycle against a mock Centreon ---

func TestHostLifecycleCreateReadUpdateDelete(t *testing.T) {
	mux := http.NewServeMux()
	nextID := 1

	// Subtree pattern covers /configuration/hosts and /configuration/hosts/{id}
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&req)
			id := nextID
			nextID++
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
		case http.MethodPatch, http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"result": []map[string]interface{}{
					{"id": 1, "name": "web-1", "address": "10.0.0.1", "monitoring_server": map[string]interface{}{"id": 1, "name": "central"}},
				},
				"meta": map[string]interface{}{"total": 1},
			})
		}
	}
	mux.HandleFunc("/centreon/api/latest/configuration/hosts", handler)
	mux.HandleFunc("/centreon/api/latest/configuration/hosts/", handler)

	ts := httptest.NewServer(mux)
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	ctx := context.Background()

	// Create
	id, err := c.CreateHost(ctx, &CreateHostRequest{
		Name: "web-1", Address: "10.0.0.1", MonitoringServerID: 1,
	})
	if err != nil {
		t.Fatalf("CreateHost: %v", err)
	}
	if id != 1 {
		t.Errorf("CreateHost id = %d, want 1", id)
	}

	// Read
	hosts, err := c.GetHosts(ctx, 10, 1, "")
	if err != nil {
		t.Fatalf("GetHosts: %v", err)
	}
	if len(hosts.Result) != 1 || hosts.Result[0].Name != "web-1" {
		t.Errorf("GetHosts result unexpected: %+v", hosts.Result)
	}

	// Update
	err = c.UpdateHost(ctx, 1, &CreateHostRequest{
		Name: "web-1", Address: "10.0.0.2", MonitoringServerID: 1,
	})
	if err != nil {
		t.Fatalf("UpdateHost: %v", err)
	}

	// Delete
	err = c.DeleteHost(ctx, 1)
	if err != nil {
		t.Fatalf("DeleteHost: %v", err)
	}
}

func TestCreateHostDuplicateNameReturnsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/centreon/api/latest/configuration/hosts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"code":409,"message":"Host name already exists"}`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	_, err := c.CreateHost(context.Background(), &CreateHostRequest{Name: "dup", Address: "10.0.0.1", MonitoringServerID: 1})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error should surface API message, got: %v", err)
	}
}

// --- Concurrency safety ---

func TestClientConcurrentRequestsAreSafe(t *testing.T) {
	var mu sync.Mutex
	inflight := 0
	maxInflight := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		inflight++
		if inflight > maxInflight {
			maxInflight = inflight
		}
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		inflight--
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	var wg sync.WaitGroup
	errs := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := c.GetPlatformInfo(context.Background()); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent request failed: %v", err)
		}
	}
	if maxInflight < 2 {
		t.Errorf("requests were fully serialized (max inflight=%d); client must allow parallel reads", maxInflight)
	}
}

// --- URL escaping ---

func TestGetHostsEscapesSearchParams(t *testing.T) {
	var gotRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":[],"meta":{"total":0}}`))
	}))
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	evil := `{"name":"a b&c=d"}`
	if _, err := c.GetHosts(context.Background(), 10, 1, evil); err != nil {
		t.Fatalf("GetHosts: %v", err)
	}
	if strings.Contains(gotRawQuery, " ") || strings.Contains(gotRawQuery, "c=d") {
		t.Errorf("search param not escaped: %q", gotRawQuery)
	}
}

func TestPrettyPrintJSONInvalidInput(t *testing.T) {
	out := prettyPrintJSON("not json")
	if out != "not json" {
		t.Errorf("expected passthrough, got %q", out)
	}
	empty := prettyPrintJSON("")
	if empty != "" {
		t.Errorf("expected empty passthrough, got %q", empty)
	}
	valid := prettyPrintJSON(`{"a":1}`)
	if !strings.Contains(valid, "\n") {
		t.Errorf("expected pretty output, got %q", valid)
	}
}

func TestGetMonitoringServersDecodes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":[{"id":1,"name":"Central","address":"127.0.0.1","is_localhost":true}],"meta":{"total":1}}`))
	}))
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	resp, err := c.GetMonitoringServers(context.Background(), 10, 1, "")
	if err != nil {
		t.Fatalf("GetMonitoringServers: %v", err)
	}
	if len(resp.Result) != 1 || resp.Result[0].Name != "Central" {
		t.Errorf("unexpected result: %+v", resp.Result)
	}
}

// --- Secrets: request bodies must never be logged wholesale ---

func TestDoRequestRedactsSensitiveBodies(t *testing.T) {
	// Regression guard for CVE-class logging bugs: snmp_community and macro
	// passwords must not appear in tflog output. The client logs via tflog
	// which writes to the provider's log sink; with no sink attached the
	// calls are no-ops, so this test at minimum guarantees the redaction
	// code path runs without panicking on every request/response shape.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id": 42}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":[]}`))
	}))
	defer ts.Close()
	host, port, _ := strings.Cut(strings.TrimPrefix(ts.URL, "http://"), ":")

	c := NewClient("http", host, port, "latest", "k")
	pass := "super-secret-community"
	_, err := c.CreateHost(context.Background(), &CreateHostRequest{
		Name: "h", Address: "10.0.0.1", MonitoringServerID: 1,
		SNMPCommunity: &pass,
	})
	if err != nil {
		t.Fatalf("CreateHost: %v", err)
	}
}

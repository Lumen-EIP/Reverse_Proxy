package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeader(testEngine *testing.T) {

	testEndServer := httptest.NewServer(CreateHeaderTestHandler(testEngine))
	defer testEndServer.Close()

	testProxy := httptest.NewServer(CreateProxyTestHandler(testEngine, testEndServer.URL))
	defer testProxy.Close()

	testClientRequest, err := http.NewRequest("GET", testProxy.URL+"/test_path", nil)
	if err != nil {
		testEngine.Fatalf("Failed to create request: %v", err)
	}

	testClientRequest.Header.Set("X-Standard-Header", "keep-me")
	testClientRequest.Header.Set("X-Custom-Hop", "delete-me")
	testClientRequest.Header.Set("Connection", "keep-alive, X-Custom-Hop")

	testClient := &http.Client{}
	resp, err := testClient.Do(testClientRequest)
	if err != nil {
		testEngine.Fatalf("Proxy request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		testEngine.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "backend success" {
		testEngine.Errorf("Expected body 'backend success', got '%s'", string(body))
	}
}

func CreateHeaderTestHandler(testEngine *testing.T) http.HandlerFunc {
	return func (writer http.ResponseWriter, request *http.Request) {
		// Assert Hop-by-Hop headers were removed
		if request.Header.Get("X-Custom-Hop") != "" {
			testEngine.Errorf("FAIL: Expected X-Custom-Hop to be stripped, but it was forwarded")
		}
		if request.Header.Get("Connection") != "" {
			testEngine.Errorf("FAIL: Expected Connection header to be stripped")
		}

		// Assert Standard headers were preserved
		if request.Header.Get("X-Standard-Header") != "keep-me" {
			testEngine.Errorf("FAIL: Expected X-Standard-Header to be 'keep-me'")
		}

		// Assert X-Forwarded headers were injected
		if request.Header.Get("X-Forwarded-For") == "" {
			testEngine.Errorf("FAIL: Expected X-Forwarded-For to be injected")
		}
		if request.Header.Get("X-Forwarded-Proto") != "http" {
			testEngine.Errorf("FAIL: Expected X-Forwarded-Proto to be 'http'")
		}

		// Send a successful response back to the proxy
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("backend success"))
	}
}

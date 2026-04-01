package test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"reverse_proxy.deban.com/internal/header_management"
)

func CreateProxyTestHandler(testEngine *testing.T, targetServer string) http.HandlerFunc {
	backendUrl, err := url.Parse(targetServer)

	if err != nil {
		testEngine.Fatalf("Failed to parse test target URL: %v", err) 
	}

	return func (writer http.ResponseWriter, request *http.Request) {
		outReq, err := http.NewRequestWithContext(request.Context(), request.Method, backendUrl.String()+request.URL.Path, request.Body)
		if err != nil {
			http.Error(writer, "Failed to create request", http.StatusInternalServerError)
			return
		}

		header_management.CopyHeader(request.Header, outReq.Header)
		header_management.DeleteHopHeader(outReq.Header)
		header_management.SetForwardHeader(request, outReq)

		client := &http.Client{}
		resp, err := client.Do(outReq)
		if err != nil {
			http.Error(writer, "Backend server is down", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		header_management.DeleteHopHeader(resp.Header)
		header_management.CopyHeader(resp.Header, writer.Header())

		writer.WriteHeader(resp.StatusCode)
		io.Copy(writer, resp.Body)
	}
}
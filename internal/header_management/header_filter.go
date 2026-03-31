package header_management

import (
	"net/http"
	"strings"
)

func DeleteHopHeader(headerMap http.Header) {
	for _, headers := range headerMap["Connection"] {
		for _, header := range strings.Split(headers, "") {
			if strings.TrimSpace(header) != "" {
				headerMap.Del(header)
			}
		}
	}

	for _, hopHeader := range HOPE_HEADERS {
		headerMap.Del(hopHeader)
	}
}

func SetForwardHeader(inRequest, outRequest *http.Request) {
	clientIpAddress := ExtractIpAddress(inRequest)

	previousForwardHeader, isPresent := outRequest.Header["X-Forwarded-For"]
	if isPresent && len(previousForwardHeader) > 0 {
		outRequest.Header.Set("X-Forwarded-For", previousForwardHeader[0]+" "+clientIpAddress)
	} else {
		outRequest.Header.Set("X-Forwarded-For", clientIpAddress)
	}

	outRequest.Header.Set("X-Forwarded-Host", inRequest.Host)
}

func CopyHeader(srcHeader, dstHeader http.Header) {
	for key, subHeaders := range srcHeader {
		for _, header := range subHeaders {
			dstHeader.Add(key, header)
		}
	}
}
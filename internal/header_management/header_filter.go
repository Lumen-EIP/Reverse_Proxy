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

func CopyHeader(srcHeader, dstHeader http.Header) {
	for key, subHeaders := range srcHeader {
		for _, header := range subHeaders {
			dstHeader.Add(key, header)
		}
	}
}
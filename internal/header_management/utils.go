package header_management

import (
	"net"
	"net/http"
)

func ExtractIpAddress(request *http.Request) string {
	ipAddress, _, err := net.SplitHostPort(request.RemoteAddr)

	if err != nil {
		ipAddress = request.RemoteAddr
	}

	return ipAddress
}
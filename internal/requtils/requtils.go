package requtils

import (
	"net"
	"net/http"
	"strings"
)

func GetRemoteIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	if xrIP := r.Header.Get("X-Real-Ip"); xrIP != "" {
		return xrIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

package requtils

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRemoteIP(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header

		expRemoteIP string
	}{
		{
			name: "only X-Forwarded-For is present",
			headers: http.Header{
				"X-Forwarded-For": {"111.111.11.11", "333.333.33.33"},
			},
			expRemoteIP: "111.111.11.11",
		},
		{
			name: "only X-Real-IP is present",
			headers: http.Header{
				"X-Real-Ip": {"222.222.22.22"},
			},
			expRemoteIP: "222.222.22.22",
		},
		{
			name:        "neither X-Forwarded-For or X-Real-IP are present",
			headers:     http.Header{},
			expRemoteIP: "000.000.00.00",
		},
		{
			name: "both X-Forwarded-For and X-Real-IP are present",
			headers: http.Header{
				"X-Forwarded-For": {"111.111.11.11"},
				"X-Real-Ip":       {"222.222.22.22"},
			},
			expRemoteIP: "111.111.11.11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = "000.000.00.00:0000"
			maps.Copy(req.Header, tt.headers)

			fmt.Println("Headers:", req.Header)
			if ip := GetRemoteIP(req); ip != tt.expRemoteIP {
				t.Errorf("expected remote IP %s, but got %s", tt.expRemoteIP, ip)
			}
		})
	}
}

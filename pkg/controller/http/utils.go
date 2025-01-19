package http

import "net/http"

func trimToken(token string) string {
	e := min(len(token), 8)
	return token[:e] + "..."
}

func cloneHeader(src http.Header) map[string]string {
	dst := map[string]string{}

	for k, values := range src {
		for _, v := range values {
			dst[k] = v
		}
	}

	return dst
}

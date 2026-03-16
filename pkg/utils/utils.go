package utils

import (
	"encoding/base64"
	"net"
	"net/http"
	"strings"
)

// GetIP extracts the real client IP (last in X-Forwarded-For chain)
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// XorEncrypt performs XOR cipher on data with key
func XorEncrypt(data, key string) string {
	if key == "" {
		return data
	}
	keyLen := len(key)
	out := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		out[i] = data[i] ^ key[i%keyLen]
	}
	return string(out)
}

// EncodeProxyURL XOR-encrypts a URL then base64-encodes it for proxy redirection
func EncodeProxyURL(rawURL, password string) string {
	encrypted := XorEncrypt(rawURL, password)
	return base64.StdEncoding.EncodeToString([]byte(encrypted))
}

// GetBDUSS extracts BDUSS and STOKEN cookie parts from a full cookie string
func GetBDUSS(cookie string, bdclnd *string) string {
	bduss := extractCookiePart(cookie, "BDUSS")
	stoken := extractCookiePart(cookie, "STOKEN")
	result := bduss + "; " + stoken + ";"
	if bdclnd != nil && *bdclnd != "" {
		result += "BDCLND=" + *bdclnd + ";"
	}
	return result
}

func extractCookiePart(cookie, name string) string {
	for _, part := range strings.Split(cookie, ";") {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if strings.EqualFold(kv[0], name) {
			return part
		}
	}
	return ""
}

// ContainsString reports whether slice contains s
func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// StripTrailingSlash removes a trailing slash from s
func StripTrailingSlash(s string) string {
	return strings.TrimRight(s, "/")
}

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"
)

// HashIP returns a hex SHA-256 hash of the client IP. Used to deduplicate views
// without storing raw addresses (cheap GDPR mitigation).
func HashIP(remoteAddr, xff string) string {
	ip := xff
	if ip == "" {
		host, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			ip = remoteAddr
		} else {
			ip = host
		}
	} else if i := strings.Index(ip, ","); i >= 0 {
		ip = strings.TrimSpace(ip[:i])
	}
	if ip == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(sum[:])
}

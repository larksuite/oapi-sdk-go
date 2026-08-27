package ws

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func validateConnectionHeaders(header http.Header) error {
	for name := range header {
		if isReservedWebSocketHeader(name) {
			return errInvalidConnectionHeaders
		}
	}
	return nil
}

func isReservedWebSocketHeader(name string) bool {
	if strings.EqualFold(name, "Host") || strings.EqualFold(name, "Connection") || strings.EqualFold(name, "Upgrade") {
		return true
	}
	const webSocketHeaderPrefix = "Sec-WebSocket-"
	return len(name) >= len(webSocketHeaderPrefix) && strings.EqualFold(name[:len(webSocketHeaderPrefix)], webSocketHeaderPrefix)
}

func validateConnectionHost(host string) error {
	if host == "" || !isASCII(host) {
		return errInvalidConnectionHost
	}

	hostname, port, ok := splitConnectionHost(host)
	if !ok || port != "" && !isValidConnectionPort(port) {
		return errInvalidConnectionHost
	}
	if strings.HasPrefix(hostname, "[") {
		return validateBracketedIPv6(hostname)
	}
	if strings.Contains(hostname, ":") || !isValidDNSOrIPv4Host(hostname) {
		return errInvalidConnectionHost
	}
	return nil
}

func isASCII(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] >= 0x80 {
			return false
		}
	}
	return true
}

func splitConnectionHost(host string) (hostname string, port string, ok bool) {
	if strings.HasPrefix(host, "[") {
		closingBracket := strings.IndexByte(host, ']')
		if closingBracket < 0 || strings.Contains(host[closingBracket+1:], "]") {
			return "", "", false
		}
		hostname = host[:closingBracket+1]
		suffix := host[closingBracket+1:]
		if suffix == "" {
			return hostname, "", true
		}
		if !strings.HasPrefix(suffix, ":") || len(suffix) == 1 {
			return "", "", false
		}
		return hostname, suffix[1:], true
	}

	colon := strings.LastIndexByte(host, ':')
	if colon < 0 {
		return host, "", true
	}
	if strings.Contains(host[:colon], ":") || colon == 0 || colon == len(host)-1 {
		return "", "", false
	}
	return host[:colon], host[colon+1:], true
}

func isValidConnectionPort(port string) bool {
	if !isDecimal(port) || len(port) > 1 && port[0] == '0' {
		return false
	}
	number, err := strconv.Atoi(port)
	return err == nil && number >= 1 && number <= 65535
}

func validateBracketedIPv6(hostname string) error {
	if len(hostname) < 4 || hostname[len(hostname)-1] != ']' {
		return errInvalidConnectionHost
	}
	literal := hostname[1 : len(hostname)-1]
	if !strings.Contains(literal, ":") || strings.Contains(literal, "%") || net.ParseIP(literal) == nil {
		return errInvalidConnectionHost
	}
	return nil
}

func isValidDNSOrIPv4Host(hostname string) bool {
	if hostname == "" || len(hostname) > 253 || strings.HasPrefix(hostname, ".") || strings.HasSuffix(hostname, ".") {
		return false
	}
	if isNumericDottedHost(hostname) {
		return isValidIPv4Host(hostname)
	}
	for _, label := range strings.Split(hostname, ".") {
		if !isValidDNSLabel(label) {
			return false
		}
	}
	return true
}

func isNumericDottedHost(hostname string) bool {
	if !strings.Contains(hostname, ".") {
		return false
	}
	for _, label := range strings.Split(hostname, ".") {
		if !isDecimal(label) {
			return false
		}
	}
	return true
}

func isValidIPv4Host(hostname string) bool {
	parts := strings.Split(hostname, ".")
	if len(parts) != net.IPv4len {
		return false
	}
	for _, part := range parts {
		if len(part) > 1 && part[0] == '0' {
			return false
		}
		number, err := strconv.Atoi(part)
		if err != nil || number < 0 || number > 255 {
			return false
		}
	}
	return true
}

func isDecimal(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return false
		}
	}
	return true
}

func isValidDNSLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 || !isASCIIAlphaNumeric(label[0]) || !isASCIIAlphaNumeric(label[len(label)-1]) {
		return false
	}
	for i := 1; i < len(label)-1; i++ {
		if !isASCIIAlphaNumeric(label[i]) && label[i] != '-' {
			return false
		}
	}
	return true
}

func isASCIIAlphaNumeric(value byte) bool {
	return value >= '0' && value <= '9' || value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z'
}

func buildConnectionURL(endpoint string, connectionHost *string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil || strings.Contains(endpoint, "#") {
		return nil, errInvalidConnectionEndpoint
	}
	if parsed.User != nil || parsed.Hostname() == "" {
		return nil, errInvalidConnectionEndpoint
	}
	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, errInvalidConnectionEndpoint
	}
	if connectionHost != nil {
		parsed.Host = *connectionHost
	}
	return parsed, nil
}

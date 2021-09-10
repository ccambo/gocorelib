/*
Copyright 2021 MSFL Authors. All right reserved.
*/

// net - Networking 관련 기능 제공 패키지
package net

import (
	"net"
	"net/http"
	"strings"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// IsValidPort - 지정한 포트 값이 올바른 범위내에 존재하는 검증
func IsValidPort(port int) bool {
	return port > 0 && port < 65535
}

// GetRequestIP - 지정한 HTTP Request 정보에서 IP 정보 추출
// 1. Check "X-Real-Ip" header
// 2. Check "X-Forwarded-For" header
// 3. Check RemoteAddr
// 오류 발생 시는 RemoteAddr 그대로 반환
func GetRequestIP(req *http.Request) string {
	// Header - Real IP 값 검증
	if addr := strings.Trim(req.Header.Get("X-Real-Ip"), ""); addr != "" {
		return addr
	}

	// Header - Forwarded IP 값 검증
	if addr := strings.Trim(req.Header.Get("X-Forwarded-For"), ""); addr != "" {
		return addr
	}

	// RemoteAddr
	addr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}

	return addr
}

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
const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

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
	if addr := strings.Trim(req.Header.Get(XRealIP), ""); addr != "" {
		return addr
	}

	// Header - Forwarded IP 값 검증
	if addr := strings.Trim(req.Header.Get(XForwardedFor), ""); addr != "" {
		return addr
	}

	// RemoteAddr
	addr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}

	return addr
}

// RemoteIp - 지정한 HTTP Request에서 Remote IP 정보 추출
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr

	if ip := req.Header.Get("XClientIP"); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

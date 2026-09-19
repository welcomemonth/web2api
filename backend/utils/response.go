package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

const upstreamTemporaryClientMessage = "上游 Qwen 请求被网络超时、连接中断或 WAF 风控拦截；网关已按当前策略重试/切换账号但仍失败。请稍后重试，或在管理页刷新/复验账号后再试。"

func WriteError(w http.ResponseWriter, status int, detail any) {
	WriteJSON(w, status, map[string]any{"detail": sanitizeClientErrorDetail(detail)})
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func sanitizeClientErrorDetail(detail any) any {
	switch v := detail.(type) {
	case string:
		return sanitizeClientErrorString(v)
	case error:
		return sanitizeClientErrorString(v.Error())
	default:
		return detail
	}
}

func sanitizeClientErrorString(msg string) string {
	if shouldMaskUpstreamErrorMessage(msg) {
		return upstreamTemporaryClientMessage
	}
	return msg
}

func shouldMaskUpstreamErrorMessage(msg string) bool {
	lower := strings.ToLower(msg)
	if strings.TrimSpace(lower) == "" {
		return false
	}
	for _, marker := range []string{
		"chat.qwen.ai",
		"create_chat",
		"stream_chat",
		"upstream",
		"qwen api",
		"aliyun_waf",
		"<!doctype",
		"<html",
		"captcha",
		"wsarecv",
		"connection attempt failed",
		"connected party did not properly respond",
		"connected host has failed to respond",
		"net/http: request canceled",
		"context deadline exceeded",
		"i/o timeout",
	} {
		if strings.Contains(lower, marker) && IsTransientUpstreamErrorMessage(lower) {
			return true
		}
	}
	return false
}

func IsTransientUpstreamErrorMessage(lower string) bool {
	lower = strings.ToLower(lower)
	if strings.TrimSpace(lower) == "" {
		return false
	}
	if strings.Contains(lower, "http 500") ||
		strings.Contains(lower, "http 502") ||
		strings.Contains(lower, "http 503") ||
		strings.Contains(lower, "http 504") ||
		strings.Contains(lower, "status 500") ||
		strings.Contains(lower, "status 502") ||
		strings.Contains(lower, "status 503") ||
		strings.Contains(lower, "status 504") ||
		strings.Contains(lower, "status=500") ||
		strings.Contains(lower, "status=502") ||
		strings.Contains(lower, "status=503") ||
		strings.Contains(lower, "status=504") {
		return true
	}
	for _, marker := range []string{
		"create_chat parse error",
		"invalid character '<'",
		"<!doctype",
		"<html",
		"aliyun_waf",
		"waf",
		"captcha",
		"security check",
		"please enable javascript",
		"context deadline exceeded",
		"i/o timeout",
		"net/http: request canceled",
		"timeout",
		"timed out",
		"wsarecv",
		"connection attempt failed",
		"connected party did not properly respond",
		"connected host has failed to respond",
		"failed to respond",
		"connection reset",
		"connection refused",
		"connection aborted",
		"connection closed",
		"connection timed out",
		"server closed idle connection",
		"unexpected eof",
		"temporary failure",
		"temporarily unavailable",
		"bad gateway",
		"gateway timeout",
		"service unavailable",
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

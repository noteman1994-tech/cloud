package service

import "strings"

const (
	defaultBrandName     = "cloud"
	defaultBrandDocURL   = "https://www.feishu.cn/"
	defaultBrandSubtitle = "AI API Gateway Platform"
)

func normalizeBrandName(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "Sub2API" {
		return defaultBrandName
	}
	return trimmed
}

func normalizeBrandSubtitle(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "Subscription to API Conversion Platform" {
		return defaultBrandSubtitle
	}
	return trimmed
}

func normalizeBrandDocURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return defaultBrandDocURL
	}
	return trimmed
}

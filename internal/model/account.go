package model

import (
	"strings"
	"time"

	"github.com/welcomemonth/web2api/internal/utils"
)

const (
	accountUsageChat     = "chat"
	accountUsageImage    = "image"
	accountUsageVideo    = "video"
	accountUsageMetadata = "metadata"
	accountUsageUnknown  = "unknown"
)

const (
	AccountMinIntervalMS  = 1000
	RateLimitCooldown     = 600
	RateLimitBaseCooldown = 600
)

type Account struct {
	Email               string  `json:"email"`
	Password            string  `json:"password"`
	Token               string  `json:"token"`
	Cookies             string  `json:"cookies"`
	Username            string  `json:"username"`
	Source              string  `json:"source,omitempty"`
	EnvName             string  `json:"env_name,omitempty"`
	ActivationPending   bool    `json:"activation_pending"`
	StatusCode          string  `json:"status_code"`
	LastError           string  `json:"last_error"`
	LastRequestStarted  float64 `json:"last_request_started"`
	LastRequestFinished float64 `json:"last_request_finished"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	RateLimitStrikes    int     `json:"rate_limit_strikes"`

	Valid            bool                        `json:"valid,omitempty"`
	Inflight         int                         `json:"inflight,omitempty"`
	RateLimitedUntil float64                     `json:"rate_limited_until,omitempty"`
	RateLimits       map[string]AccountRateLimit `json:"rate_limits,omitempty"`
}

func (a *Account) Normalize() {
	if strings.TrimSpace(a.Source) == "" {
		a.Source = "file"
	}
	if a.StatusCode == "" {
		if a.ActivationPending {
			a.StatusCode = "pending_activation"
		} else {
			a.StatusCode = "valid"
		}
	}
	a.Valid = !a.ActivationPending && a.StatusCode != "invalid" && a.StatusCode != "auth_error" && a.StatusCode != "banned"
	a.compactRateLimits()
}

func (a *Account) AvailableFor(usage string) bool {
	if a == nil || !a.Valid || a.Token == "" {
		return false
	}
	now := float64(time.Now().UnixNano()) / 1e9
	if a.RateLimitedUntilFor(usage) > now {
		return false
	}
	minInterval := float64(utils.MaxInt(AccountMinIntervalMS, 0)) / 1000.0
	return a.LastRequestStarted+minInterval <= now
}

func (a *Account) Status() string {
	if a.ActivationPending {
		return "pending_activation"
	}
	if a.RateLimitedUntilFor(accountUsageChat) > float64(time.Now().UnixNano())/1e9 {
		return "rate_limited"
	}
	if a.Valid {
		return "valid"
	}
	if a.StatusCode != "" {
		return a.StatusCode
	}
	return "invalid"
}

func (a *Account) MigrateLegacyRateLimit() {
	if a == nil {
		return
	}
	now := float64(time.Now().UnixNano()) / 1e9
	until := a.RateLimitedUntil
	if until <= now && isRateLimitErrorMessage(a.LastError) {
		cooldown := rateLimitCooldownSeconds(a.LastError)
		until = float64(time.Now().Add(time.Duration(cooldown)*time.Second).UnixNano()) / 1e9
	}
	if until > now {
		usage := inferRateLimitUsage(a.LastError)
		if usage != accountUsageUnknown {
			a.SetRateLimitFor(usage, until, a.LastError)
		}
	}
	a.syncLegacyRateLimit()
	a.compactRateLimits()
}

func (a *Account) RateLimitedUntilFor(usage string) float64 {
	if a == nil {
		return 0
	}
	state, ok := a.RateLimits[normalizeAccountUsage(usage)]
	if !ok {
		return 0
	}
	return state.Until
}

func (a *Account) SetRateLimitFor(usage string, until float64, message string) {
	if a == nil || until <= 0 {
		return
	}
	usage = normalizeAccountUsage(usage)
	if a.RateLimits == nil {
		a.RateLimits = map[string]AccountRateLimit{}
	}
	state := a.RateLimits[usage]
	state.Until = until
	state.Reason = rateLimitReasonForUsage(usage)
	state.LastError = message
	state.Strikes++
	if state.Strikes < a.RateLimitStrikes {
		state.Strikes = a.RateLimitStrikes
	}
	a.RateLimits[usage] = state
}

func (a *Account) ClearRateLimitFor(usage string) {
	if a == nil || a.RateLimits == nil {
		return
	}
	delete(a.RateLimits, normalizeAccountUsage(usage))
	if len(a.RateLimits) == 0 {
		a.RateLimits = nil
	}
}

func (a *Account) syncLegacyRateLimit() {
	if a == nil {
		return
	}
	a.RateLimitedUntil = a.RateLimitedUntilFor(accountUsageChat)
}

func (a *Account) compactRateLimits() {
	if a == nil || len(a.RateLimits) == 0 {
		return
	}
	now := float64(time.Now().UnixNano()) / 1e9
	for usage, state := range a.RateLimits {
		normalized := normalizeAccountUsage(usage)
		if normalized != usage {
			delete(a.RateLimits, usage)
		}
		if normalized == accountUsageUnknown && state.Reason == "legacy_unknown_quota_limited" {
			delete(a.RateLimits, usage)
			continue
		}
		if state.Until <= 0 && state.LastError == "" && state.Reason == "" {
			continue
		}
		if state.Until > 0 && state.Until <= now {
			delete(a.RateLimits, usage)
			continue
		}
		state.Reason = utils.FirstNonEmpty(state.Reason, rateLimitReasonForUsage(normalized))
		a.RateLimits[normalized] = state
	}
	if len(a.RateLimits) == 0 {
		a.RateLimits = nil
	}
}

func rateLimitCooldownSeconds(lower string) int {
	lower = strings.ToLower(lower)
	if strings.Contains(lower, "today's usage") ||
		strings.Contains(lower, "todays usage") ||
		strings.Contains(lower, "daily usage") ||
		strings.Contains(lower, "daily limit") ||
		strings.Contains(lower, "upper limit for today") {
		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 10, 0, 0, now.Location())
		return utils.MaxInt(int(time.Until(nextDay).Seconds()), RateLimitBaseCooldown)
	}
	return RateLimitBaseCooldown
}

func inferRateLimitUsage(message string) string {
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "image") || strings.Contains(lower, "image_gen") || strings.Contains(lower, "t2i") || strings.Contains(lower, "picture") || strings.Contains(lower, "photo") || strings.Contains(lower, "图片") || strings.Contains(lower, "图像") || strings.Contains(lower, "cdn.qwenlm.ai"):
		return accountUsageImage
	case strings.Contains(lower, "video") || strings.Contains(lower, "t2v") || strings.Contains(lower, ".mp4") || strings.Contains(lower, "视频"):
		return accountUsageVideo
	case strings.Contains(lower, "chat") || strings.Contains(lower, "t2t") || strings.Contains(lower, "message") || strings.Contains(lower, "completion") || strings.Contains(lower, "对话"):
		return accountUsageChat
	default:
		return accountUsageUnknown
	}
}

func isRateLimitErrorMessage(lower string) bool {
	lower = strings.ToLower(lower)
	if strings.Contains(lower, "http 429") ||
		strings.Contains(lower, "status 429") ||
		strings.Contains(lower, "status=429") ||
		strings.Contains(lower, "code=429") ||
		strings.Contains(lower, "code 429") {
		return true
	}
	for _, marker := range []string{
		"ratelimited",
		"rate_limited",
		"rate limited",
		"rate limit",
		"too many requests",
		"upper limit",
		"usage limit",
		"today's usage",
		"todays usage",
		"daily usage",
		"daily limit",
		"quota",
		"free quota",
		"insufficient quota",
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func normalizeAccountUsage(usage string) string {
	switch strings.ToLower(strings.TrimSpace(usage)) {
	case "", accountUsageChat, "completion", "conversation", "message", "messages", "text", "t2t":
		return accountUsageChat
	case accountUsageImage, "images", "image_gen", "t2i", "picture", "photo":
		return accountUsageImage
	case accountUsageVideo, "videos", "t2v", "video_gen":
		return accountUsageVideo
	case accountUsageMetadata, "models", "model", "account", "verify", "verification":
		return accountUsageMetadata
	case accountUsageUnknown, "legacy", "global":
		return accountUsageUnknown
	default:
		return strings.ToLower(strings.TrimSpace(usage))
	}
}

func rateLimitReasonForUsage(usage string) string {
	switch normalizeAccountUsage(usage) {
	case accountUsageImage:
		return "image_quota_limited"
	case accountUsageVideo:
		return "video_quota_limited"
	case accountUsageMetadata:
		return "metadata_rate_limited"
	case accountUsageUnknown:
		return "legacy_unknown_quota_limited"
	default:
		return "chat_rate_limited"
	}
}

type AccountRateLimit struct {
	Until     float64 `json:"until,omitempty"`
	Reason    string  `json:"reason,omitempty"`
	LastError string  `json:"last_error,omitempty"`
	Strikes   int     `json:"strikes,omitempty"`
}

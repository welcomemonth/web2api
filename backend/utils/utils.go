package utils

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func DeepCopyMap(src map[string]any) map[string]any {
	if src == nil {
		return map[string]any{}
	}
	raw, _ := json.Marshal(src)
	var dst map[string]any
	_ = json.Unmarshal(raw, &dst)
	if dst == nil {
		dst = map[string]any{}
	}
	return dst
}

func NormalizeLower(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func RandomID() string {
	buf := make([]byte, 16)
	if _, err := cryptorand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// func writeError(w http.ResponseWriter, status int, detail any) {
// 	utils.WriteJSON(w, status, map[string]any{"detail": sanitizeClientErrorDetail(detail)})
// }

// func utils.DecodeJSON(r *http.Request, dst any) error {
// 	defer r.Body.Close()
// 	dec := json.NewDecoder(io.LimitReader(r.Body, 256<<20))
// 	dec.UseNumber()
// 	return dec.Decode(dst)
// }

func StringValue(m map[string]any, key, fallback string) string {
	if m == nil {
		return fallback
	}
	v, ok := m[key]
	if !ok {
		return fallback
	}
	return AnyString(v, fallback)
}

func AnyString(v any, fallback string) string {
	switch x := v.(type) {
	case string:
		if x != "" {
			return x
		}
	case json.Number:
		return x.String()
	case fmt.Stringer:
		return x.String()
	}
	return fallback
}

func IntValue(m map[string]any, key string, fallback int) int {
	v, ok := m[key]
	if !ok {
		return fallback
	}
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case json.Number:
		if i, err := strconv.Atoi(x.String()); err == nil {
			return i
		}
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(x)); err == nil {
			return i
		}
	}
	return fallback
}

func BoolValue(v any) bool {
	b, _ := v.(bool)
	return b
}

func CoerceBool(v any) *bool {
	switch x := v.(type) {
	case bool:
		return &x
	case float64:
		b := x != 0
		return &b
	case json.Number:
		i, _ := strconv.Atoi(x.String())
		b := i != 0
		return &b
	case string:
		switch NormalizeLower(x) {
		case "1", "true", "yes", "on", "enable", "enabled", "auto", "thinking":
			b := true
			return &b
		case "0", "false", "no", "off", "disable", "disabled", "fast", "none":
			b := false
			return &b
		}
	}
	return nil
}

func AnyList(v any) []any {
	if list, ok := v.([]any); ok {
		return list
	}
	return nil
}

func FirstStringAny(values ...any) string {
	for _, v := range values {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func ExtractModelList(decoded any) []map[string]any {
	switch v := decoded.(type) {
	case []any:
		return MapList(v)
	case map[string]any:
		if out := MapList(v["data"]); len(out) > 0 {
			return out
		}
		if out := MapList(v["models"]); len(out) > 0 {
			return out
		}
	}
	return nil
}

func MapList(v any) []map[string]any {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func PseudoEmbedding(text string) []float64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(text))
	base := float64(h.Sum64()%math.MaxUint32) / float64(math.MaxUint32)
	vec := make([]float64, 1536)
	for i := range vec {
		vec[i] = (base*float64(i%10))/10.0 - 0.5
	}
	return vec
}

func SplitExts(value string) map[string]bool {
	out := map[string]bool{}
	for _, item := range strings.Split(value, ",") {
		item = strings.Trim(strings.ToLower(item), " .")
		if item != "" {
			out[item] = true
		}
	}
	return out
}

func FileExt(name string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
}

func Truncate(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit]
}

func Trim(text string, limit int) string {
	return Truncate(text, limit)
}

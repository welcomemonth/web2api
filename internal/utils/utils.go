package utils

import (
	"crypto/md5"
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

// MaxInt 返回 a 和 b 中较大的整数值。
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FirstNonEmpty 返回 values 中第一个非空字符串。
//
// 判断字符串是否为空时，会忽略字符串首尾的空白字符；
// 但返回值保持原字符串内容不变。
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

// normalizeLower 将字符串去除首尾空白后转换为小写。
//
// 该函数主要用于字符串标准化处理和比较。
func normalizeLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// RandomID 生成一个随机 ID。
//
// 正常情况下使用 crypto/rand 生成 16 字节随机数据，
// 并将其编码为十六进制字符串返回。
// 如果加密安全随机数生成失败，则使用当前 Unix 纳秒时间戳作为备用值。
func RandomID() string {
	buf := make([]byte, 16)

	if _, err := cryptorand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(buf)
}

// StringValue 从 map 中获取指定 key 对应的字符串值。
//
// 如果 map 为 nil、key 不存在，或者对应值无法转换为字符串，
// 则返回 fallback。
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

// AnyString 尝试将任意值转换为字符串。
//
// 支持以下类型：
//   - string
//   - json.Number
//   - fmt.Stringer
//
// 如果无法转换为有效字符串，则返回 fallback。
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

// IntValue 从 map 中获取指定 key 对应的整数值。
//
// 支持以下类型：
//   - int
//   - float64
//   - json.Number
//   - string
//
// 如果 key 不存在或无法转换为整数，则返回 fallback。
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

// BoolValue 尝试将任意值转换为 bool。
//
// 如果 v 本身是 bool 类型，则返回其值；
// 如果不是 bool 类型，则返回 false。
func BoolValue(v any) bool {
	b, _ := v.(bool)
	return b
}

// CoerceBool 尝试将任意值转换为 bool。
//
// 支持以下类型及表示方式：
//   - bool
//   - float64
//   - json.Number
//   - string
//
// 字符串支持常见的布尔表示，例如：
// "1"、"true"、"yes"、"on"、"enable"、"enabled"、"auto"、"thinking"
// "0"、"false"、"no"、"off"、"disable"、"disabled"、"fast"、"none"
//
// 如果无法识别该值，则返回 nil。
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
		switch normalizeLower(x) {
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

// AnyList 尝试将任意值转换为 []any。
//
// 如果 v 的实际类型不是 []any，则返回 nil。
func AnyList(v any) []any {
	if list, ok := v.([]any); ok {
		return list
	}

	return nil
}

// FirstStringAny 返回 values 中第一个非空的字符串值。
//
// 判断字符串是否为空时，会忽略首尾空白字符；
// 返回值本身也会去除首尾空白字符。
func FirstStringAny(values ...any) string {
	for _, v := range values {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}

	return ""
}

// ExtractModelList 从解析后的数据中提取模型列表。
//
// 支持以下几种数据格式：
//   - []any：直接作为模型列表
//   - map[string]any：从 "data" 字段中提取模型列表
//   - map[string]any：从 "models" 字段中提取模型列表
//
// 如果无法找到有效的模型列表，则返回 nil。
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

// MapList 将 []any 转换为 []map[string]any。
//
// 只有实际类型为 map[string]any 的元素才会被保留，
// 其他类型的元素会被忽略。
//
// 如果 v 不是 []any 类型，则返回 nil。
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

// PseudoEmbedding 根据文本生成一个伪 Embedding 向量。
//
// 函数首先使用 FNV-1a 对文本进行哈希，
// 再根据哈希结果生成一个固定长度为 1536 的 float64 向量。
//
// 该函数并不是真正的语义 Embedding，生成的向量不具备真实的语义信息，
// 主要用于开发测试、占位或不需要真实向量模型的场景。
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

// SplitExts 将逗号分隔的文件扩展名字符串转换为集合。
//
// 函数会对每个扩展名执行以下处理：
//   - 转换为小写
//   - 去除首尾空白
//   - 去除首尾的点号
//
// 空字符串会被忽略。
//
// 例如：
//
//	"jpg, PNG, .gif"
//
// 最终会转换为：
//
//	map[string]bool{
//		"jpg": true,
//		"png": true,
//		"gif": true,
//	}
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

// FileExt 获取文件名的扩展名，并转换为小写。
//
// 返回结果不包含开头的 "."。
//
// 例如：
// "example.TXT" -> "txt"
// "test.mp4"    -> "mp4"
func FileExt(name string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
}

// Truncate 去除字符串首尾空白，并限制字符串的最大长度。
//
// 当 limit 小于等于 0，或者字符串长度未超过 limit 时，
// 直接返回处理后的完整字符串。
//
// 注意：该函数通过 len 计算字符串长度，因此限制的是字节数，
// 而不是 Unicode 字符数。如果处理中文等多字节字符，
// 可能会出现截断到半个字符的情况。
func Truncate(text string, limit int) string {
	text = strings.TrimSpace(text)

	if limit <= 0 || len(text) <= limit {
		return text
	}

	return text[:limit]
}

// Trim 是 Truncate 的别名。
//
// 用于去除字符串首尾空白，并将字符串长度限制为指定的字节数。
func Trim(text string, limit int) string {
	return Truncate(text, limit)
}

func GetEmailHashFilename(email string) string {
	// 步骤 A：标准化处理（极其重要！）
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	// 步骤 B：计算 MD5 哈希
	hasher := md5.New()
	hasher.Write([]byte(normalizedEmail))
	hashString := hex.EncodeToString(hasher.Sum(nil))

	// 步骤 C：拼接扩展名
	return hashString + ".json"
}

func VerifyEmail(email string) bool {
	if !strings.Contains(email, "@") || strings.Contains(email, "/") || strings.Contains(email, "\\") {
		return false
	}

	return true
}

func RedactToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "-"
	}
	if len(token) <= 12 {
		return "token-hidden"
	}
	return token[:6] + "..." + token[len(token)-4:]
}

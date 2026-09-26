package utils

import "encoding/json"

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

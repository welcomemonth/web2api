package model

import "testing"

func TestAPIKeyStatusValues(t *testing.T) {
	if APIKeyStatusEnabled != "enabled" || APIKeyStatusDisabled != "disabled" {
		t.Errorf("状态常量值错误")
	}
}

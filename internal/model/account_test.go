package model

import "testing"

func TestAccountStatusValues(t *testing.T) {
	if AccountStatusNormal != "normal" ||
		AccountStatusError != "error" ||
		AccountStatusDisabled != "disabled" {
		t.Errorf("状态常量值错误")
	}
}

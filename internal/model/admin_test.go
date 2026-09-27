package model

import (
	"encoding/json"
	"testing"
)

func TestAdminJSON(t *testing.T) {
	in := Admin{Username: "admin", Password: "admin123"}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(b); got != `{"username":"admin","password":"admin123"}` {
		t.Errorf("Marshal = %s", got)
	}

	var out Admin
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if out != in {
		t.Errorf("roundtrip = %+v, want %+v", out, in)
	}
}

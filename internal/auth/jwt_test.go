package auth

import (
	"errors"
	"testing"
	"time"
)

const testSecret = "test-secret"

func TestJWTGenerateAndParse(t *testing.T) {
	svc := NewJWTService(testSecret, 24*time.Hour)

	token, err := svc.Generate("admin")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %q, want admin", claims.Username)
	}
	if claims.ExpiresAt == nil || time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Errorf("Token 应在将来过期: %+v", claims.ExpiresAt)
	}
}

func TestJWTParseExpired(t *testing.T) {
	// 有效期为负 → 签发即过期
	svc := NewJWTService(testSecret, -time.Hour)

	token, err := svc.Generate("admin")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if _, err := svc.Parse(token); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("过期 Token 应返回 ErrInvalidToken, got %v", err)
	}
}

func TestJWTParseWrongSecret(t *testing.T) {
	token, err := NewJWTService("secret-a", time.Hour).Generate("admin")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewJWTService("secret-b", time.Hour).Parse(token); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("换密钥后应校验失败, got %v", err)
	}
}

func TestJWTParseTampered(t *testing.T) {
	svc := NewJWTService(testSecret, time.Hour)
	token, err := svc.Generate("admin")
	if err != nil {
		t.Fatal(err)
	}

	// 篡改签名段最后一个字符
	tampered := token[:len(token)-1] + "X"
	if tampered == token {
		tampered = token[:len(token)-1] + "Y"
	}
	if _, err := svc.Parse(tampered); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("被篡改的 Token 应校验失败, got %v", err)
	}
}

func TestJWTParseGarbage(t *testing.T) {
	svc := NewJWTService(testSecret, time.Hour)

	for _, tok := range []string{"", "not-a-token", "a.b.c"} {
		if _, err := svc.Parse(tok); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("Parse(%q) 应返回 ErrInvalidToken, got %v", tok, err)
		}
	}
}

// alg=none 是经典的 JWT 绕过攻击：攻击者把算法声明为 none 并去掉签名。
// 校验方若不加限制就会接受任意伪造的 Token。
func TestJWTParseRejectsAlgNone(t *testing.T) {
	// {"alg":"none","typ":"JWT"} . {"username":"hacker"} . <空签名>
	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0." +
		"eyJ1c2VybmFtZSI6ImhhY2tlciIsInN1YiI6ImhhY2tlciJ9."

	svc := NewJWTService(testSecret, time.Hour)
	if _, err := svc.Parse(noneToken); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("alg=none 的 Token 必须被拒绝, got %v", err)
	}
}

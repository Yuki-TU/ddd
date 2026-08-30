package user_test

import (
	"errors"
	"testing"

	domain "github.com/Yuki-TU/ddd/domain/user"
)

// 値オブジェクトのルール（1文字以上）をコンストラクタで確認する。
func TestNewUserName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr error
	}{
		{name: "1文字以上なら生成できる", in: "taro", want: "taro"},
		{name: "空文字はエラー", in: "  ", wantErr: domain.ErrInvalidUserName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewUserName(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.String() != tt.want {
				t.Errorf("got %q, want %q", got.String(), tt.want)
			}
		})
	}
}

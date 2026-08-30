package user

import (
	"strings"
	"unicode/utf8"
)

// UserName はユーザ名の値オブジェクト。
// プリミティブな string ではなく、名前のルールをここに閉じ込める。
type UserName struct {
	value string
}

// NewUserName はユーザ名を生成する。ルール違反ならエラーを返す。
func NewUserName(value string) (UserName, error) {
	name := strings.TrimSpace(value)
	// 業務ルール: ユーザ名は1文字以上
	if utf8.RuneCountInString(name) < 1 {
		return UserName{}, ErrInvalidUserName
	}
	return UserName{value: name}, nil
}

// String は永続化や表示用に内部の文字列を返す。
func (n UserName) String() string {
	return n.value
}

// Equal は等価性で比較する（値オブジェクトなので中身が同じなら同じ値）
func (n UserName) Equal(other UserName) bool {
	return n.value == other.value
}

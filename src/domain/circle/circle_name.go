package circle

import (
	"strings"
	"unicode/utf8"
)

// CircleName はサークル名の値オブジェクト。
type CircleName struct {
	value string
}

// NewCircleName はサークル名を生成する。ルール違反ならエラーを返す。
func NewCircleName(value string) (CircleName, error) {
	name := strings.TrimSpace(value)
	// 業務ルール: サークル名は1文字以上
	if utf8.RuneCountInString(name) < 1 {
		return CircleName{}, ErrInvalidCircleName
	}
	return CircleName{value: name}, nil
}

// String は永続化や表示用に内部の文字列を返す。
func (n CircleName) String() string {
	return n.value
}

// Equal は等価性で比較する。
func (n CircleName) Equal(other CircleName) bool {
	return n.value == other.value
}

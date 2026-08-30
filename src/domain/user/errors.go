package user

import "errors"

// ドメインのルール違反を表すエラー。
// HTTP のステータスに変換するのはプレゼンテーション層の仕事。
var (
	ErrInvalidUserName   = errors.New("ユーザ名は1文字以上である必要があります")
	ErrUserIDRequired    = errors.New("idが必要です")
	ErrUserNameRequired  = errors.New("名前が必要です")
	ErrUserAlreadyExists = errors.New("ユーザ名は既に存在します")
	ErrUserNotFound      = errors.New("ユーザが見つかりません")
	ErrCannotDeleteUser  = errors.New("サークルのオーナーのため退会できません")
)

package circle

import "errors"

// ドメインのルール違反を表すエラー。
var (
	ErrInvalidCircleName   = errors.New("サークル名は1文字以上である必要があります")
	ErrCircleIDRequired    = errors.New("idが必要です")
	ErrCircleNameRequired  = errors.New("名前が必要です")
	ErrOwnerRequired       = errors.New("オーナーが必要です")
	ErrCircleAlreadyExists = errors.New("サークル名は既に存在します")
	ErrCircleNotFound      = errors.New("サークルが見つかりません")
	ErrCircleFull          = errors.New("サークルのメンバー数が上限に達しています")
	ErrAlreadyMember       = errors.New("すでに参加しています")
)

package user

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// UserData は転送専用の DTO。
// ドメインオブジェクトをアプリケーションの外へ出さないための入れ物。
// 振る舞い（ChangeName など）は持たせず、値の受け渡しだけをする。
type UserData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewUserData はドメインオブジェクトから DTO を作る。
// 引数に User ごと渡すことで、項目が増えても修正はここだけになる。
func NewUserData(u *domain.User) UserData {
	return UserData{
		ID:   u.ID().String(),
		Name: u.Name().String(),
	}
}

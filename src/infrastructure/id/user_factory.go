package id

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// UserFactory は UUID でユーザ ID を採番する。
// ドメインの Factory インターフェースを実現する。
type UserFactory struct{}

// NewUserFactory はユーザファクトリを生成する。
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// Create は採番してからユーザオブジェクトを作る。
func (f *UserFactory) Create(name domain.UserName) (*domain.User, error) {
	// 保存前に ID を採番する
	id, err := newUUID()
	if err != nil {
		return nil, err
	}
	userID, err := domain.NewUserID(id)
	if err != nil {
		return nil, err
	}
	return domain.NewUser(userID, name)
}

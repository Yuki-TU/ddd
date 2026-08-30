package user

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// GetCommand はユーザ検索に必要なデータ。
type GetCommand struct {
	ID   string
	Name string
}

// GetService はユーザ情報確認ユースケース。
type GetService struct {
	// ユーザリポジトリのインターフェースに依存する
	users domain.Repository
}

// NewGetService はユーザ情報確認ユースケースを生成する。
func NewGetService(users domain.Repository) *GetService {
	return &GetService{users: users}
}

// Get はユーザ ID またはユーザ名から検索し、DTO を返す。
// User エンティティをそのまま返さない（ドメインの振る舞いを外に漏らさない）。
func (s *GetService) Get(cmd GetCommand) (UserData, error) {
	var (
		u   *domain.User
		err error
	)
	switch {
	case cmd.ID != "":
		id, idErr := domain.NewUserID(cmd.ID)
		if idErr != nil {
			return UserData{}, idErr
		}
		// ID で再構築する
		u, err = s.users.FindByID(id)
	case cmd.Name != "":
		name, nameErr := domain.NewUserName(cmd.Name)
		if nameErr != nil {
			return UserData{}, nameErr
		}
		// 名前で再構築する
		u, err = s.users.FindByName(name)
	default:
		// ID も名前も無い
		return UserData{}, domain.ErrUserNotFound
	}
	if err != nil {
		return UserData{}, err
	}
	if u == nil {
		return UserData{}, domain.ErrUserNotFound
	}
	return NewUserData(u), nil
}

package user

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// DeleteCommand は退会に必要なデータ。
type DeleteCommand struct {
	ID string
}

// DeleteService は退会ユースケース。
type DeleteService struct {
	// ユーザリポジトリのインターフェースに依存する
	users domain.Repository
}

// NewDeleteService は退会ユースケースを生成する。
func NewDeleteService(users domain.Repository) *DeleteService {
	return &DeleteService{users: users}
}

// Delete はユーザを退会させる。
func (s *DeleteService) Delete(cmd DeleteCommand) error {
	// コマンドの ID を値オブジェクトにする
	id, err := domain.NewUserID(cmd.ID)
	if err != nil {
		return err
	}
	// 削除対象を再構築する
	u, err := s.users.FindByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return domain.ErrUserNotFound
	}
	// 削除はリポジトリの責務
	return s.users.Delete(u)
}

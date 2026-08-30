package user

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// UpdateCommand はユーザ更新に必要なデータ。
// メールや住所が増えても、メソッド引数ではなくこのオブジェクトに足す。
type UpdateCommand struct {
	ID   string
	Name string
}

// UpdateService はユーザ情報更新ユースケース。
type UpdateService struct {
	// ユーザリポジトリのインターフェースに依存する
	users domain.Repository
	// 名前の重複確認
	service *domain.Service
}

// NewUpdateService はユーザ更新ユースケースを生成する。
func NewUpdateService(users domain.Repository, service *domain.Service) *UpdateService {
	return &UpdateService{users: users, service: service}
}

// Update はユーザ情報を一括更新する。
func (s *UpdateService) Update(cmd UpdateCommand) error {
	// コマンドの ID を値オブジェクトにする
	id, err := domain.NewUserID(cmd.ID)
	if err != nil {
		return err
	}
	// ID で再構築する
	u, err := s.users.FindByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return domain.ErrUserNotFound
	}

	// コマンドの名前を値オブジェクトにする
	name, err := domain.NewUserName(cmd.Name)
	if err != nil {
		return err
	}
	// 値の更新はエンティティの振る舞い（リポジトリの UpdateName にはしない）
	if err := u.ChangeName(name); err != nil {
		return err
	}

	// ドメインサービスで重複確認
	exists, err := s.service.Exists(u)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrUserAlreadyExists
	}

	// 上書き保存（永続化）
	return s.users.Save(u)
}

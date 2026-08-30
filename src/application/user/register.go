package user

import (
	domain "github.com/Yuki-TU/ddd/domain/user"
)

// RegisterCommand はユーザ登録に必要なデータ。
type RegisterCommand struct {
	Name string
}

// RegisterService はユーザ登録ユースケース。
// ドメインオブジェクトとリポジトリを組み立てて、登録という機能を実現する。
type RegisterService struct {
	// ユーザリポジトリのインターフェースに依存する
	users domain.Repository
	// 採番はファクトリに任せる
	factory domain.Factory
	// 重複確認はドメインサービス（具象でよい）
	service *domain.Service
}

// NewRegisterService はユーザ登録ユースケースを生成する。
func NewRegisterService(users domain.Repository, factory domain.Factory, service *domain.Service) *RegisterService {
	return &RegisterService{
		users:   users,
		factory: factory,
		service: service,
	}
}

// Register はユーザの重複がないかを確認し、重複がなければ登録する。
func (s *RegisterService) Register(cmd RegisterCommand) (UserData, error) {
	// プリミティブを値オブジェクトにする（ルールは UserName 側）
	name, err := domain.NewUserName(cmd.Name)
	if err != nil {
		return UserData{}, err
	}

	// ファクトリ経由でユーザ作成（ID 採番をドメインから分離）
	u, err := s.factory.Create(name)
	if err != nil {
		return UserData{}, err
	}

	// ドメインサービスで重複確認
	exists, err := s.service.Exists(u)
	if err != nil {
		return UserData{}, err
	}
	if exists {
		return UserData{}, domain.ErrUserAlreadyExists
	}

	// 保存（永続化）
	if err := s.users.Save(u); err != nil {
		return UserData{}, err
	}
	// エンティティではなく DTO を返す
	return NewUserData(u), nil
}

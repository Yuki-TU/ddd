package circle

import (
	domain "github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

// CircleData は転送専用の DTO。
type CircleData struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	OwnerID string   `json:"ownerId"`
	Members []string `json:"members"`
}

// NewCircleData はドメインオブジェクトから DTO を作る。
func NewCircleData(c *domain.Circle) CircleData {
	members := make([]string, 0, len(c.Members()))
	for _, m := range c.Members() {
		members = append(members, m.String())
	}
	return CircleData{
		ID:      c.ID().String(),
		Name:    c.Name().String(),
		OwnerID: c.OwnerID().String(),
		Members: members,
	}
}

// CreateCommand はサークル作成に必要なデータ。
type CreateCommand struct {
	UserID string
	Name   string
}

// CreateService はサークル作成ユースケース。
type CreateService struct {
	// サークルリポジトリのインターフェースに依存する
	circles domain.Repository
	// オーナー存在確認のためユーザリポジトリにも依存する
	users   user.Repository
	factory domain.Factory
	service *domain.Service
}

// NewCreateService はサークル作成ユースケースを生成する。
func NewCreateService(circles domain.Repository, users user.Repository, factory domain.Factory, service *domain.Service) *CreateService {
	return &CreateService{
		circles: circles,
		users:   users,
		factory: factory,
		service: service,
	}
}

// Create はオーナーを確認し、名前が重複していなければサークルを登録する。
func (s *CreateService) Create(cmd CreateCommand) (CircleData, error) {
	// コマンドのユーザIDを値オブジェクトにする
	userID, err := user.NewUserID(cmd.UserID)
	if err != nil {
		return CircleData{}, err
	}
	// オーナーとなるユーザを再構築する
	owner, err := s.users.FindByID(userID)
	if err != nil {
		return CircleData{}, err
	}
	if owner == nil {
		return CircleData{}, user.ErrUserNotFound
	}

	// コマンドの名前を値オブジェクトにする
	name, err := domain.NewCircleName(cmd.Name)
	if err != nil {
		return CircleData{}, err
	}
	// ファクトリ経由でサークル作成
	c, err := s.factory.Create(name, owner)
	if err != nil {
		return CircleData{}, err
	}

	// ドメインサービスで名前の重複確認
	exists, err := s.service.Exists(c)
	if err != nil {
		return CircleData{}, err
	}
	if exists {
		return CircleData{}, domain.ErrCircleAlreadyExists
	}

	// 保存（永続化）
	if err := s.circles.Save(c); err != nil {
		return CircleData{}, err
	}
	// エンティティではなく DTO を返す
	return NewCircleData(c), nil
}

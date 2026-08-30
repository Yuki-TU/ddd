package circle

import (
	domain "github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

// GetCommand はサークル取得に必要なデータ。
type GetCommand struct {
	ID string
}

// GetService はサークル情報確認ユースケース。
type GetService struct {
	circles domain.Repository
}

// NewGetService はサークル情報確認ユースケースを生成する。
func NewGetService(circles domain.Repository) *GetService {
	return &GetService{circles: circles}
}

// Get はサークルを再構築し、DTO を返す。
func (s *GetService) Get(cmd GetCommand) (CircleData, error) {
	id, err := domain.NewCircleID(cmd.ID)
	if err != nil {
		return CircleData{}, err
	}
	// ID で再構築する
	c, err := s.circles.FindByID(id)
	if err != nil {
		return CircleData{}, err
	}
	if c == nil {
		return CircleData{}, domain.ErrCircleNotFound
	}
	// エンティティではなく DTO を返す
	return NewCircleData(c), nil
}

// JoinCommand はサークル参加に必要なデータ。
type JoinCommand struct {
	CircleID string
	UserID   string
}

// JoinService はサークル参加ユースケース。
type JoinService struct {
	circles domain.Repository
	users   user.Repository
}

// NewJoinService はサークル参加ユースケースを生成する。
func NewJoinService(circles domain.Repository, users user.Repository) *JoinService {
	return &JoinService{circles: circles, users: users}
}

// Join はユーザをサークル集約に参加させ、保存する。
func (s *JoinService) Join(cmd JoinCommand) error {
	// コマンドを値オブジェクトにする
	circleID, err := domain.NewCircleID(cmd.CircleID)
	if err != nil {
		return err
	}
	userID, err := user.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	// 参加するユーザを再構築する
	member, err := s.users.FindByID(userID)
	if err != nil {
		return err
	}
	if member == nil {
		return user.ErrUserNotFound
	}

	// サークル集約を再構築する
	c, err := s.circles.FindByID(circleID)
	if err != nil {
		return err
	}
	if c == nil {
		return domain.ErrCircleNotFound
	}

	// 満員・既参加などのルールは集約ルートの Join に任せる
	if err := c.Join(member.ID()); err != nil {
		return err
	}
	// 保存（永続化）
	return s.circles.Save(c)
}

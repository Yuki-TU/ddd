package id

import (
	"github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

// CircleFactory は UUID でサークル ID を採番する。
type CircleFactory struct{}

// NewCircleFactory はサークルファクトリを生成する。
func NewCircleFactory() *CircleFactory {
	return &CircleFactory{}
}

// Create は採番し、オーナーを最初のメンバーにしたサークルを作る。
func (f *CircleFactory) Create(name circle.CircleName, owner *user.User) (*circle.Circle, error) {
	// 保存前に ID を採番する
	id, err := newUUID()
	if err != nil {
		return nil, err
	}
	circleID, err := circle.NewCircleID(id)
	if err != nil {
		return nil, err
	}
	// オーナーを最初のメンバーとして含める
	return circle.NewCircle(circleID, name, owner.ID(), []user.UserID{owner.ID()})
}

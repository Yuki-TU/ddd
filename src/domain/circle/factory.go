package circle

import "github.com/Yuki-TU/ddd/domain/user"

// Factory はサークルの生成を担う。
// ID 採番と「オーナーを最初のメンバーにする」生成ルールをここに寄せる。
type Factory interface {
	Create(name CircleName, owner *user.User) (*Circle, error)
}

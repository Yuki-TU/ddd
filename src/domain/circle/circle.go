package circle

import (
	"github.com/Yuki-TU/ddd/domain/user"
)

// メンバー数の上限（オーナーを含む）
const memberLimit = 30

// Circle はサークル集約のルート。
// メンバーは User 本体ではなく UserID だけ持つ。
// こうすると、Circle 経由で User の名前変更などが起きない。
type Circle struct {
	id      CircleID
	name    CircleName
	ownerID user.UserID
	members []user.UserID
}

// NewCircle はサークル集約を生成する。ID・名前・オーナーは必須。
func NewCircle(id CircleID, name CircleName, ownerID user.UserID, members []user.UserID) (*Circle, error) {
	if id == (CircleID{}) {
		return nil, ErrCircleIDRequired
	}
	if name == (CircleName{}) {
		return nil, ErrCircleNameRequired
	}
	if ownerID == (user.UserID{}) {
		return nil, ErrOwnerRequired
	}
	// スライスは呼び出し側と共有しない
	copied := make([]user.UserID, len(members))
	copy(copied, members)
	return &Circle{id: id, name: name, ownerID: ownerID, members: copied}, nil
}

// ID はサークルを識別する固有のIDを返す。
func (c *Circle) ID() CircleID {
	return c.id
}

// Name はサークル名を返す。
func (c *Circle) Name() CircleName {
	return c.name
}

// OwnerID はオーナーのユーザIDを返す。
func (c *Circle) OwnerID() user.UserID {
	return c.ownerID
}

// Members はメンバーIDのコピーを返す（呼び出し側が内部スライスを書き換えないようにする）。
func (c *Circle) Members() []user.UserID {
	out := make([]user.UserID, len(c.members))
	copy(out, c.members)
	return out
}

// Equal は同一性を比較する。
func (c *Circle) Equal(other *Circle) bool {
	if c == nil || other == nil {
		return c == other
	}
	// 同一性は id で判定する
	return c.id.Equal(other.id)
}

// IsFull は満員かどうかを集約自身が判断する（ルールを外に出さない）。
func (c *Circle) IsFull() bool {
	return len(c.members) >= memberLimit
}

// Join はメンバー追加。集約ルート経由でのみメンバーを変える。
func (c *Circle) Join(userID user.UserID) error {
	if userID == (user.UserID{}) {
		return user.ErrUserIDRequired
	}
	// 既に参加しているか
	if c.hasMember(userID) {
		return ErrAlreadyMember
	}
	// 満員か
	if c.IsFull() {
		return ErrCircleFull
	}
	c.members = append(c.members, userID)
	return nil
}

// hasMember は指定ユーザが既にメンバーかを判定する。
func (c *Circle) hasMember(userID user.UserID) bool {
	for _, m := range c.members {
		if m.Equal(userID) {
			return true
		}
	}
	return false
}

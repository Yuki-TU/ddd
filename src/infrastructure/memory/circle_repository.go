package memory

import (
	"sync"

	"github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

// CircleRepository はテスト用のインメモリ実装。
type CircleRepository struct {
	mu    sync.Mutex
	store map[string]*circle.Circle
}

// NewCircleRepository は空のインメモリリポジトリを生成する。
func NewCircleRepository() *CircleRepository {
	return &CircleRepository{store: make(map[string]*circle.Circle)}
}

// FindByID は ID で再構築する。
func (r *CircleRepository) FindByID(id circle.CircleID) (*circle.Circle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.store[id.String()]
	if !ok {
		return nil, nil
	}
	return cloneCircle(c)
}

// FindByName はサークル名で再構築する。
func (r *CircleRepository) FindByName(name circle.CircleName) (*circle.Circle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.store {
		if c.Name().Equal(name) {
			return cloneCircle(c)
		}
	}
	return nil, nil
}

// Save は保存（永続化）する。同じ名前の別サークルなら重複エラー。
func (r *CircleRepository) Save(c *circle.Circle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, stored := range r.store {
		if stored.Name().Equal(c.Name()) && !stored.Equal(c) {
			return circle.ErrCircleAlreadyExists
		}
	}
	cloned, err := cloneCircle(c)
	if err != nil {
		return err
	}
	r.store[c.ID().String()] = cloned
	return nil
}

// 再構築時にコピーを返す（集約の中身をストアと共有しない）
func cloneCircle(c *circle.Circle) (*circle.Circle, error) {
	members := make([]user.UserID, len(c.Members()))
	copy(members, c.Members())
	return circle.NewCircle(c.ID(), c.Name(), c.OwnerID(), members)
}

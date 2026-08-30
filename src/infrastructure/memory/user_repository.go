package memory

import (
	"sync"

	domain "github.com/Yuki-TU/ddd/domain/user"
)

// UserRepository はテスト用のインメモリ実装。
// ドメインの Repository を実現するので、アプリケーションサービスは変更不要。
type UserRepository struct {
	mu    sync.Mutex
	store map[string]domain.User
}

// NewUserRepository は空のインメモリリポジトリを生成する。
func NewUserRepository() *UserRepository {
	return &UserRepository{
		store: make(map[string]domain.User),
	}
}

// FindByID は ID で再構築する。
func (r *UserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.store[id.String()]
	if !ok {
		return nil, nil
	}
	cloned := r.clone(user)
	return &cloned, nil
}

// FindByName はユーザ名で再構築する。
func (r *UserRepository) FindByName(name domain.UserName) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.store {
		if name.Equal(user.Name()) {
			cloned := r.clone(user)
			return &cloned, nil
		}
	}
	// 見つからないときはエラーにせず nil（再構築できなかった）
	return nil, nil
}

// Save は保存（永続化）する。同じ名前の別人なら重複エラー。
func (r *UserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ユニーク制約に相当するチェック（自分以外が同じ名前なら重複）
	for _, stored := range r.store {
		if stored.Name().Equal(user.Name()) && !stored.Equal(user) {
			return domain.ErrUserAlreadyExists
		}
	}
	r.store[user.ID().String()] = r.clone(*user)
	return nil
}

// Delete はユーザを削除する。
func (r *UserRepository) Delete(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, user.ID().String())
	return nil
}

// 値型なので代入でコピーになる（呼び出し側がストアを直接書き換えないようにする）
func (r *UserRepository) clone(user domain.User) domain.User {
	return user
}

package user

// Repository はユーザの永続化と再構築だけを担う。
// 存在判定などのドメイン知識はここには置かない（ドメインサービス側）。
// インターフェースにして、MySQL / メモリ実装を差し替えられるようにする。
type Repository interface {
	// 保存（永続化）。更新も Save で上書きする
	Save(user *User) error
	// ID で再構築する
	FindByID(id UserID) (*User, error)
	// 名前で再構築する
	FindByName(name UserName) (*User, error)
	// 削除はリポジトリの責務
	Delete(user *User) error
}

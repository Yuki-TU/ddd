package circle

// Repository はサークル集約の永続化と再構築だけを担う。
// リポジトリは集約の単位で作る。
type Repository interface {
	// 保存（永続化）
	Save(circle *Circle) error
	// ID で再構築する
	FindByID(id CircleID) (*Circle, error)
	// 名前で再構築する
	FindByName(name CircleName) (*Circle, error)
}

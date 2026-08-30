package user

// Factory はユーザの生成を担う。
// ID の採番（UUID や採番テーブル）をドメインオブジェクトから分離する。
// 実装を差し替えられるようインターフェースにする。
type Factory interface {
	Create(name UserName) (*User, error)
}

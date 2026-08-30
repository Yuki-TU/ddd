package user

// Service はユーザのドメインサービス。
// 「ユーザが自分自身の重複を確認する」のは不自然なので、ここに置く。
// リポジトリのインターフェースに依存する（実装には依存しない）。
type Service struct {
	users Repository
}

// NewService はドメインサービスを生成する。リポジトリのインターフェースを受け取る。
func NewService(users Repository) *Service {
	return &Service{users: users}
}

// Exists はユーザ名から再構築し、既に存在するかを判断する。
func (s *Service) Exists(u *User) (bool, error) {
	// ユーザ名からユーザを再構築する
	found, err := s.users.FindByName(u.Name())
	if err != nil {
		return false, err
	}
	if found == nil {
		return false, nil
	}
	// 自分自身は重複とみなさない（名前変更せずに更新したときのため）
	return !found.Equal(u), nil
}

package circle

// Service はサークルのドメインサービス。
// 名前の重複確認はサークル自身の振る舞いではないのでここに置く。
type Service struct {
	// サークルリポジトリのインターフェースに依存する
	circles Repository
}

// NewService はドメインサービスを生成する。
func NewService(circles Repository) *Service {
	return &Service{circles: circles}
}

// Exists はサークル名から再構築し、既に存在するかを判断する。
func (s *Service) Exists(c *Circle) (bool, error) {
	found, err := s.circles.FindByName(c.Name())
	if err != nil {
		return false, err
	}
	if found == nil {
		return false, nil
	}
	// 自分自身は重複とみなさない
	return !found.Equal(c), nil
}

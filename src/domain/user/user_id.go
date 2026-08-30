package user

// UserID はユーザを識別する固有のID（値オブジェクト）。
type UserID struct {
	value string
}

// NewUserID はユーザIDを生成する。空文字は許可しない。
func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, ErrUserIDRequired
	}
	return UserID{value: value}, nil
}

// String は永続化や表示用に内部の文字列を返す。
func (id UserID) String() string {
	return id.value
}

// Equal は等価性で比較する。
func (id UserID) Equal(other UserID) bool {
	return id.value == other.value
}

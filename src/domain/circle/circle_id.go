package circle

// CircleID はサークルを識別する固有のID（値オブジェクト）。
type CircleID struct {
	value string
}

// NewCircleID はサークルIDを生成する。空文字は許可しない。
func NewCircleID(value string) (CircleID, error) {
	if value == "" {
		return CircleID{}, ErrCircleIDRequired
	}
	return CircleID{value: value}, nil
}

// String は永続化や表示用に内部の文字列を返す。
func (id CircleID) String() string {
	return id.value
}

// Equal は等価性で比較する。
func (id CircleID) Equal(other CircleID) bool {
	return id.value == other.value
}

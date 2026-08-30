package user

// User はユーザエンティティ。
// 値オブジェクトと違い、名前が変わっても同一人物として扱う。
// 同一性は ID で判定する。
type User struct {
	id   UserID   // ユーザを識別する固有のID
	name UserName // 可変な属性（ニックネームなど）
}

// NewUser はユーザエンティティを生成する。ID と名前は必須。
func NewUser(id UserID, name UserName) (*User, error) {
	if id == (UserID{}) {
		return nil, ErrUserIDRequired
	}
	if name == (UserName{}) {
		return nil, ErrUserNameRequired
	}
	return &User{id: id, name: name}, nil
}

// ID はユーザを識別する固有のIDを返す。
func (u *User) ID() UserID {
	return u.id
}

// Name は現在のユーザ名を返す。
func (u *User) Name() UserName {
	return u.name
}

// ChangeName は制限をかけて名前を変える。
// 素の値をそのまま入れるセッターにはしない。
func (u *User) ChangeName(name UserName) error {
	if name == (UserName{}) {
		return ErrUserNameRequired
	}
	u.name = name
	return nil
}

// Equal は同一性を比較する。
func (u *User) Equal(other *User) bool {
	if u == nil || other == nil {
		return u == other
	}
	// 同一性は id で判定する（名前が同じでも別人になり得る）
	return u.id.Equal(other.id)
}

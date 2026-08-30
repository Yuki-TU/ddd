package mysql

import (
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"

	domain "github.com/Yuki-TU/ddd/domain/user"
)

// UserRepository は MySQL によるユーザリポジトリの実現。
// SQL などデータベース特有の書き方をしてよい。
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository は MySQL ユーザリポジトリを生成する。
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID は ID で再構築する。
func (r *UserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	// ID で再構築する
	return r.find(`SELECT id, name FROM users WHERE id = ?`, id.String())
}

// FindByName はユーザ名で再構築する。
func (r *UserRepository) FindByName(name domain.UserName) (*domain.User, error) {
	// 名前で再構築する
	return r.find(`SELECT id, name FROM users WHERE name = ?`, name.String())
}

// find は1行読んでユーザを再構築する。無いときは nil。
func (r *UserRepository) find(query, arg string) (*domain.User, error) {
	var id, storedName string
	err := r.db.QueryRow(query, arg).Scan(&id, &storedName)
	if errors.Is(err, sql.ErrNoRows) {
		// 再構築できなかった
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return reconstructUser(id, storedName)
}

// Save は保存（永続化）する。既存なら更新、無ければ挿入する。
func (r *UserRepository) Save(user *domain.User) error {
	// 更新も Save で上書きする。
	// 名前の UNIQUE に ON DUPLICATE KEY を使うと、別人の行を更新してしまうので分ける。
	// 既存かどうかを見て INSERT / UPDATE を分ける
	var n int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ?`, user.ID().String()).Scan(&n); err != nil {
		return err
	}
	var err error
	if n > 0 {
		_, err = r.db.Exec(
			`UPDATE users SET name = ? WHERE id = ?`,
			user.Name().String(),
			user.ID().String(),
		)
	} else {
		_, err = r.db.Exec(
			`INSERT INTO users (id, name) VALUES (?, ?)`,
			user.ID().String(),
			user.Name().String(),
		)
	}
	if isDuplicate(err) {
		return domain.ErrUserAlreadyExists
	}
	return err
}

// Delete はユーザを削除する。
func (r *UserRepository) Delete(user *domain.User) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, user.ID().String())
	if isForeignKey(err) {
		// オーナーが残っているサークルがあると外部キーで落ちる
		return domain.ErrCannotDeleteUser
	}
	return err
}

// 行データからドメインオブジェクトを再構築する
func reconstructUser(id, name string) (*domain.User, error) {
	userID, err := domain.NewUserID(id)
	if err != nil {
		return nil, err
	}
	userName, err := domain.NewUserName(name)
	if err != nil {
		return nil, err
	}
	return domain.NewUser(userID, userName)
}

// isDuplicate はユニーク制約違反（1062）かを判定する。
func isDuplicate(err error) bool {
	var me *mysqldriver.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}

// isForeignKey は外部キー制約違反（1451）かを判定する。
func isForeignKey(err error) bool {
	var me *mysqldriver.MySQLError
	return errors.As(err, &me) && me.Number == 1451
}

package mysql

import (
	"database/sql"
	"errors"

	"github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

// CircleRepository は MySQL によるサークルリポジトリの実現。
type CircleRepository struct {
	db *sql.DB
}

// NewCircleRepository は MySQL サークルリポジトリを生成する。
func NewCircleRepository(db *sql.DB) *CircleRepository {
	return &CircleRepository{db: db}
}

// FindByID は ID で再構築する。
func (r *CircleRepository) FindByID(id circle.CircleID) (*circle.Circle, error) {
	return r.find(`SELECT id, name, owner_id FROM circles WHERE id = ?`, id.String())
}

// FindByName はサークル名で再構築する。
func (r *CircleRepository) FindByName(name circle.CircleName) (*circle.Circle, error) {
	return r.find(`SELECT id, name, owner_id FROM circles WHERE name = ?`, name.String())
}

// find はサークル本体とメンバーを読んで集約を再構築する。
func (r *CircleRepository) find(query, arg string) (*circle.Circle, error) {
	var id, name, ownerID string
	err := r.db.QueryRow(query, arg).Scan(&id, &name, &ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		// 再構築できなかった
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 集約の部分（メンバー）もまとめて再構築する
	rows, err := r.db.Query(`SELECT user_id FROM circle_members WHERE circle_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []user.UserID
	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); err != nil {
			return nil, err
		}
		uid, err := user.NewUserID(memberID)
		if err != nil {
			return nil, err
		}
		members = append(members, uid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reconstructCircle(id, name, ownerID, members)
}

// Save はサークル集約を保存（永続化）する。
func (r *CircleRepository) Save(c *circle.Circle) error {
	// サークル本体とメンバーを同じトランザクションで永続化する
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 既存かどうかを見て INSERT / UPDATE を分ける
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM circles WHERE id = ?`, c.ID().String()).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		_, err = tx.Exec(
			`UPDATE circles SET name = ?, owner_id = ? WHERE id = ?`,
			c.Name().String(),
			c.OwnerID().String(),
			c.ID().String(),
		)
	} else {
		_, err = tx.Exec(
			`INSERT INTO circles (id, name, owner_id) VALUES (?, ?, ?)`,
			c.ID().String(),
			c.Name().String(),
			c.OwnerID().String(),
		)
	}
	if isDuplicate(err) {
		return circle.ErrCircleAlreadyExists
	}
	if err != nil {
		return err
	}

	// メンバーは一度消してから集約の状態で入れ直す
	if _, err := tx.Exec(`DELETE FROM circle_members WHERE circle_id = ?`, c.ID().String()); err != nil {
		return err
	}
	for _, member := range c.Members() {
		if _, err := tx.Exec(
			`INSERT INTO circle_members (circle_id, user_id) VALUES (?, ?)`,
			c.ID().String(),
			member.String(),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// reconstructCircle は行データからサークル集約を再構築する。
func reconstructCircle(id, name, ownerID string, members []user.UserID) (*circle.Circle, error) {
	circleID, err := circle.NewCircleID(id)
	if err != nil {
		return nil, err
	}
	circleName, err := circle.NewCircleName(name)
	if err != nil {
		return nil, err
	}
	owner, err := user.NewUserID(ownerID)
	if err != nil {
		return nil, err
	}
	return circle.NewCircle(circleID, circleName, owner, members)
}

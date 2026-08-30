package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// compose.yaml の MySQL に合わせた接続文字列
const defaultDSN = "ddd:ddd@tcp(127.0.0.1:3306)/ddd?parseTime=true&charset=utf8mb4&timeout=2s"

// DSN は MySQL 接続文字列を返す。MYSQL_DSN があればそれを使う。
func DSN() string {
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return dsn
	}
	return defaultDSN
}

// Open は MySQL に1回接続し、テーブルを用意する。
func Open(dsn string) (*sql.DB, error) {
	return open(dsn, 1)
}

// OpenWithRetry は起動直後の MySQL 待ち用。
func OpenWithRetry(dsn string) (*sql.DB, error) {
	return open(dsn, 30)
}

// open は接続を試し、成功したらマイグレーションする。
func open(dsn string, attempts int) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 起動直後は接続できないことがあるので、指定回数まで待つ
	var pingErr error
	for range attempts {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr = db.PingContext(ctx)
		cancel()
		if pingErr == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if pingErr != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql ping: %w", pingErr)
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// migrate は学習用のテーブルを作る。
func migrate(db *sql.DB) error {
	// ユニーク制約はドメインサービスの重複確認と両方置く
	// （同時登録の穴埋め + コード上でルールが見えるようにする）
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uk_users_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS circles (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			owner_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uk_circles_name (name),
			CONSTRAINT fk_circles_owner FOREIGN KEY (owner_id) REFERENCES users (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS circle_members (
			circle_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (circle_id, user_id),
			CONSTRAINT fk_members_circle FOREIGN KEY (circle_id) REFERENCES circles (id) ON DELETE CASCADE,
			CONSTRAINT fk_members_user FOREIGN KEY (user_id) REFERENCES users (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

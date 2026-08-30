package mysql_test

import (
	"errors"
	"testing"

	domain "github.com/Yuki-TU/ddd/domain/user"
	"github.com/Yuki-TU/ddd/infrastructure/id"
	"github.com/Yuki-TU/ddd/infrastructure/mysql"
)

// リポジトリ実装のテスト。MySQL が無いときはスキップする。
func TestUserRepository_SaveAndFind(t *testing.T) {
	db, err := mysql.Open(mysql.DSN())
	if err != nil {
		t.Skipf("MySQL に接続できないためスキップ: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := mysql.NewUserRepository(db)
	factory := id.NewUserFactory()

	name, err := domain.NewUserName("taro-mysql-test")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = db.Exec(`DELETE FROM users WHERE name IN (?, ?)`, name.String(), "jiro-mysql-test")

	user, err := factory.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(user); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM users WHERE id = ?`, user.ID().String())
	})

	got, err := repo.FindByName(name)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || !got.Equal(user) {
		t.Fatalf("got %+v", got)
	}

	byID, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if byID == nil || !byID.Equal(user) {
		t.Fatalf("got %+v", byID)
	}

	otherName, err := domain.NewUserName("jiro-mysql-test")
	if err != nil {
		t.Fatal(err)
	}
	other, err := factory.Create(otherName)
	if err != nil {
		t.Fatal(err)
	}
	if err := other.ChangeName(name); err != nil {
		t.Fatal(err)
	}
	err = repo.Save(other)
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Fatalf("duplicate name: %v", err)
	}
}

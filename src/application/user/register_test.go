package user_test

import (
	"errors"
	"testing"

	userapp "github.com/Yuki-TU/ddd/application/user"
	domain "github.com/Yuki-TU/ddd/domain/user"
	"github.com/Yuki-TU/ddd/infrastructure/id"
	"github.com/Yuki-TU/ddd/infrastructure/memory"
)

// アプリケーションサービスのテストはメモリリポジトリでよい。
// インターフェースに依存しているので、本番の MySQL 実装に差し替えなくて済む。
func newUserServices() (*userapp.RegisterService, *userapp.GetService, *userapp.UpdateService, *userapp.DeleteService) {
	users := memory.NewUserRepository()
	userService := domain.NewService(users)
	return userapp.NewRegisterService(users, id.NewUserFactory(), userService),
		userapp.NewGetService(users),
		userapp.NewUpdateService(users, userService),
		userapp.NewDeleteService(users)
}

// TestRegisterService_Register は登録ユースケースの正常系と重複・空名前を確認する。
func TestRegisterService_Register(t *testing.T) {
	t.Parallel()

	t.Run("重複がなければ登録できる", func(t *testing.T) {
		t.Parallel()
		register, _, _, _ := newUserServices()
		got, err := register.Register(userapp.RegisterCommand{Name: "taro"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "taro" || got.ID == "" {
			t.Fatalf("got %+v", got)
		}
	})

	t.Run("同じユーザ名は登録できない", func(t *testing.T) {
		t.Parallel()
		register, _, _, _ := newUserServices()
		if _, err := register.Register(userapp.RegisterCommand{Name: "taro"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err := register.Register(userapp.RegisterCommand{Name: "taro"})
		if !errors.Is(err, domain.ErrUserAlreadyExists) {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("空のユーザ名は登録できない", func(t *testing.T) {
		t.Parallel()
		register, _, _, _ := newUserServices()
		_, err := register.Register(userapp.RegisterCommand{Name: ""})
		if !errors.Is(err, domain.ErrInvalidUserName) {
			t.Fatalf("got %v", err)
		}
	})
}

// TestGetService_Get は ID / 名前検索と未存在を確認する。
func TestGetService_Get(t *testing.T) {
	t.Parallel()
	register, get, _, _ := newUserServices()
	created, err := register.Register(userapp.RegisterCommand{Name: "taro"})
	if err != nil {
		t.Fatal(err)
	}

	byID, err := get.Get(userapp.GetCommand{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if byID != created {
		t.Fatalf("got %+v", byID)
	}

	byName, err := get.Get(userapp.GetCommand{Name: "taro"})
	if err != nil {
		t.Fatal(err)
	}
	if byName != created {
		t.Fatalf("got %+v", byName)
	}

	_, err = get.Get(userapp.GetCommand{ID: "missing-id"})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("got %v", err)
	}
}

// TestUpdateAndDelete は名前変更と退会を確認する。
func TestUpdateAndDelete(t *testing.T) {
	t.Parallel()
	register, get, update, deleteUser := newUserServices()
	created, err := register.Register(userapp.RegisterCommand{Name: "taro"})
	if err != nil {
		t.Fatal(err)
	}

	if err := update.Update(userapp.UpdateCommand{ID: created.ID, Name: "jiro"}); err != nil {
		t.Fatal(err)
	}
	got, err := get.Get(userapp.GetCommand{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "jiro" {
		t.Fatalf("got %+v", got)
	}

	if err := deleteUser.Delete(userapp.DeleteCommand{ID: created.ID}); err != nil {
		t.Fatal(err)
	}
	_, err = get.Get(userapp.GetCommand{ID: created.ID})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("got %v", err)
	}
}

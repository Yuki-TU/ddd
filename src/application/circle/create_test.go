package circle_test

import (
	"errors"
	"testing"

	circleapp "github.com/Yuki-TU/ddd/application/circle"
	userapp "github.com/Yuki-TU/ddd/application/user"
	circledomain "github.com/Yuki-TU/ddd/domain/circle"
	userdomain "github.com/Yuki-TU/ddd/domain/user"
	"github.com/Yuki-TU/ddd/infrastructure/id"
	"github.com/Yuki-TU/ddd/infrastructure/memory"
)

// TestCircleCreateAndJoin はサークル作成・重複名・参加・既参加を確認する。
func TestCircleCreateAndJoin(t *testing.T) {
	t.Parallel()

	// ユースケースの組み立ては main と同じ。リポジトリだけメモリに差し替える。
	users := memory.NewUserRepository()
	circles := memory.NewCircleRepository()
	userService := userdomain.NewService(users)
	register := userapp.NewRegisterService(users, id.NewUserFactory(), userService)
	create := circleapp.NewCreateService(circles, users, id.NewCircleFactory(), circledomain.NewService(circles))
	get := circleapp.NewGetService(circles)
	join := circleapp.NewJoinService(circles, users)

	owner, err := register.Register(userapp.RegisterCommand{Name: "taro"})
	if err != nil {
		t.Fatal(err)
	}
	member, err := register.Register(userapp.RegisterCommand{Name: "jiro"})
	if err != nil {
		t.Fatal(err)
	}

	created, err := create.Create(circleapp.CreateCommand{UserID: owner.ID, Name: "surfing"})
	if err != nil {
		t.Fatal(err)
	}
	if created.OwnerID != owner.ID || len(created.Members) != 1 {
		t.Fatalf("got %+v", created)
	}

	_, err = create.Create(circleapp.CreateCommand{UserID: owner.ID, Name: "surfing"})
	if !errors.Is(err, circledomain.ErrCircleAlreadyExists) {
		t.Fatalf("got %v", err)
	}

	if err := join.Join(circleapp.JoinCommand{CircleID: created.ID, UserID: member.ID}); err != nil {
		t.Fatal(err)
	}
	got, err := get.Get(circleapp.GetCommand{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Members) != 2 {
		t.Fatalf("got %+v", got)
	}

	err = join.Join(circleapp.JoinCommand{CircleID: created.ID, UserID: member.ID})
	if !errors.Is(err, circledomain.ErrAlreadyMember) {
		t.Fatalf("got %v", err)
	}
}

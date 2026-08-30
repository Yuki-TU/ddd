package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	circleapp "github.com/Yuki-TU/ddd/application/circle"
	userapp "github.com/Yuki-TU/ddd/application/user"
	circledomain "github.com/Yuki-TU/ddd/domain/circle"
	userdomain "github.com/Yuki-TU/ddd/domain/user"
	"github.com/Yuki-TU/ddd/infrastructure/id"
	"github.com/Yuki-TU/ddd/infrastructure/memory"
	httpapi "github.com/Yuki-TU/ddd/presentation/http"
)

// HTTP のテストもメモリリポジトリで組み立てる。
// Controller はユースケースを呼ぶだけなので、DB がなくても経路を確認できる。
func newTestServer() http.Handler {
	users := memory.NewUserRepository()
	circles := memory.NewCircleRepository()
	userService := userdomain.NewService(users)
	circleService := circledomain.NewService(circles)
	return httpapi.NewServer(
		userapp.NewRegisterService(users, id.NewUserFactory(), userService),
		userapp.NewGetService(users),
		userapp.NewUpdateService(users, userService),
		userapp.NewDeleteService(users),
		circleapp.NewCreateService(circles, users, id.NewCircleFactory(), circleService),
		circleapp.NewGetService(circles),
		circleapp.NewJoinService(circles, users),
	).Handler()
}

// TestUserAndCircleAPI は登録・検索・サークル作成・参加の HTTP 経路を確認する。
func TestUserAndCircleAPI(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(newTestServer())
	t.Cleanup(ts.Close)
	client := ts.Client()

	resp, err := client.Post(ts.URL+"/users", "application/json", bytes.NewBufferString(`{"name":"taro"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var owner userapp.UserData
	if err := json.NewDecoder(resp.Body).Decode(&owner); err != nil {
		t.Fatal(err)
	}

	resp, err = client.Get(ts.URL + "/users?name=taro")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	resp, err = client.Get(ts.URL + "/users/" + owner.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	resp, err = client.Post(ts.URL+"/users", "application/json", bytes.NewBufferString(`{"name":"jiro"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var member userapp.UserData
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		t.Fatal(err)
	}

	body := `{"userId":"` + owner.ID + `","name":"surfing"}`
	resp, err = client.Post(ts.URL+"/circles", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var circle circleapp.CircleData
	if err := json.NewDecoder(resp.Body).Decode(&circle); err != nil {
		t.Fatal(err)
	}

	joinBody := `{"userId":"` + member.ID + `"}`
	resp, err = client.Post(ts.URL+"/circles/"+circle.ID+"/members", "application/json", bytes.NewBufferString(joinBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

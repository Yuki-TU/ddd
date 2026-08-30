package main

import (
	"log"
	"net/http"
	"os"

	circleapp "github.com/Yuki-TU/ddd/application/circle"
	userapp "github.com/Yuki-TU/ddd/application/user"
	circledomain "github.com/Yuki-TU/ddd/domain/circle"
	userdomain "github.com/Yuki-TU/ddd/domain/user"
	"github.com/Yuki-TU/ddd/infrastructure/id"
	"github.com/Yuki-TU/ddd/infrastructure/mysql"
	httpapi "github.com/Yuki-TU/ddd/presentation/http"
)

func main() {
	// 起動時に依存関係を組み立てる（コンストラクタ注入）
	db, err := mysql.OpenWithRetry(mysql.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// リポジトリの実現（MySQL）
	users := mysql.NewUserRepository(db)
	circles := mysql.NewCircleRepository(db)
	// ドメインサービスはリポジトリのインターフェースを受け取る
	userService := userdomain.NewService(users)
	circleService := circledomain.NewService(circles)

	// アプリケーションサービス（ユースケース）を組み立てて Controller に渡す
	srv := httpapi.NewServer(
		userapp.NewRegisterService(users, id.NewUserFactory(), userService),
		userapp.NewGetService(users),
		userapp.NewUpdateService(users, userService),
		userapp.NewDeleteService(users),
		circleapp.NewCreateService(circles, users, id.NewCircleFactory(), circleService),
		circleapp.NewGetService(circles),
		circleapp.NewJoinService(circles, users),
	)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	log.Printf("listen %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}

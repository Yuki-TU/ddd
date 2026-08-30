package httpapi

import (
	"net/http"

	circleapp "github.com/Yuki-TU/ddd/application/circle"
	userapp "github.com/Yuki-TU/ddd/application/user"
)

// Server はクライアントとアプリケーションサービスの仲介（Controller）。
// リクエストをコマンドに変え、ユースケースを呼び、HTTP として返すだけ。
// ドメインルールは書かない。
type Server struct {
	registerUser *userapp.RegisterService
	getUser      *userapp.GetService
	updateUser   *userapp.UpdateService
	deleteUser   *userapp.DeleteService
	createCircle *circleapp.CreateService
	getCircle    *circleapp.GetService
	joinCircle   *circleapp.JoinService
}

// NewServer は各ユースケースを受け取って Controller を生成する。
func NewServer(
	registerUser *userapp.RegisterService,
	getUser *userapp.GetService,
	updateUser *userapp.UpdateService,
	deleteUser *userapp.DeleteService,
	createCircle *circleapp.CreateService,
	getCircle *circleapp.GetService,
	joinCircle *circleapp.JoinService,
) *Server {
	return &Server{
		registerUser: registerUser,
		getUser:      getUser,
		updateUser:   updateUser,
		deleteUser:   deleteUser,
		createCircle: createCircle,
		getCircle:    getCircle,
		joinCircle:   joinCircle,
	}
}

// Handler は HTTP ルートを組み立てる。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /users", s.createUser)
	mux.HandleFunc("GET /users/{id}", s.getUserByID)
	mux.HandleFunc("GET /users", s.searchUser)
	mux.HandleFunc("PUT /users/{id}", s.updateUserByID)
	mux.HandleFunc("DELETE /users/{id}", s.deleteUserByID)
	mux.HandleFunc("POST /circles", s.createCircleHandler)
	mux.HandleFunc("GET /circles/{id}", s.getCircleByID)
	mux.HandleFunc("POST /circles/{id}/members", s.joinCircleHandler)
	return mux
}

// health は起動確認用。
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type nameBody struct {
	Name string `json:"name"`
}

type userIDBody struct {
	UserID string `json:"userId"`
}

type createCircleBody struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

// createUser はユーザ登録リクエストをコマンドに変えてユースケースへ渡す。
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var body nameBody
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "不正なリクエストです"})
		return
	}
	// リクエストからコマンドを作り、ユースケースに渡す
	got, err := s.registerUser.Register(userapp.RegisterCommand{Name: body.Name})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, got)
}

// getUserByID は ID でユーザを検索する。
func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	got, err := s.getUser.Get(userapp.GetCommand{ID: r.PathValue("id")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, got)
}

// searchUser は名前でユーザを検索する。
func (s *Server) searchUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "name を指定してください"})
		return
	}
	got, err := s.getUser.Get(userapp.GetCommand{Name: name})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, got)
}

// updateUserByID はユーザ更新リクエストをコマンドに変えてユースケースへ渡す。
func (s *Server) updateUserByID(w http.ResponseWriter, r *http.Request) {
	var body nameBody
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "不正なリクエストです"})
		return
	}
	if err := s.updateUser.Update(userapp.UpdateCommand{ID: r.PathValue("id"), Name: body.Name}); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteUserByID は退会リクエストをユースケースへ渡す。
func (s *Server) deleteUserByID(w http.ResponseWriter, r *http.Request) {
	if err := s.deleteUser.Delete(userapp.DeleteCommand{ID: r.PathValue("id")}); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// createCircleHandler はサークル作成リクエストをコマンドに変えてユースケースへ渡す。
func (s *Server) createCircleHandler(w http.ResponseWriter, r *http.Request) {
	var body createCircleBody
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "不正なリクエストです"})
		return
	}
	got, err := s.createCircle.Create(circleapp.CreateCommand{UserID: body.UserID, Name: body.Name})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, got)
}

// getCircleByID は ID でサークルを取得する。
func (s *Server) getCircleByID(w http.ResponseWriter, r *http.Request) {
	got, err := s.getCircle.Get(circleapp.GetCommand{ID: r.PathValue("id")})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, got)
}

// joinCircleHandler はサークル参加リクエストをコマンドに変えてユースケースへ渡す。
func (s *Server) joinCircleHandler(w http.ResponseWriter, r *http.Request) {
	var body userIDBody
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "不正なリクエストです"})
		return
	}
	if err := s.joinCircle.Join(circleapp.JoinCommand{CircleID: r.PathValue("id"), UserID: body.UserID}); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Yuki-TU/ddd/domain/circle"
	"github.com/Yuki-TU/ddd/domain/user"
)

type errorBody struct {
	Error string `json:"error"`
}

// decodeJSON はリクエストボディを構造体へ読む。
func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

// writeJSON は JSON レスポンスを書く。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError はドメインエラーを HTTP エラー応答にする。
func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, statusOf(err), errorBody{Error: err.Error()})
}

// statusOf はドメインエラーを HTTP ステータスに変換する。
// ドメイン層は HTTP を知らない。変換はプレゼンテーション層の仕事。
func statusOf(err error) int {
	switch {
	case errors.Is(err, user.ErrUserNotFound), errors.Is(err, circle.ErrCircleNotFound):
		return http.StatusNotFound
	case errors.Is(err, user.ErrUserAlreadyExists),
		errors.Is(err, circle.ErrCircleAlreadyExists),
		errors.Is(err, circle.ErrAlreadyMember),
		errors.Is(err, circle.ErrCircleFull),
		errors.Is(err, user.ErrCannotDeleteUser):
		return http.StatusConflict
	case errors.Is(err, user.ErrInvalidUserName),
		errors.Is(err, user.ErrUserIDRequired),
		errors.Is(err, user.ErrUserNameRequired),
		errors.Is(err, circle.ErrInvalidCircleName),
		errors.Is(err, circle.ErrCircleIDRequired),
		errors.Is(err, circle.ErrCircleNameRequired),
		errors.Is(err, circle.ErrOwnerRequired):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

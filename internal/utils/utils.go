package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type UserAuthKey int8

func UserIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(UserAuthKey(0))
	id, ok := v.(string)
	return id, ok
}

func RespondWithError(resp http.ResponseWriter, code int, message string) {
	RespondWithJSON(resp, code, map[string]string{"error": message})
}

func RespondWithJSON(resp http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	resp.Write(response)
}

func Value(req *http.Request, p string) sql.NullString {
	return sql.NullString{
		String: req.FormValue(p),
		Valid:  true,
	}
}
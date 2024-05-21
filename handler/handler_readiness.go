package handler

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {

	RespondWithJSON(w, 200, struct{}{})
}

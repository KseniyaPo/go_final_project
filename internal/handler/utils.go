package handler

import (
	"io"
	"net/http"

	"go_final_project/pkg/api"
)

type UtilsHandler struct{}

func NewUtilsHandler() *UtilsHandler {
	return &UtilsHandler{}
}

func (u *UtilsHandler) NextDate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now, err := api.ParseDate(r.FormValue("now"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		date, err := api.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
		if err != nil {
			api.Failed(w, err, http.StatusBadRequest)
			return
		}

		io.WriteString(w, date)
	})
}

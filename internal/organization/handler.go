package organization

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hell sir"))
		return
	})

	return mux
}

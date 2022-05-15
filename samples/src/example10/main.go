package main

import "net/http"

type Handler struct { // want `struct Handler doesn't verify interface compliance for http.Handler`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

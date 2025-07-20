package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte("Retrieving book with ID: " + id))
	})

	mux.HandleFunc("GET /books/dir/{d...}", func(w http.ResponseWriter, r *http.Request) {
		dir := r.PathValue("d")
		w.Write([]byte("Retrieving books in directory: " + dir))
	})

	http.ListenAndServe(":8080", mux)
}

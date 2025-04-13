package main

import "net/http"

func main() {
	fileServer := http.FileServer(http.Dir("./public"))

	http.Handle("/", fileServer)
	http.ListenAndServe(":8080", nil)
}

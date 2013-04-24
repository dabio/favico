package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", track(index))
	http.HandleFunc("/favicon", track(favicon))

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "5000"
	}

	fmt.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

/*
Serves the favicon from the following sources:

1. Memcache
2. Postgres
3. Domain-Root
4. Page Link
5. Default Icon
*/

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "/favicon?domain=google.com")
}

func favicon(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, req.FormValue("domain"))
}

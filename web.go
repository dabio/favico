package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	log.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

/*
Serves the favicon from the following sources:

1. Memcache
2. Postgres
3. Googles Service: https://plus.google.com/_/favicon?domain=flickr.com
3. Domain-Root
4. Page Link
5. Default Icon
*/

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/favicon?domain=google.com")
}

func favicon(w http.ResponseWriter, r *http.Request) {
	icon, err := fromGoogle(r.FormValue("domain"))

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	fmt.Fprintf(w, "%s", icon)
}

func fromGoogle(domain string) ([]byte, error) {
	log.Println(fmt.Sprintf("Fetch from Google for domain '%s'", domain))

	service := fmt.Sprintf("https://plus.google.com/_/favicon?domain=%s", domain)
	response, err := http.Get(service)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

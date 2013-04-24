package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var remote_service_url string

func init() {
	remote_service_url = "https://plus.google.com/_/favicon?domain=%s"
}

func main() {
	http.HandleFunc("/", track(index))
	http.HandleFunc("/favicon", track(favicon))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// index shows the homepage. A small reminder how to use this service.
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/favicon?domain=google.com")
}

// favicon tries to get the favicon from these sources:
// 1. Memcache
// 2. Google Service
func favicon(w http.ResponseWriter, r *http.Request) {
	domain := r.FormValue("domain")

	source := "Google"
	icon, err := fromGoogle(domain)
	if err != nil {
		panic(err)
	}

	w.Header().Set("X-Source", source)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	fmt.Fprintf(w, "%s", icon)

	go saveIcon(icon)
}

// fromGoogle connects to the google favicon service and tries to fetch the
// favicon.
func fromGoogle(domain string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf(remote_service_url, domain))
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

// saveIcon tries to save the icon to a memcache storage. The call of this
// function should be done in a Goroutine.
func saveIcon(icon []byte) {
	log.Println("Save to memcached")
}

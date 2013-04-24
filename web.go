package main

import (
	"fmt"
	"github.com/bmizerany/mc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var cn mc.Conn

// init sets the variables needed for our program
func init() {
	cn, err := mc.Dial("tcp", os.Getenv("MEMCACHIER_SERVERS"))
	if err != nil {
		panic(err)
	}

	err = cn.Auth(os.Getenv("MEMCACHIER_USERNAME"), os.Getenv("MEMCACHIER_PASSWORD"))
	if err != nil {
		panic(err)
	}
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

	source := "Cache"
	icon, err := fromCache(domain)
	if err != nil {
		source = "Google"
		icon, err = fromGoogle(domain)
		if err != nil {
			panic(err)
		}
	}

	w.Header().Set("X-Source", source)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	fmt.Fprintf(w, "%s", icon)

	go saveIcon(domain, icon)
}

func fromCache(domain string) ([]byte, error) {
	val, _, _, err := cn.Get(domain)
	if err != nil {
		panic(err)
	}

	return []byte(val), err
}

// fromGoogle connects to the google favicon service and tries to fetch the
// favicon.
func fromGoogle(domain string) ([]byte, error) {
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

// saveIcon tries to save the icon to a memcache storage. The call of this
// function should be done in a Goroutine.
func saveIcon(domain string, icon []byte) {
	err := cn.Set(domain, string(icon), 0, 0, 86400)
	if err != nil {
		return
	}

	log.Println("Save to memcached")
}

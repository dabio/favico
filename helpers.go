package main

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s request took %s", name, elapsed)
}

func track(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer timeTrack(time.Now(), getFunctionName(fn))
		fn(w, req)
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

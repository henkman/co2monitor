package main

import (
	"flag"
	"net/http"

	"github.com/henkman/co2monitor/web"
)

func main() {
	var config struct {
		listen string
	}
	flag.StringVar(&config.listen, "l", ":8080", "listen address")
	flag.Parse()

	var jcm web.JsonCO2Monitor
	if err := jcm.Open(); err != nil {
		panic(err)
	}
	defer jcm.Close()
	go jcm.Run()

	http.Handle("/read", &jcm)
	http.Handle("/", http.FileServer(http.Dir("www")))

	if err := http.ListenAndServe(config.listen, nil); err != nil {
		panic(err)
	}
}

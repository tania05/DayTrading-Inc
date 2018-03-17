package main

import (
  "net/http"
  "io"
  "encoding/json"
  "github.com/gorilla/mux"
  "log"
  "webserver/internal/money"
  "webserver/internal/trigger"
  "webserver/internal/logger"
  "webserver/internal/context"
  "io/ioutil"
	"webserver/internal/config"
	"strconv"
)

var ipMap []string = {"128.333.3.3:9090"}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello World!")
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	ip = ipMap[0];

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Write(ip)
}

func RegisterIPHandler(w http.ResponseWriter, r *http.Request) {
	var payload AddCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.Path("/users").Methods("GET").HandlerFunc(GetServerHandler);
	r.Path("/register/{ip}/{port}").Methods("POST").HandlerFunc(RegisterIPHandler);

	http.Handle("/", r)
	addr := ":5555"
	log.Fatal(http.ListenAndServe(addr, nil))
}

package main

import (
  "net/http"
  "common/config"
  "io"
  "github.com/gorilla/mux"
  "log"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "strconv"
)

type RegisterServerCommand struct {
	IP string
	Port int
}

var ipMap []string = []string{}
var rr int = 0

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello World!")
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	//body, _ := ioutil.ReadAll(r.Body)
	//err := json.Unmarshal(body, &payload)
	defer r.Body.Close()
	if (len(ipMap) < 1){
		fmt.Println("no servers registered")
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		ip := ipMap[rr];
		rr = (rr + 1) % len(ipMap)

		fmt.Println("sent back")
		fmt.Println(ip);
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ip))
	}	
}

func RegisterIPHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterServerCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	link := payload.IP + ":" + strconv.Itoa(payload.Port)
	ipMap = append(ipMap, link)
	w.WriteHeader(http.StatusOK)
	fmt.Println("server registered: " + link)
	fmt.Println(ipMap)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.Path("/users").Methods("GET").HandlerFunc(GetServerHandler);
	r.Path("/register").Methods("POST").HandlerFunc(RegisterIPHandler);

	http.Handle("/", r)
	port := config.GlobalConfig.LoadBalancer.Port
	addr := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

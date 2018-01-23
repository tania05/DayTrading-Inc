package main

import (
	"net"
	"log"
	"strconv"
)

const domain = "localhost"
const port = 8080

func getQuote(user string) string {
	addr := (domain + ":" + strconv.Itoa(port))

	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	conn.Write([]byte(user))
	conn.Write([]byte("\n"))

	buff := make([]byte, 1024)
	n, _ := conn.Readline(buff)
	log.Printf("Recieve: %s", buff[:n])
	return string(buff)
}

func main() {

	c := getQuote("TST,user")
	log.Printf("%s", c)
}
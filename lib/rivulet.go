package rivulet

import (
	"fmt"
	"net"
	"net/http"
)

func http_chat_handler(w http.ResponseWriter, r *http.Request, server *Server) {
	fmt.Fprintf(w, "Debug: /%s", r.URL.Path[1:])
	server.Broadcast("Request URL received: " + r.URL.Path[1:] + "\n")
}

func NewRivulet(pwd string) {
	fmt.Println("Got working directory: " + pwd)
	db, err := NewDatabase(pwd)
	if err != nil {
		fmt.Println("Error initializing database.")
		return
	}

	server := NewServer("default", *db)
	listener, _ := net.Listen("tcp", ":6666")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http_chat_handler(w, r, server)
	})
	go http.ListenAndServe(":8066", nil)

	for {
		conn, _ := listener.Accept()
		server.joins <- conn
	}
}

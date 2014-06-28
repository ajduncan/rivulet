package rivulet

import (
	"fmt"
	"net"
)

type Server struct {
	id       string
	db       DB
	clients  []*RivuletClient
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (server *Server) Broadcast(data string) {
	for _, client := range server.clients {
		client.outgoing <- data
	}
}

func (server *Server) Join(connection net.Conn) {
	client := NewClient(connection)
	server.clients = append(server.clients, client)
	fmt.Println("Connection from: ", connection.RemoteAddr().String())

	a, ok := server.db.assets["motd.txt"]
	if ok {
		client.outgoing <- string(a.data)
	}

	go func() {
		for {
			server.incoming <- <-client.incoming
		}
	}()
}

func (server *Server) Listen() {
	go func() {
		for {
			select {
			case data := <-server.incoming:
				server.Broadcast(data)
			case conn := <-server.joins:
				server.Join(conn)
			}
		}
	}()
}

func NewServer(name string, db DB) *Server {
	server := &Server{
		id:       name,
		db:       db,
		clients:  make([]*RivuletClient, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	server.Listen()

	return server
}

package rivulet

import (
	"fmt"
	"net"
)

type RivuletServer struct {
	id       string
	db       DB
	clients  []*RivuletClient
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (server *RivuletServer) Broadcast(data string) {
	for _, client := range server.clients {
		client.outgoing <- data
	}
}

func (server *RivuletServer) Join(connection net.Conn) {
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

func (server *RivuletServer) Listen() {
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

func (server *RivuletServer) Run() {
	listener, _ := net.Listen("tcp", ":6666")

	for {
		conn, _ := listener.Accept()
		server.joins <- conn
	}
}

func NewRivuletServer(name string, db DB) *RivuletServer {
	server := &RivuletServer{
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

/* Just a simple echo chat server, heavily cribbed from from https://gist.github.com/drewolson/3950226
 * 
 */

package main

import (
	"bufio"
	"fmt"
	"net"
)

type RivuletClient struct {
	id		 string
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

type Channel struct {
	id		 string
	clients  []*RivuletClient
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (client *RivuletClient) Read() {
	for {
		line, _ := client.reader.ReadString('\n')
		client.incoming <- client.id + " " + line
	}
}

func (client *RivuletClient) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *RivuletClient) Listen() {
	go client.Read()
	go client.Write()
}

func (channel *Channel) Broadcast(data string) {
	for _, client := range channel.clients {
		client.outgoing <- data
	}
}

func (channel *Channel) Join(connection net.Conn) {
	client := NewClient(connection)
	channel.clients = append(channel.clients, client)
	go func() {
		for {
			channel.incoming <- <-client.incoming
		}
	}()
}

func (channel *Channel) Listen() {
	go func() {
		for {
			select {
			case data := <-channel.incoming:
				channel.Broadcast(data)
			case conn := <-channel.joins:
				channel.Join(conn)
			}
		}
	}()
}

func NewClient(connection net.Conn) *RivuletClient {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &RivuletClient{
		id: connection.RemoteAddr().String(),
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()

	return client
}

func NewChannel(name string) *Channel {
	channel := &Channel{
		id:		  name,
		clients:  make([]*RivuletClient, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	channel.Listen()

	return channel
}

func main() {
	channel := NewChannel("default")

	listener, _ := net.Listen("tcp", ":6666")

	for {
		conn, _ := listener.Accept()
		fmt.Println("Connection from: ", conn.RemoteAddr().String())
		channel.joins <- conn
	}
}
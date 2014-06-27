/* Just a simple echo chat server, heavily cribbed from from https://gist.github.com/drewolson/3950226
 * 
 */

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

type Asset struct {
	id		string
	data	[]byte
}

type DB struct {
	assets map[string]Asset
}

func (db *DB) load_assets(path string) (error) {
	// initialize assets;
	db_assets := make(map[string]Asset)

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		filedata, err := ioutil.ReadFile(path + "/" + f.Name())
		if err == nil {
			a := &Asset{id: f.Name(), data: filedata}
			db_assets[f.Name()] = *a
		}
	}
	db.assets = db_assets
	return nil
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

func NewServer() {
	channel := NewChannel("default")

	listener, _ := net.Listen("tcp", ":6666")

	for {
		conn, _ := listener.Accept()
		fmt.Println("Connection from: ", conn.RemoteAddr().String())
		channel.joins <- conn
	}
}

func (db *DB) print_assets() {
	for key, value := range db.assets {
		fmt.Println("Key: ", key)
		fmt.Println("Value: ")
		fmt.Println(string(value.data))
	}
}

func (db *DB) print_motd() {
	a, ok := db.assets["motd.txt"]
	if ok {
		fmt.Println("Message of the day:")
		fmt.Println(string(a.data))
	}
}

func main() {
	db := DB{}
	err := db.load_assets("./static/db")
	if err != nil {
		fmt.Println("Error reading assets.")
		return
	}
	db.print_motd()
}

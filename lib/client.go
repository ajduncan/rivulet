package rivulet

import (
	"bufio"
	"net"
)

type RivuletClient struct {
	id       string
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
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

func NewClient(connection net.Conn) *RivuletClient {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &RivuletClient{
		id:       connection.RemoteAddr().String(),
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()

	return client
}

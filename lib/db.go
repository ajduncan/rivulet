package rivulet

import (
	"fmt"
	"io/ioutil"
	"net"
)

type Asset struct {
	id   string
	data []byte
}

type DB struct {
	assets   map[string]Asset
	connects chan net.Conn
}

func (db *DB) load_assets(path string) error {
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

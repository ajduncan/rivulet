package rivulet

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	_ "github.com/mattn/go-sqlite3"
)

type Asset struct {
	id   string
	data []byte
}

type DB struct {
	root     string
	assets   map[string]Asset
	rdb      *sql.DB
	connects chan net.Conn
}

func (db *DB) load_assets(path string) error {
	// initialize assets;
	db_assets := make(map[string]Asset)
	asset_path := db.root + path
	files, _ := ioutil.ReadDir(asset_path)
	for _, f := range files {
		filedata, err := ioutil.ReadFile(asset_path + "/" + f.Name())
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

func (db *DB) initialize_db() error {
	sqldb, err := sql.Open("sqlite3", db.root + "/static/db/rivulet.db")
	if err != nil {
		fmt.Println("Error reading database.")
		return err
	}
	db.rdb = sqldb
	defer db.rdb.Close()
	return err
}

func NewDatabase(pwd string) (*DB, error) {
	db := &DB{root: pwd}
	err := db.load_assets("/static/assets")
	if err != nil {
		fmt.Println("Error reading assets.")
		return nil, err
	}
	/*
	err = db.initialize_db()
	if err != nil {
		fmt.Println("Error reading database.")
		return nil, err
	}
	*/

/*
	sql := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.rdb.Exec(sql)
*/
	return db, nil
}

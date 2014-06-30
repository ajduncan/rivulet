package rivulet

import (
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
)

type API struct {
	server *RivuletServer
	db     *DB
	http   *martini.ClassicMartini
}

func http_index_handler(w http.ResponseWriter, r *http.Request, db *DB, server *RivuletServer) {
	t, _ := template.ParseFiles(db.root + "/templates/index.html")
	t.Execute(w, db)
}

func (api *API) Init() {
	api.http.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http_index_handler(w, r, api.db, api.server)
	})
}

func (api *API) Run() {
	go http.ListenAndServe(":8066", api.http)
}

func NewAPI(server *RivuletServer, rdb *DB) *API {
	api := &API{server: server, db: rdb, http: martini.Classic()}

	return api
}

package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nabam/koel-server/storage"
	"log"
	"net/http"
)

var db storage.Storage

func ping(w rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.(http.ResponseWriter).Write([]byte("ok"))
}

func services(w rest.ResponseWriter, req *rest.Request) {
	svcs, err := db.GetServices()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(svcs)
}

func createService(w rest.ResponseWriter, req *rest.Request) {
	if err := db.CreateService(req.PathParam("service")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).Write([]byte("ok"))
}

func deleteService(w rest.ResponseWriter, req *rest.Request) {
	if err := db.DeleteService(req.PathParam("service")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).Write([]byte("ok"))
}

func main() {
	configFileName := flag.String("config", "koel-server.conf", "path to config")
	flag.Parse()

	var appConfig config
	if _, err := toml.DecodeFile(*configFileName, &appConfig); err != nil {
		fmt.Println(fmt.Sprintf("Can't parse config: %s", err))
		return
	}

	var err error
	db, err = storage.InitSqlite(appConfig.Sqlite.Path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not initialize sqlite %s\n", err))
	}

	defer db.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/v1/ping", ping),
		rest.Get("/v1/services", services),
		rest.Put("/v1/services/#service", createService),
		rest.Delete("/v1/services/#service", deleteService),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", appConfig.Http.Host, appConfig.Http.Port), api.MakeHandler()))
}

package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	. "./profiler"
	"compress/gzip"
	"time"
)

var rawStorage = LocalRawStorage{}
var process = LocalProcessStorage{}
var dbh = DBHelper{}
var config = Config{}

func init() {
	config.Init()

	dbh.Server = config.MongoHost
	dbh.Database = config.DBName
	dbh.Connect()

	go ProccessLoad(&rawStorage, &process, &dbh)
	ticker := time.NewTicker(time.Duration(config.LoadInterval) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				//Call the periodic function here.
				ProccessLoad(&rawStorage, &process, &dbh)
			}
		}
	}()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/filter", Index).Methods("GET")
	r.HandleFunc("/table", Table).Methods("POST")
	r.HandleFunc("/reduce", Reduce).Methods("POST")
	//s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	//r.PathPrefix("/static/").Handler(s)

	r.PathPrefix("/static").Handler(http.FileServer(http.Dir("./static/")))
	//r.PathPrefix("/static/").Handler(
	//	http.StripPrefix("/static/", http.FileServer(http.Dir("/static/static"))))
	//r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "./static/index.html")
	//})

	host := "192.168.2.102:8002"
	log.Println(fmt.Sprintf("[x] started server on address %s", host))
	if err := http.ListenAndServe(host, handlers.CompressHandlerLevel(r, gzip.BestCompression)); err != nil {
		log.Fatal(err)
	}

	//http.Handle("/", r)
}

func Index(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path[1:]
	http.ServeFile(w, r, "./static/index.html")
}

func Table(w http.ResponseWriter, r *http.Request) {
	s1 := time.Now()

	if len(rawStorage.Items) == 0 {
		RespondWithJson(w, http.StatusOK, M{"code": 2, "response": nil})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var raw HandlerTableParams
	err := decoder.Decode(&raw)
	if err != nil {
		panic(err)
	}

	var filtered []ProfileRecordView
	for _, view := range process.Data {
		if raw.CheckTime(view.Time) {
			filtered = append(filtered, view)
		}
	}
	s2 := time.Now()
	log.Println("Table: ", s2.Sub(s1))
	RespondWithStream(w, http.StatusOK, M{"code": 0, "response": filtered})
}

func Reduce(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Gorilla!\n"))
	s1 := time.Now()

	if len(rawStorage.Items) == 0 {
		RespondWithJson(w, http.StatusOK, M{"code": 2, "response": nil})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var raw HandlerReduceParams
	err := decoder.Decode(&raw)
	if err != nil {
		panic(err)
	}

	items, values := ReduceRecordsByFields(raw.Fields, &raw, rawStorage.Items, rawStorage.Values)
	views := Process(items, values)
	s2 := time.Now()
	log.Println("Reduce: ", s2.Sub(s1))
	RespondWithStream(w, http.StatusOK, M{"code": 0, "response": views})
}

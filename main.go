package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type dataTemp struct {
	TempSearch []tempSuggested `json:"tempSearches"`
}

type dataBot struct {
	BotSearch []groupBots `json:"botSearches"`
}

func main() {
	r := mux.NewRouter()
	// routes consist of a path and a handler function.
	r.HandleFunc("/xdccTempSearch", xdccTempSearch).Methods("GET")
	r.HandleFunc("/xdccBotSearch", xdccBotSearch).Methods("GET")

	// bind to a port and pass our router in
	http.Handle("/", &middleWareServer{r})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func xdccTempSearch(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form

		fmt.Println("card type", query["cardType"])
		fmt.Println("url value", query["pageLink"])

		t := tempSearchMain()

		d := dataTemp{TempSearch: t}

		tempSearchJSON, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(tempSearchJSON)

		// temp := []byte("this is searchXDCC")
		// w.Write(temp)
	}
}

func xdccBotSearch(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form

		fmt.Println("card type", query["cardType"])
		fmt.Println("url value", query["pageLink"])

		b := botSearchMain()

		d := dataBot{BotSearch: b}

		botSearchJSON, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(botSearchJSON)

		// temp := []byte("this is searchXDCC for bots")
		// w.Write(temp)
	}
}

type middleWareServer struct {
	r *mux.Router
}

func (s *middleWareServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

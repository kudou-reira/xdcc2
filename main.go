package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type dataTemp struct {
	TempSearch []episodesSuggested `json:"tempSearches"`
	TempErrors errorMessage        `json:"tempErrors"`
}

type dataBot struct {
	BotSearch []bots `json:"botSearches"`
}

type dataOptimizedBot struct {
	OptimizedBot [][]uniqueBot `json:"optimizedBots"`
}

type dataMedia struct {
	Media []Media `json:"media"`
}

func main() {
	r := mux.NewRouter()
	// routes consist of a path and a handler function.
	r.HandleFunc("/xdccTempSearch", xdccTempSearch).Methods("GET")
	r.HandleFunc("/xdccBotSearch", xdccBotSearch).Methods("GET")
	r.HandleFunc("/xdccOptimizeDL", xdccOptimizeDL).Methods("GET")
	r.HandleFunc("/xdccAnilist", xdccAnilist).Methods("GET")
	r.HandleFunc("/", xdccRoot).Methods("GET")

	// bind to a port and pass our router in
	http.Handle("/", &middleWareServer{r})

	// log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func xdccAnilist(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form

		tempSeason := query["season"][0]
		tempStartDate := query["year"][0]

		a := anilistMain(tempSeason, tempStartDate)

		m := dataMedia{Media: a}

		mediaJSON, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mediaJSON)

		// temp := []byte("this is xdccAnilist")
		// w.Write(temp)
	}
}

func xdccRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	temp := []byte("this is xdccRoot")
	w.Write(temp)
}

func xdccTempSearch(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form

		tempReceived := query["suggestion"][0]
		fmt.Println("this is the query suggestion", tempReceived)
		// fmt.Sprintf("%T", tempReceived)

		// change tempReceived
		// fake := "hi"
		t, errors := tempSearchMain(tempReceived)

		d := dataTemp{TempSearch: t, TempErrors: errors}

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

		// fmt.Println("this is the botSearch", query)
		tempReceived := query["queue[]"]
		// fmt.Println("this is the stack", tempReceived)

		b := botSearchMain(tempReceived)

		// fmt.Println("this is botSearch not printing", b)

		d := dataBot{BotSearch: b}

		botSearchJSON, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("this is botSearch JSON", botSearchJSON)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(botSearchJSON)

		// temp := []byte("this is searchXDCC for bots")
		// w.Write(temp)
	}
}

//rename OPTIMIZE

func xdccOptimizeDL(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form

		// fmt.Println("this is the optimizeDL", query)

		tempReceived := query["downloads[]"]

		var receivedBots []bots

		for i := range tempReceived {
			var singleReceived bots
			tempReceivedBytes := []byte(tempReceived[i])

			err = json.Unmarshal(tempReceivedBytes, &singleReceived)
			if err != nil {
				panic(err)
			}
			receivedBots = append(receivedBots, singleReceived)
		}

		o := optimizeDLMain(receivedBots)

		// slcT, _ := json.MarshalIndent(o, "", " ")
		// fmt.Println(string(slcT))

		d := dataOptimizedBot{OptimizedBot: o}

		optimizedBotJSON, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(optimizedBotJSON)
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

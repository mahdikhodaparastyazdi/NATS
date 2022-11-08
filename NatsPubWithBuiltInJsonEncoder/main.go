package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

type Order struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)

	if err := ec.Publish("foo", &Order{Id: order.Id, Name: order.Name}); err != nil {
		log.Fatal(err)
	}
	nc.Flush()
	json.NewEncoder(w).Encode(order)
}

var nc, _ = nats.Connect("0.0.0.0:4222")

var ec, _ = nats.NewEncodedConn(nc, nats.JSON_ENCODER)

func main() {
	defer nc.Close()
	defer ec.Close()
	args := flag.Args()
	if len(args) != 2 {
		//showUsageAndExit(1)
	}

	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	opts = append(opts, nats.UserCredentials("aa"))

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", "foo", []byte("Hello World"))
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/order", createOrder).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

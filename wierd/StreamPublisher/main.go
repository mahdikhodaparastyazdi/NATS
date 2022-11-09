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

//var nc, _ = nats.Connect("0.0.0.0:4222")

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)
	//payload, _ := json.Marshal(order)

	//_, err := js.Publish("foo", []byte("hello"))
	if err != nil {
		panic(err)
	}

	//err1 := nc.Publish("foo", payload)
	//if err1 != nil {log.Print("before for")
	//	panic(err1)
	//}
	nc.Flush()
	json.NewEncoder(w).Encode(order)
}

var nc, err = nats.Connect("0.0.0.0:4222")
var js, _ = nc.JetStream(nats.PublishAsyncMaxPending(256))

//if err != nil {
//log.Fatal(err)
//}

func main() {
	js.AddStream(&nats.StreamConfig{
		Name:     "FOO",
		Subjects: []string{"foo"},
		MaxBytes: 1024,
	})
	//js.AddConsumer("ORDERS", &nats.ConsumerConfig{
	//	Durable: "MONITOR",
	//})
	defer nc.Close()
	args := flag.Args()
	if len(args) != 2 {
		//showUsageAndExit(1)
	}

	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	opts = append(opts, nats.UserCredentials("aa"))

	//subj, msg := args[0], []byte(args[1])
	//nc.Publish(subj, msg)
	//nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", "foo", []byte("Hello World"))
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/order", createOrder).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/codingsince1985/geo-golang/openstreetmap"
)

type Reverse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	DisplayName string `json:"displayname"`
}

func reverse(w http.ResponseWriter, req *http.Request) {
	var t Reverse
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		io.WriteString(w, `{"error": "malformed json?"}`)
		return
	}

	address, _ := openstreetmap.Geocoder().ReverseGeocode(t.Lat, t.Lng)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(Address{address})
	if err != nil {
		panic("can't marshal our own stuff?")
	}

	fmt.Printf("Address of (%f,%f) is %s\n", t.Lat, t.Lng, address)
}

func main() {
	http.HandleFunc("/reverse", reverse)
	http.Handle("/", http.FileServer(http.Dir("root")))
	log.Fatal(http.ListenAndServe("localhost:1337", nil)) // nil for DefaultServeMux
}

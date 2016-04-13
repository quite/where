package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/codingsince1985/geo-golang/openstreetmap"
)

const lat, lng = 55.6195815, 12.9779243

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

func parseFloat(value interface{}) float64 {
	f, _ := strconv.ParseFloat(value.(string), 64)
	return f
}

func main() {
	http.HandleFunc("/reverse", reverse)
	http.Handle("/", http.FileServer(http.Dir("root")))
	log.Fatal(http.ListenAndServe("localhost:1337", nil)) // nil for DefaultServeMux
}

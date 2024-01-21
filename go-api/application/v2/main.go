package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Payload struct {
	IP         string `json:"PodIP"`
	PodName    string `json:"PodName"`
	Path       string `json:"Path"`
	Time       string `json:"Time"`
	Method     string `json:"HTTPMethod"`
	APIVersion string `json:"APIVersion"`
}

func main() {
	http.HandleFunc("/", handleAll)
	http.ListenAndServe(":8080", nil)
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	log.Println("PATH : ", r.RequestURI, " | METHOD : ", r.Method, " | SOURCE : ", r.RemoteAddr)
	t := time.Now()
	p := Payload{
		IP:         os.Getenv("POD_IP"),
		PodName:    os.Getenv("POD_NAME"),
		Path:       r.RequestURI,
		Time:       fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
		Method:     r.Method,
		APIVersion: "v2",
	}
	b, e := json.Marshal(p)
	if e != nil {
		fmt.Println("ERROR : ", e)
	}
	fmt.Fprintf(w, string(b))
}

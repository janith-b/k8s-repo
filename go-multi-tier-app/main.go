package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var externalPath string
var externalDomain string = "http://" + os.Getenv("EXTERNAL_DOMAIN")

func main() {
	// fmt.Print(os.Args)
	if os.Args[1] == "app01" {
		externalPath = "app02"
	} else {
		externalPath = "app01"
	}

	r := http.NewServeMux()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy: " + os.Args[1] + "\n"))
	})
	r.HandleFunc("/"+os.Args[1], func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello From " + strings.ToUpper(os.Args[1]) + "\n"))
	})
	r.HandleFunc("/"+externalPath, func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(r.Method, externalDomain+"/"+externalPath, nil)
		if err != nil {
			fmt.Println("ERROR : ", err)
		}
		client := http.Client{}
		resp, e := client.Do(req)
		if e != nil {
			fmt.Println(e)
		}
		b := bytes.Buffer{}
		io.Copy(&b, resp.Body)
		time.Sleep(1000)

		w.Write([]byte(resp.Status + " | " + string(b.Bytes())))

	})
	http.ListenAndServe(":"+os.Args[2], r)

}

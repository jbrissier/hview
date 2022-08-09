package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HViewResult struct {
	Time    string
	Headers map[string]string
	IP      string
}

func (h *HViewResult) Writer() (int, error) {
	return 0, nil
}

func jsonHeaders(w http.ResponseWriter, req *http.Request) {
	// return the same as headers but in json format

	w.Header().Set("Content-Type", "application/json")
	headerMap := make(map[string]string)

	for name, headers := range req.Header {
		for _, h := range headers {
			headerMap[name] = h
		}
	}

	hv := HViewResult{
		Time:    time.Now().Format(time.RFC3339),
		IP:      req.RemoteAddr,
		Headers: headerMap,
	}

	res, _ := json.Marshal(hv)

	w.Write(res)

}

func textHeaders(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hview\n\n")
	fmt.Fprintf(w, "Time:\n-------------------\n")
	fmt.Fprintf(w, "%s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "Header:\n-------------------\n")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "\nIP:\n-------------------\n")
	fmt.Fprintf(w, "%s\n", req.RemoteAddr)
}

func headers(w http.ResponseWriter, req *http.Request) {
	// This handler does something a little more sophisticated by reading all the HTTP request headers and echoing them into the response body.

	json := req.URL.Query().Get("format")

	if json == "json" {
		jsonHeaders(w, req)
		return
	}
	// return as text
	textHeaders(w, req)
}

func main() {

	port := os.Getenv("H_VIEW_PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", headers)
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

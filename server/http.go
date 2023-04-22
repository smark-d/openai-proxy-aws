package server

import (
	"fmt"
	"github.com/smark-d/openai-proxy-aws/server/api"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var target = "https://api.openai.com"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if !api.Filter(r) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	_, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("Error parsing URL: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	targetURL := target + r.URL.Path
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Println("Error creating proxy request: ", err.Error())
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	for headerKey, headerValues := range r.Header {
		for _, headerValue := range headerValues {
			proxyReq.Header.Add(headerKey, headerValue)
		}
	}
	// Replace the Authorization header with Turing's API key
	ikey, err := api.GetOpenAIkey()
	if err != nil {
		log.Println("Error getting openai key: ", err.Error())
		http.Error(w, "Error getting openai key", http.StatusInternalServerError)
		return
	}
	proxyReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ikey))
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Println("Error sending proxy request: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.CopyBuffer(w, resp.Body, make([]byte, 256))
}

func Start(port string) {
	http.HandleFunc("/", handlerFunc)
	log.Printf("Starting server on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

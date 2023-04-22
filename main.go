package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

func notFound(response http.ResponseWriter, request *http.Request) {
	result := `{"status": 404, "message": "404 NOT FOUND"}`

	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)

	json.NewEncoder(response).Encode(finalResult)
}

func main() {
	fmt.Println("Starting ByteTrack API...")

	route := mux.NewRouter()
	route.Use(commonMiddleware)

	router := cors.Default().Handler(route)

	route.NotFoundHandler = http.HandlerFunc(notFound)

	http.ListenAndServe(":31475", router)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fmt.Print(request.URL.Path)

		var globallimiter = rate.NewLimiter(50, 110)

		if !globallimiter.Allow() {
			ratelimited(response, request)
			return
		}

		response.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(response, request)
	})
}

func ratelimited(response http.ResponseWriter, request *http.Request) {
	result := `{"status":429, "message":"You are requesting too quickly!"}`
	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)
	json.NewEncoder(response).Encode(finalResult)
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusTooManyRequests)
}

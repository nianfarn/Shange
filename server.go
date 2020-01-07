package main

import (
	"encoding/json"
	"fmt"
	"log"
	. "net/http"
)

func userHandler(w ResponseWriter, r *Request) {
	switch r.Method {
	case MethodPost:
		bodyValidation(r, w)

		var u user
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			responseError(w, StatusBadRequest, "unprocessable body")
		}
		defer r.Body.Close()

		user := createUser(&u)
		responseSuccess(w, StatusOK, user)
	}
}

func bodyValidation(r *Request, w ResponseWriter) {
	value := r.Header.Get("Content-Type")
	if value != "application/json" {
		Error(w, "Content-Type header is not application/json", StatusUnsupportedMediaType)
	}
	if r.Body == nil {
		responseError(w, StatusBadRequest, "empty body")
	}
}

func entryHandler(w ResponseWriter, r *Request) {
	log.Print("from log")   // fixme: remove
	fmt.Print("from print") // fixme: remove
	w.Write([]byte("<h1>Hello there, i'm started listening</h1>"))
}

// Makes the response with payload in JSON
func responseSuccess(w ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Makes the error response with payload in JSON
func responseError(w ResponseWriter, code int, message string) {
	responseSuccess(w, code, map[string]string{"error": message})
}

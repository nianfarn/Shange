package main

import (
	"encoding/json"
	. "net/http"
)

type httpError struct {
	code    int
	message string
}

func (e *httpError) response(w ResponseWriter) {
	response(w, e.code, map[string]string{"error": e.message})
}

func userHandler(w ResponseWriter, r *Request) {
	//vars := mux.Vars(r)

	//id := vars["id"]

	switch r.Method {
	case MethodPost:
		if err := validateCreateRequest(r); err != nil {
			err.response(w)
		}

		var u user
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			e := httpError{code: StatusBadRequest, message: "unprocessable body"}
			e.response(w)
		}
		defer r.Body.Close()

		user := createUser(&u)

		response(w, StatusOK, user)
	case MethodGet:
		if err := validateGetRequest(r); err != nil {
			err.response(w)
		}

		//user := findUser(r.RequestURI)
	}
}

func validateGetRequest(r *Request) *httpError {
	if r.Body != nil {
		return &httpError{
			code:    StatusBadRequest,
			message: "Unexpected body",
		}
	}
	return nil
}

func validateCreateRequest(r *Request) *httpError {
	if r.Body == nil {
		return &httpError{
			code:    StatusBadRequest,
			message: "Empty body",
		}
	}

	value := r.Header.Get("Content-Type")
	if value != "application/json" {
		return &httpError{
			code:    StatusUnsupportedMediaType,
			message: "Content-Type header is not application/json",
		}
	}

	return nil
}

// Makes the response with payload in JSON
func response(w ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if status == 0 {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

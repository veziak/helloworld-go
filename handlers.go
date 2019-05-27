package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var db = NewDB()

type DateOfBirthRequest struct {
	DateOfBirth string `json:"dateOfBirth"`
}

func getBirthdayMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	log.Printf("get username %v", username)

	exists, err := db.UserExist(username)
	if err != nil {
		log.Printf("Error checking if user exist, error: %v", err)
		errorResponse(w, 500, "Internal error.")
		return
	}
	if !exists {
		errorResponse(w, 404, fmt.Sprintf("User %s not found.", username))
		return
	}

	user, err := db.GetUser(username)
	if err != nil {
		errorResponse(w, 500, "error")
		return
	}

	// all time should be in UTC, otherwise time difference will calculated incorrectly
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	message := GetBirthdayMessage(username, user.DateOfBirth, now)

	result := map[string]string{"message": message}
	Response(w, http.StatusOK, result)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	err := db.CheckDBConnection()
	if err != nil {
		log.Printf("Healtheck failed, error: %v", err)
		errorResponse(w, http.StatusInternalServerError, "Service is not available at the moment.")
		return
	}
	Response(w, http.StatusOK, "")
}

func version(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "na"
	}
	Response(w, http.StatusOK, version)
}

func createOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	err := ValidateUsername(username)
	if err != nil {
		log.Printf("Wrong username format %v", username)
		errorResponse(w, http.StatusBadRequest, "Can't read request body.")
		return
	}

	data := DateOfBirthRequest{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Can't read request body, Error: %v", err)
		errorResponse(w, http.StatusBadRequest, "Can't read request body.")
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("can't parse json, Error: %v , body: %v", err, string(body))
		errorResponse(w, http.StatusBadRequest, "Can't parse json.")
		return
	}

	exist, err := db.UserExist(username)
	if err != nil {
		log.Printf("Can't check id user exist, error: %v", err)
		errorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	if exist {
		// update user
		user, err := db.GetUser(username)
		if err != nil {
			log.Printf("Error getting a user from database, error: %v", err)
			errorResponse(w, http.StatusInternalServerError, "internal error")
			return
		}
		dob, err := ValidateDateOfBirth(data.DateOfBirth)
		if err != nil {
			log.Printf("Can't validate dateOfBirth: %v", err)
			errorResponse(w, http.StatusBadRequest, "Can't validate dateOfBirth.")
			return
		}
		user.DateOfBirth = dob

		err = db.UpdateUser(user)
		if err != nil {
			log.Printf("Error updating a user to database, error: %v", err)
			errorResponse(w, http.StatusInternalServerError, "internal error")
			return
		}

		Response(w, http.StatusNoContent, "")
	} else {
		// create a new user
		dob, err := ValidateDateOfBirth(data.DateOfBirth)
		if err != nil {
			log.Printf("Can't validate dateOfBirth: %v", err)
			errorResponse(w, http.StatusBadRequest, "Can't validate dateOfBirth.")
			return
		}

		_, err = db.CreateUser(username, dob)
		if err != nil {
			log.Printf("Error saving a new user to database, error: %v", err)
			errorResponse(w, http.StatusInternalServerError, "internal error")
			return
		}
		Response(w, http.StatusNoContent, "")
	}
}

//  write json response
func Response(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code != http.StatusNoContent {
		_, err := w.Write(response)
		if err != nil {
			log.Printf("Error writing to response. error: %v", err)
		}
	}
}

// return error message
func errorResponse(w http.ResponseWriter, code int, msg string) {
	Response(w, code, map[string]string{"message": msg})
}

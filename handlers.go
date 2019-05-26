package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"
)

var db = NewDB()

type DateOfBirthRequest struct {
	DateOfBirth string `json:"dateOfBirth"`
}

func ValidateUsername(username string) error {
	if username == "" {
		log.Print("ValidateUsername failed, username is empty")
		return fmt.Errorf("username can't be empty.")
	}
	match, _ := regexp.MatchString("^[a-zA-Z]+$", username)
	if match {
		log.Printf("ValidateUsername successful, username: %s", username)
		return nil
	} else {
		log.Printf("ValidateUsername failed, username: %s", username)
		return fmt.Errorf("username is in a wrong format.")
	}
}

func ValidateDateOfBirth(dateOfBirth string) (time.Time, error) {
	match, _ := regexp.MatchString("\\d{4}-\\d{2}-\\d{2}", dateOfBirth)
	if match {
		t, err := time.Parse("2006-01-02", dateOfBirth)
		if err == nil {
			// date must be at least one day before current date
			if t.Unix() < time.Now().Add(-24*time.Hour).Unix() {
				log.Printf("ValidateDateOfBirth successful, dateOfBirth: %s", dateOfBirth)
				return t, nil
			} else {
				log.Printf("ValidateDateOfBirth failed, dateOfBirth is in the future: %s", dateOfBirth)
			}

		} else {
			log.Printf("ValidateDateOfBirth failed, dateOfBirth: %s error: %v", dateOfBirth, err)
		}
	} else {
		log.Printf("ValidateDateOfBirth regexp validation failed, dateOfBirth: %s", dateOfBirth)
	}
	return time.Time{}, fmt.Errorf("dateOfBirth: %s is in a wrong format", dateOfBirth)
}

func GetBirthdayMessage(username string, birthDate time.Time, today time.Time) string {

	log.Printf("b: %v, t: %v", birthDate, today)
	if today.Month() == birthDate.Month() && today.Day() == birthDate.Day() {
		return fmt.Sprintf("Hello, %s! Happy birthday!", username)
	}
	nextBirthday := time.Date(today.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.UTC)

	if nextBirthday.Unix() < today.Unix() {
		nextBirthday = nextBirthday.AddDate(1, 0, 0)
	}

	log.Printf("hours left: %v", nextBirthday.Sub(today).Hours())
	// even if 1 hour left we should return "1 day left"
	days := int64(math.Ceil(nextBirthday.Sub(today).Hours() / 24))
	message := fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, days)
	log.Print(message)
	return message
}

func getBirthdayMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	log.Printf("get username %v", username)

	exists, err := db.UserExist(username)
	if err != nil {
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

		Response(w, http.StatusNoContent, "aa")
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

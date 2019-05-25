package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestUsernameValidation(t *testing.T) {
	err := ValidateUsername("asasas")
	if err != nil {
		t.Error(err)
	}

	err = ValidateUsername("asasas2232223")
	if err == nil {
		t.Error("Only characters should be allowed in username.")
	}

	err = ValidateUsername("asasas2 232223")
	if err == nil {
		t.Error("Only characters should be allowed in username.")
	}

	err = ValidateUsername("")
	if err == nil {
		t.Error("Username can't be empty.")
	}

	err = ValidateUsername(" ssdf  ")
	if err == nil {
		t.Error("No whitespace allowed in username.")
	}

	err = ValidateUsername("გამარჯობა")
	if err == nil {
		t.Error("Only latin characters are allowed in username.")
	}
}


func TestDateOfBirthValidation(t *testing.T) {
	_, err := ValidateDateOfBirth("asasas")
	if err == nil {
		t.Error("Wrong date format.")
	}

	_, err = ValidateDateOfBirth("2012-31-01")
	if err == nil {
		t.Error("Wrong date format.")
	}

	_, err = ValidateDateOfBirth("01/02/2013")
	if err == nil {
		t.Error("Wrong date format.")
	}

	_, err = ValidateDateOfBirth("")
	if err == nil {
		t.Error("dateOfBirth can't be empty.")
	}

	_, err = ValidateDateOfBirth("2029-01-01")
	if err == nil {
		t.Error("Date should be in the past.")
	}

	now := time.Now().Truncate(24 * time.Hour)
	_, err = ValidateDateOfBirth(fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day()))
	if err == nil {
		t.Error("Date should be in the past.")
	}

	_, err = ValidateDateOfBirth("2000-01-40")
	if err == nil {
		t.Error("Wrong date format.")
	}

	_, err = ValidateDateOfBirth("2012-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestBirthMessage(t *testing.T) {
	username := "test"
	today := time.Date(2016, 5, 26, 0, 0, 0, 0, time.UTC)

	birthDate := time.Date(2010, 5, 26, 0, 0, 0, 0, time.UTC)
	message := GetBirthdayMessage(username, birthDate, today)
	if message != "Hello, test! Happy birthday!" {
		t.Error("message is incorrect")
	}
	log.Print(message)

	birthDate = time.Date(2010, 5, 27, 0, 0, 0, 0, time.UTC)
	message = GetBirthdayMessage(username, birthDate, today)
	if message != "Hello, test! Your birthday is in 1 day(s)" {
		t.Error("message is incorrect")
	}
	log.Print(message)

	birthDate = time.Date(2010, 6, 26, 0, 0, 0, 0, time.UTC)
	message = GetBirthdayMessage(username, birthDate, today)
	if message != "Hello, test! Your birthday is in 31 day(s)" {
		t.Error("message is incorrect")
	}
	log.Print(message)

	// leap year
	today = time.Date(2016, 2, 2, 0, 0, 0, 0, time.UTC)
	birthDate = time.Date(2010, 2, 1, 0, 0, 0, 0, time.UTC)
	message = GetBirthdayMessage(username, birthDate, today)
	if message != "Hello, test! Your birthday is in 365 day(s)" {
		t.Error("message is incorrect")
	}
	log.Print(message)

	// normal year
	today = time.Date(2017, 2, 2, 0, 0, 0, 0, time.UTC)
	birthDate = time.Date(2010, 2, 1, 0, 0, 0, 0, time.UTC)
	message = GetBirthdayMessage(username, birthDate, today)
	if message != "Hello, test! Your birthday is in 364 day(s)" {
		t.Error("message is incorrect")
	}
	log.Print(message)
}
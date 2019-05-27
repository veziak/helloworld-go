package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"time"
)

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

	// even if 1 hour left we should return "1 day left"
	days := int64(math.Ceil(nextBirthday.Sub(today).Hours() / 24))
	message := fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, days)
	log.Print(message)
	return message
}

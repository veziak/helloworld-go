package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"os"
	"time"
)

type User struct {
	Username    string    `sql:"username,pk"`
	DateOfBirth time.Time `sql:"dateofbirth"`
}

type DB struct {
	DB *pg.DB
}

func NewDB() *DB {
	defaultHost := "localhost:5432"
	defaultUser := "postgres"
	defaultPassword := "postgres"
	host := os.Getenv("POSTGRES_DB_HOST")
	user := os.Getenv("POSTGRES_DB_USER")
	password := os.Getenv("POSTGRES_DB_PASSWORD")

	if host == "" {
		host = defaultHost
	}
	if user == "" {
		user = defaultUser
	}
	if password == "" {
		password = defaultPassword
	}

	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: "hello",
		Addr:     host})

	return &DB{db}
}

// save a new user to database
func (db DB) CreateUser(username string, dateOfBirth time.Time) (*User, error) {

	user := User{username, dateOfBirth}
	err := db.DB.Insert(&user)
	if err != nil {
		return nil, fmt.Errorf("Could not save user into the database: %s", err)
	}
	return &user, nil
}

// update user date of birth in database
func (db DB) UpdateUser(user *User) error {
	res, err := db.DB.Model(user).Column("dateofbirth").WherePK().Update()
	if err != nil {
		return fmt.Errorf("Could not update %s, error: %v", user.Username, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("Could not update user, %s not found.", user.Username)
	}
	return nil
}

// get user by username
func (db DB) GetUser(username string) (*User, error) {
	var user User
	err := db.DB.Model(&user).Where("username = ?", username).First()
	if err != nil {
		return nil, fmt.Errorf("Could not get user from the database: %s", err)
	}
	return &user, nil
}

// check if user exist
func (db DB) UserExist(username string) (bool, error) {
	var user User
	exist, err := db.DB.Model(&user).Where("username = ?", username).Exists()
	if err != nil {
		return true, fmt.Errorf("Could not get user from the database: %s", err)
	}
	return exist, nil
}

// check database connection
func (db DB) CheckDBConnection() error {
	_, err := db.DB.Exec("select 1")
	return err
}

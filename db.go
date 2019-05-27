package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
	host := os.Getenv("POSTGRES_DB_HOST")
	user := os.Getenv("POSTGRES_DB_USER")
	password := os.Getenv("POSTGRES_DB_PASSWORD")

	if host == "" {
		host = "172.17.0.2:5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}

	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: "hello",
		Addr:     host})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	return &DB{db}
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
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

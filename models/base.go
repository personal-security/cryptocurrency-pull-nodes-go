package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	type ConfigDB struct {
		User     string `json:"user"`
		Pass     string `json:"password"`
		Database string `json:"db"`
	}

	var config ConfigDB

	file, _ := ioutil.ReadFile("config/db.json")

	json.Unmarshal([]byte(file), &config)

	username := config.User   // os.Getenv("db_user")
	password := config.Pass   // os.Getenv("db_pass")
	dbName := config.Database // os.Getenv("db_name")
	dbHost := "localhost"     //os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	//fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	db = conn
	db.Debug().AutoMigrate(
		&CallbackWallet{},
	)
}

func GetDB() *gorm.DB {
	return db
}

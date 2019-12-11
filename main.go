package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config ...
type Config struct {
	DBName string `json:"db_name"`
	DBHost string `json:"db_host"`
	DBPass string `json:"db_pass"`
	DBPort string `json:"db_port"`
	DBUser string `json:"db_user"`
}

var conf Config
var dsn string

func check(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func readTextFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func init() {
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	check("error", err)

	byteValue, err := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &conf)

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)
}

func main() {
	exportFileName := "results.csv"
	args := os.Args

	if len(args) <= 1 {
		log.Fatalln("SQL source is required")
	}

	if len(args) == 3 {
		exportFileName = args[2]
	}

	sqlFile := args[1]
	sqlCommand, err := readTextFile(sqlFile)
	check("error", err)

	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	check("error", err)

	rows, err := db.Query(sqlCommand)
	defer rows.Close()
	check("error", err)

	csvFile, err := os.Create(exportFileName)
	defer csvFile.Close()
	check("error", err)

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	columns, err := rows.Columns()
	check("error", err)

	totalColumn := len(columns)
	rawEntry := make([][]byte, totalColumn)
	entry := make([]interface{}, totalColumn)
	entryResult := make([]string, totalColumn)
	for index := range rawEntry {
		entry[index] = &rawEntry[index]
	}

	writer.WriteAll([][]string{columns})
	for rows.Next() {
		err = rows.Scan(entry...)
		check("error", err)
		for i, raw := range rawEntry {
			entryResult[i] = string(raw)
		}
		err := writer.Write(entryResult)
		check("error", err)
	}
	log.Print("Done...")
}

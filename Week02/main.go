package main

import "github.com/pkg/errors"
import _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"
import "log"

var db *sql.DB

func main() {
	var initDbErr error
	db, initDbErr = sql.Open("mysql", "root:123@/mydb?charset=utf8")
	if initDbErr != nil {
		log.Fatal(initDbErr)
	}
	fmt.Printf("%+v", service())
}

func dao(sqlStr string) error {
	rows, err := db.Query(sqlStr)
	if err != nil {
		return errors.Wrap(err, "product data not found!!")
	}
	fmt.Println(rows)
	return nil
}

func service() error {
	return dao("select id,name,price from product")
}

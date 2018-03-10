package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("mysql", "star:starstar@tcp(116.196.123.49:3306)/star?charset=utf8")
	if err != nil {
		fmt.Println("出现错误")
		fmt.Println(err)
	}
	db.SetMaxOpenConns(30)
}

func getChatID(id string) string {

	fmt.Println("id:"+id)
	query, err := db.Query("SELECT chatID FROM user WHERE id = " + id)
	if err != nil {
		fmt.Println("查询出现错误")

		fmt.Println(err)
	}
	var chatID string

	for query.Next() {
		if err := query.Scan(&chatID); err != nil {
			return ""
		}
	}
	return chatID
}

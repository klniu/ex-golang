package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func main() {
	start := time.Now()
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db1?charset=utf8")
	if err != nil {
		fmt.Printf("sql.Open: %v\r\n", err)
		return
	}
	defer db.Close()

	userList, err := getUserList()
	if err != nil {
		fmt.Printf("getUserList: %v\r\n", err)
		return
	}

	for _, item := range userList {
		var transactionCount int
		row := db.QueryRow("SELECT COUNT(*) FROM `mysql`.`user` WHERE set_id = ?", item)
		if err := row.Scan(&transactionCount); err != nil {
			fmt.Printf("row.Scan: %v\r\n", err)
			return
		}

		if transactionCount <= 0 {
			db.Exec("DELETE FROM `mysql`.`user` WHERE set_id = ?", item)
			db.Exec("DELETE FROM `mysql`.`user` WHERE set_id = ?", item)
			fmt.Println(item)
		}
	}

	end := time.Now()
	fmt.Println(start.Format("2006-01-02 15:04:05") + " ~ " + end.Format("2006-01-02 15:04:05"))
}

func getUserList() ([]int, error) {
	rows, err := db.Query("SELECT `user_id` FROM `user`")
	if err != nil {
		fmt.Printf("db.Query: %v\r\n", err)
		return nil, err
	}
	defer rows.Close()

	var userList []int
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			fmt.Printf("rows.Scan: %v\r\n", err)
			return nil, err
		}
		userList = append(userList, userId)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("rows.Err: %v\r\n", err)
		return nil, err
	}

	return userList, nil
}

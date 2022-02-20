package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/tenntenn/sqlite"
)

type Record struct {
	ID    int64
	Name  string
	Phone string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func run() error {
	var input string
	db, err := sql.Open(sqlite.DriverName, "addressbook.db")
	if err != nil {
		return err
	}

	if err := createTable(db); err != nil {
		return err
	}

	for {
		if err := showRecords(db); err != nil {
			return err
		}

		fmt.Print("Change Record? (Y/n)")
		fmt.Scan(&input)

		if input == "y" {
			if err := editRecord(db); err != nil {
				return err
			}
		}

		fmt.Print("Add Record? (Y/n)")
		fmt.Scan(&input)
		if input == "y" {
			if err := inputRecord(db); err != nil {
				return err
			}
		}
	}

	return nil
}

func createTable(db *sql.DB) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS addressbook (
			id    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name  TEXT NOT NULL,
			phone TEXT NOT NULL
	);`

	if _, err := db.Exec(sql); err != nil {
		return err
	}

	return nil
}

func showRecords(db *sql.DB) error {
	fmt.Println("全件表示")
	rows, err := db.Query("SELECT * FROM addressbook")
	if err != nil {
		return err
	}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Name, &r.Phone); err != nil {
			return err
		}
		fmt.Printf("[%d]  Name:%s TEL:%s\n", r.ID, r.Name, r.Phone)
	}
	fmt.Println("--------")

	return nil
}

func editRecord(db *sql.DB) error {
	var input string
	var r Record
	fmt.Print("ID >")
	fmt.Scan(&r.ID)
	// IDに紐ずくレコードを取得
	rows := db.QueryRow("SELECT * FROM addressbook WHERE ID = ?", r.ID)
	if err := rows.Scan(&r.ID, &r.Name, &r.Phone); err != nil {
		return err
	}
	fmt.Printf("[%d]  Name:%s TEL:%s\n", r.ID, r.Name, r.Phone)

	fmt.Print("Which Column Change? (name / TEL)")
	fmt.Scan(&input)

	switch input {
	case "name":
		fmt.Print("Enter the changed name ? >")
		fmt.Scan(&r.Name)
		r, err := db.Exec("UPDATE addressbook SET name = ? WHERE ID = ? ", r.Name, r.ID)
		if err != nil {
			log.Println(err)
		}
		cnt, err := r.RowsAffected()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Affected rows:", cnt)
	case "TEL":
		fmt.Print("Enter the changed TEL ? >")
		fmt.Scan(&r.Phone)
		r, err := db.Exec("UPDATE addressbook SET TEL = ? WHERE ID = ? ", r.Phone, r.ID)
		if err != nil {
			log.Println(err)
		}
		cnt, err := r.RowsAffected()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Affected rows:", cnt)
	default:
		fmt.Print("Sorry, Invalid command")
		return nil
	}
	return nil
}

func inputRecord(db *sql.DB) error {
	var r Record

	fmt.Print("Name >")
	fmt.Scan(&r.Name)

	fmt.Print("TEL >")
	fmt.Scan(&r.Phone)
	const sql = "INSERT INTO addressbook(name, phone) values (?,?)"
	_, err := db.Exec(sql, r.Name, r.Phone)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func createdb(filepath string) {
	var err error
	app.DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Can not allocate database on disk")
		panic(err)
	}
	if app.DB == nil {
		panic("Failed to create database")
	}
}

func createTable() {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS vehicles(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Uid TEXT NOT NULL,
		State TEXT,
		Battery INTEGER,
		CreatedAt TEXT ,
		UpdatedAt TEXT
	);
	`

	_, err := app.DB.Exec(sql_table)
	if err != nil {
		log.Fatal("Can not create table")
		panic(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM vehicles")
	app.DB.Exec("ALTER SEQUENCE vehicles_id_seq RESTART WITH 1")
}

func SaveVehicle(v Vehicle) {
	statement, err := app.DB.Prepare("INSERT INTO vehicles (Uid,State, Battery,CreatedAt,UpdatedAt) VALUES (?,?,?,?,?)")
	if (err != nil) {
		log.Print(err)
		return
	} else {
		_, err := statement.Exec(v.Uid,v.inState(), v.ChargeLevel(), v.CreatedAt, v.UpdatedAt)
		if (err != nil) {
			panic("could not save vehicle to database...")
		}
	}
}



func RestoreVehiclesFromDB() {
	rows, err := app.DB.Query("SELECT * FROM vehicles")
	if (err != nil) {
		exit(err)
	}
	var Id int
	var Uid string
	var State uint8
	var Battery int
	var createdAt string
	var updatedAt string


	for rows.Next() {
		err = rows.Scan(&Id, &Uid, &State, &Battery,&createdAt,&updatedAt)
		if (err != nil) {
			exit(err)
		}
		fmt.Println(Id)
		fmt.Println(Uid)
		fmt.Println(State)
		fmt.Println(Battery)
		fmt.Println(createdAt)
		fmt.Println(updatedAt)

	}

	rows.Close() //good habit to close
}

func init() {
	fmt.Print("Configuring the database...\n")
	createdb("./v.db")
	createTable()
//	RestoreVehiclesFromDB()
	go func() {
		<-app.start
		for {
			select {
			case v := <-app.store:
				fmt.Print("Presisting Vehicle id: ",v.Uid + "\n")
				SaveVehicle(*v);
			case <-app.quit:
				break
			}
		}
	}()
}

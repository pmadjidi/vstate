package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"vehicles/vstate"
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
		Uid TEXT NOT NULL UNIQUE,
		State INTEGER ,
		Battery INTEGER,
		CreatedAt INTEGER ,
		UpdatedAt INTEGER 
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

func SaveVehicle(v *Vehicle) {

	uid,state,battery,createdAt,updatedAt := v.readValues();
	statement, err := app.DB.Prepare("INSERT OR REPLACE INTO vehicles (Uid,State, Battery,CreatedAt,UpdatedAt) VALUES (?,?,?,?,?)")
	if (err != nil) {
		log.Print(err.Error())
		return
	} else {
		_, err := statement.Exec(uid,state,battery,createdAt,updatedAt)
		if (err != nil) {
			log.Print(err.Error())
			panic("could not save vehicle to database...")
		}
	}
}

func DeleteVehicle(v *Vehicle) {
	statement, err := app.DB.Prepare("DELETE  FROM vehicles where Uid = (Uid) VALUES (?)")
	if (err != nil) {
		log.Print(".............." + err.Error())
		return
	} else {
		_, err := statement.Exec(v.getUid())
		if (err != nil) {
			panic("could not delete the vehicle from database...")
		}
	}
}

func RestoreVehiclesFromDB() {
	rows, err := app.DB.Query("SELECT * FROM vehicles")
	if (err != nil) {
		exit(err)
	}
	var id int
	var uid string
	var state vstate.State
	var battery int
	var createdAt int64
	var updatedAt int64

	for rows.Next() {
		err = rows.Scan(&id, &uid, &state, &battery, &createdAt, &updatedAt)
		if (err != nil) {
			exit(err)
		}
		fmt.Println(id)
		fmt.Println(uid)
		fmt.Println(state)
		fmt.Println(battery)
		fmt.Println(createdAt)
		fmt.Println(updatedAt)
		v := Vehicle{Id: id, Uid: uid, State: state, Battery: battery, CreatedAt: createdAt, UpdatedAt: updatedAt, Port: make(chan *Request)}
		app.garage.Set(v.Uid, &v)
		v.listen()
	}

	rows.Close()
}

func init() {
	fmt.Print("Configuring the database...\n")
	createdb("./v.db")
	createTable()
	RestoreVehiclesFromDB()
	go func() {
		<-app.start
	Loop:
		for {
			select {
			case v := <-app.store:
				fmt.Print("\n Presisting Vehicle id: ", v.Uid+"\n")
				SaveVehicle(v);
			case v := <-app.delete:
				fmt.Print("\n Deleting Vehicle id: ", v.Uid+"\n")
				DeleteVehicle(v);
				fmt.Print("\n Removing Vehicle id: ", v.Uid+"from Garage....\n")
				app.garage.Delete(v.Uid)
			case <-app.quit:
				fmt.Print("\nStopping database event loop...\n")
				break Loop
			}
		}
	}()
}

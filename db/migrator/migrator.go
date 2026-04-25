package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func performMigration(currentenv map[string]any, direction string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		currentenv["username"], currentenv["password"], currentenv["host"], currentenv["port"], currentenv["database"])

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, errinstance := mysql.WithInstance(db, &mysql.Config{})
	if(errinstance != nil){
		panic(errinstance)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "mysql", driver)
	if err != nil {
		panic(err)
	}

	if direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
	} else {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
	}
	fmt.Printf("Migration %s finished successfully\n", direction)
}

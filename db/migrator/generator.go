package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func generateMigrationFile(currentenv map[string]any, tableName string) {
	if tableName == "" {
		fmt.Println("Error: -table flag is required for creating migrations")
		return
	}

	version := getLatestVersion(currentenv) + 1
	
	upFileName := fmt.Sprintf("%06d_%s.up.sql", version, tableName)
	downFileName := fmt.Sprintf("%06d_%s.down.sql", version, tableName)

	os.Create(filepath.Join(migrationsDir, upFileName))
	os.Create(filepath.Join(migrationsDir, downFileName))

	// Update local tracker
	os.WriteFile(versionFile, []byte(strconv.Itoa(version)), 0644)

	fmt.Printf("Generated: %s and %s\n", upFileName, downFileName)
}

func getLatestVersion(currentenv map[string]any) int {
	// 1. Try local file
	data, err := os.ReadFile(versionFile)
	if err == nil {
		v, _ := strconv.Atoi(strings.TrimSpace(string(data)))
		return v
	}

	// 2. Fallback to Database
	fmt.Println("Local version file missing, checking database...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		currentenv["username"], currentenv["password"], currentenv["host"], currentenv["port"], currentenv["database"])
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return 0
	}
	defer db.Close()

	var version int
	err = db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil {
		return 0
	}
	return version
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type DatabaseConfiguration struct {
	Development map[string]any `yaml:"development"`
	Test        map[string]any `yaml:"test"`
	Production  map[string]any `yaml:"production"`
}

const (
	DEVELOPMENT int = iota + 1
	TEST
	PRODUCTION
)

const migrationsDir = "db/migrator/migrations"
const versionFile = migrationsDir + "/.version"

func main() {
	dbconfig := readDatabaseConfiguration()
	flags := takeInputFlags()
	
	validateFlags(flags)

	action := *flags["action"].(*string)
	envVal := *flags["environment"].(*int)
	currentenv := getAppropriateDbConfig(dbconfig, envVal)

	switch action {
	case "run":
		direction := *flags["direction"].(*string)
		performMigration(currentenv, direction)
	case "create":
		tableName := *flags["table"].(*string)
		generateMigrationFile(currentenv, tableName)
	default:
		fmt.Println("Invalid action. Use -action=run or -action=create")
	}
}

func validateFlags(flags map[string]any) {
	env := *flags["environment"].(*int)
	direction := *flags["direction"].(*string)
	action := *flags["action"].(*string)
	table := *flags["table"].(*string)

	if env < DEVELOPMENT || env > PRODUCTION {
		panic(fmt.Sprintf("Invalid environment: %d. Must be 1 (Dev), 2 (Test), or 3 (Prod)", env))
	}

	if action != "run" && action != "create" {
		panic(fmt.Sprintf("Invalid action: %s. Must be 'run' or 'create'", action))
	}

	if action == "run" && direction != "up" && direction != "down" {
		panic(fmt.Sprintf("Invalid direction: %s. Must be 'up' or 'down'", direction))
	}

	if action == "create" && strings.TrimSpace(table) == "" {
		panic("Table name (-table) is required when action is 'create'")
	}
}

func readDatabaseConfiguration() *DatabaseConfiguration {
	rawDbConfig, err := os.ReadFile("configs/databases/database.yml")
	dbconfig := DatabaseConfiguration{}
	if err != nil {
		panic(err)
	}
	if errun := yaml.Unmarshal(rawDbConfig, &dbconfig); errun != nil {
		panic(errun)
	}
	return &dbconfig
}

func takeInputFlags() map[string]any {
	result := make(map[string]any)
	env := flag.Int("env", 1, "Database environment to use (1: development, 2: test, 3: production)")
	direction := flag.String("direction", "up", "Migration direction to execute (up: apply migrations, down: rollback migrations)")
	action := flag.String("action", "run", "Task to perform (run: execute migrations on the database, create: generate new empty migration files)")
	table := flag.String("table", "", "Name of the table or feature for the new migration (required and only used when -action=create)")

	flag.Parse()

	result["environment"] = env
	result["direction"] = direction
	result["action"] = action
	result["table"] = table
	return result
}

func getAppropriateDbConfig(DBConfig *DatabaseConfiguration, environment int) map[string]any {
	var dbconfig map[string]any
	switch environment {
	case DEVELOPMENT:
		dbconfig = DBConfig.Development
	case TEST:
		dbconfig = DBConfig.Test
	case PRODUCTION:
		dbconfig = DBConfig.Production
	}
	return dbconfig
}

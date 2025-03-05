// migrator.go

package database_actions

import (
	"encoding/csv"
	"os"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var driver database.Driver

// InitMigrator initiates values essential for migrations
func InitMigrator(dsnMigrate string) error {
	var err error
	db, err := sql.Open("mysql", dsnMigrate)
	if err != nil {
		return fmt.Errorf("error while opening db connection: %w", err)
	}
	driver, err = mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("error while instanciating migration driver: %w", err)
	}

	return nil
}

// RunMigrate performs all or only some up/down migrations
//
// Default 'steps' as 0 (runs all migrations)
func RunMigrate(migrationType string, steps int) (string, error) {
	m, err := migrate.NewWithDatabaseInstance(
		"file://database_actions/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return "", fmt.Errorf("error while instanciating new migration ("+migrationType+") with DB : %w", err)
	}

	if steps != 0 {
		m.Steps(steps)
	} else {
		if migrationType == "up" {
			err = m.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				return "Migration(s) : " + migrate.ErrNoChange.Error(), nil
			}
			if err != nil {
				return "", fmt.Errorf("error while running up migration(s): %w", err)
			}
		} else if migrationType == "down" {
			err = m.Down()
			if errors.Is(err, migrate.ErrNoChange) {
				return "Migration(s) : " + migrate.ErrNoChange.Error(), nil
			}
			if err != nil {
				return "", fmt.Errorf("error while running down migration(s): %w", err)
			}
		} else {
			return "", fmt.Errorf("error unknown migration type: " + migrationType)
		}
	}

	return migrationsSuccessMessage(migrationType, steps), nil
}

func migrationsSuccessMessage(migrationType string, steps int) string {
	msg := "Successfully ran"
	if steps == 0 {
		return msg + " all " + migrationType + " migrations"
	}
	if steps == 1 {
		return msg + " 1 " + migrationType + " migration"
	}

	return msg + " " + strings.Trim(strconv.Itoa(steps), "-") + " " + migrationType + " migrations"
}

func InsertBreedsFromCSV(dsn string, csvFilePath string) error {
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}
	defer db.Close()

	// Prepare the insert statement
	stmt, err := db.Prepare("INSERT IGNORE INTO breeds (id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
			return fmt.Errorf("error preparing SQL statement: %w", err)
	}
	defer stmt.Close()

	// Iterate over the records and insert them into the database
	for _, record := range records[1:] { // Skip the header row
		_, err := stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			fmt.Printf("Error inserting record: %v\n", err)
			continue
		}
	}

	fmt.Println("Data successfully inserted into the database.")
	return nil
}
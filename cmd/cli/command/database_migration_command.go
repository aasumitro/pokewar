package command

import (
	"fmt"
	"github.com/aasumitro/pokewar/pkg/appconfig"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"strconv"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var migrateCmd = &cobra.Command{
	Use:  "migration",
	Long: `migrate cmd is used for database migration: migrate < up | down >`,
}

var migrateUpCmd = &cobra.Command{
	Use:  "up",
	Long: `Command to upgrade database migration`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			fmt.Printf("migrate down error: %v \n", err)
			return
		}

		if err := migration.Up(); err != nil {
			fmt.Printf("migrate up error: %v \n", err)
			return
		}

		fmt.Println("Migrate up done with success")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:  "down",
	Long: `Command to downgrade database`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			fmt.Printf("migrate down error: %v \n", err)
			return
		}

		if err := migration.Down(); err != nil {
			fmt.Printf("migrate down error: %v \n", err)
			return
		}

		fmt.Println("Migrate down done with success")
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:  "version",
	Long: `Command to see database migration version`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			fmt.Printf("migrate down error: %v \n", err)
			return
		}

		version, dirty, err := migration.Version()
		if err != nil {
			fmt.Printf("migrate up error: %v \n", err)
			return
		}

		fmt.Printf("Database Version %d is dirty: %s",
			version, strconv.FormatBool(dirty))
	},
}

func initGoMigrate() (instance *migrate.Migrate, err error) {
	fileSource, err := (&file.File{}).Open("file://db/migrations")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(
		appconfig.DbPool, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	instance, err = migrate.NewWithInstance(
		"file", fileSource, "pokewar", driver)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func init() {
	CliCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	"github.com/wavinamayola/user-management/internal/utils"
)

const driver = "mysql"

// format: goose mysql "root:secret@tcp(localhost:3306)/assessment_task?parseTime=true" status
func main() {
	flag.Parse()
	flag.Usage = usage
	args := flag.Args()

	cfg, err := utils.LoadConfig("./env/")
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}

	dataSource, err := utils.NewDBStringFromDBConfig(cfg)
	if err != nil {
		log.Fatalf("failed to set data source: %+v", err)
	}

	db, err := sql.Open(cfg.Database.Driver, dataSource)
	if err != nil {
		log.Fatalf("failed to open db conn: %+v", err)
	}
	defer func() { _ = db.Close() }()

	if err = goose.SetDialect(driver); err != nil {
		log.Fatal("failed to set driver for goose")
	}

	if len(args) == 0 {
		log.Fatal("expected at least one arg")
	}

	command := args[0]

	if err = goose.Run(command, db, cfg.Database.MigrationDir, args[1:]...); err != nil {
		log.Fatalf("goose run err: %+v", err)
	}

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	const (
		usageRun      = `goose [OPTIONS] COMMAND`
		usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version`
	)
	fmt.Println(usageRun)
	flag.PrintDefaults()
	fmt.Println(usageCommands)
}

package apiserver

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/assanoff/http-crud-server/internal/app/store/sqlstore"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq" // ...
	log "github.com/sirupsen/logrus"
)

// Start ...
func Start(config *Config) error {
	pool, err := pgxpool.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()
	log.Infof("Connected!")

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
	}
	// migrateDatabase(conn.Conn())
	MigrateDB(conn.Conn(), config)

	conn.Release()

	if err != nil {
		return err
	}
	store := sqlstore.New(pool)

	srv := newServer(store, config.Endpoint)

	return http.ListenAndServe(":"+config.Port, srv)

}

// MigrateDB ...
func MigrateDB(conn *pgx.Conn, config *Config) {

	templ, err := ioutil.ReadFile("./migrations/0001_create_user.sql")
	if err != nil {
		log.Fatalf("Unable to open a migrate file: %v", err)
	}

	var templateScript bytes.Buffer
	var migrateTemplate = template.Must(template.New("migrateTempl").Parse(string(templ)))
	err = migrateTemplate.Execute(&templateScript, config)
	if err != nil {
		log.Fatalf("Unable to templ a migrate: %v", err)
	}
	migrateScript := templateScript.String()

	_, err = conn.Exec(context.Background(), migrateScript)
	if err != nil {
		log.Errorf("Unable to migrate: %v\n", err)
		return
	}

}

// func migrateDatabase(conn *pgx.Conn) {
// 	migrator, err := migrate.NewMigrator(conn, "test")
// 	if err != nil {
// 		log.Fatalf("Unable to create a migrator: %v", err)
// 	}

// 	err = migrator.LoadMigrations("./migrations")
// 	if err != nil {
// 		log.Fatalf("Unable to load migrations: %v", err)
// 	}

// 	err = migrator.Migrate(func(err error) (retry bool) {
// 		log.Infof("Commit failed during migration, retrying. Error: %v", err)
// 		return true
// 	})

// 	if err != nil {
// 		log.Fatalf("Unable to migrate: %v", err)
// 	}

// 	ver, err := migrator.GetCurrentVersion()
// 	if err != nil {
// 		log.Fatalf("Unable to get current schema version: %v", err)
// 	}

// 	log.Infof("Migration done. Current schema version: %v", ver)
// }

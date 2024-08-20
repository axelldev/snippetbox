package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/axelldev/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger          *slog.Logger
	snippets        *models.SnippetModel
	temaplatesCache templateCache
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dns := flag.String("dns", "web:password@/snippetbox?parseTime=true", "MYSQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	db, err := openDb(*dns)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	app := &application{
		logger:          logger,
		snippets:        &models.SnippetModel{DB: db},
		temaplatesCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.router())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}

package main

//to run both files simultaneously go run .
import (
	"database/sql"
	"flag"
	"github.com/yash/snippetbox/pkg/models/postgres"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "user=dummy dbname=snippetbox password=pass sslmode=disable", "Postgres Datasource Name")
	flag.Parse()

	//question 1: log.New function is expected to have io.Writer type data or we can say something which has write method so that it can implement writer interface
	// but here os.Stdout does not have a write method so cannot implement writer interface so how is this working or valid?
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &postgres.SnippetModel{Pool: db},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	//writing messages into the log using new loggers
	infoLog.Printf("Starting server on :%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

/*
f, err:= os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err!=nil{
		log.Fatal(err)
	}
	defer f.Close()
This can be used if u dont want to do logging at runtime
*/

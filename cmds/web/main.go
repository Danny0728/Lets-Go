package main

//to run both files simultaneously go run .
import (
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/yash/snippetbox/pkg/models/postgres"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	session       *sessions.Session
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "user=dummy dbname=snippetbox password=pass sslmode=disable", "Postgres Datasource Name")
	secret := flag.String("secret", "s6Nd%+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

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
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour //session will expire after 12 hours

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
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
This can be used if u don't want to do logging at runtime
*/

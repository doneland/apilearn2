package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	host     = "apipostgres"
	port     = 5432
	user     = "apiadmin"
	password = "Werty132"
	dbname   = "apidb"
)

var (
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db      *sqlx.DB
)

func init() {
	// Load values from evironment file .env.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Will move forward with environment variables or default values.")
		cwd, _ := os.Getwd()
		log.Println("Current diretory:", cwd)
	}

	host, _ = os.LookupEnv("SQL_HOST")
	p, _ := os.LookupEnv("SQL_PORT")
	port, _ = strconv.Atoi(p)
	user, _ = os.LookupEnv("SQL_USER")
	password, _ = os.LookupEnv("SQL_PASSWORD")
	dbname, _ = os.LookupEnv("SQL_DATABASE")

	log.Printf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

// App main environment component.
type App struct {
	DB *sqlx.DB
}

// Get returns base data.
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": "works"}`))
}

// Cats return information about cats.
func (a *App) Cats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": "bad cats"}`))
}

// Trxs returns transactions.
func (a *App) Trxs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	trxs, err := FetchTrxs(a.DB)
	if err != nil {
		log.Fatalln(err)
	}
	json.NewEncoder(w).Encode(trxs)
}

// TrxsAdd creates a transaction.
func (a *App) TrxsAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received.")
	t := &Trx{}

	fmt.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Start sleep.")
	time.Sleep(time.Second * 10)
	log.Print("Finished sleep.")

	err = FetchTrxSave(a.DB, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func main() {

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Connection to DB has failed, err: %s", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("DB Ping has failed.")
	}
	defer db.Close()

	err = FetchCreateTrxTable(db)
	if err != nil {
		log.Fatalln(err)
	}

	app := &App{DB: db}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", app.Get).Methods(http.MethodGet)
	api.HandleFunc("/cats", app.Cats).Methods(http.MethodGet)
	api.HandleFunc("/trxs", app.Trxs).Methods(http.MethodGet)
	api.HandleFunc("/trxsadd", app.TrxsAdd).Methods(http.MethodPost)

	log.Fatalln(http.ListenAndServe(":8080", r))
}

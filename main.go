package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	host     = "apipostgres"
	port     = 5433
	user     = "apiadmin"
	password = "Werty132"
	dbname   = "apidb"
)

var (
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db      *sqlx.DB
)

func init() {
	/*host, _ = os.LookupEnv("DB_HOST")
	p, _ := os.LookupEnv("DB_PORT")
	port, _ = strconv.Atoi(p)
	user, _ = os.LookupEnv("DB_USER")
	password, _ = os.LookupEnv("DB_PASSWORD")
	dbname, _ = os.LookupEnv("DB_DBNAME")*/

	log.Printf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

	log.Fatalln(http.ListenAndServe(":8080", r))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "apilearn_admin"
	password = "Werty132"
	dbname   = "apilearn"
)

var (
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db      *sqlx.DB
)

/*func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Cannot opent .env file.")
	}

	host, _ = os.LookupEnv("DB_HOST")
	p, _ := os.LookupEnv("DB_PORT")
	port, _ = strconv.Atoi(p)
	user, _ = os.LookupEnv("DB_USER")
	password, _ = os.LookupEnv("DB_PASSWORD")
	dbname, _ = os.LookupEnv("DB_DBNAME")

	log.Printf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}*/

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": "works"}`))
}

func cats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": "found cats"}`))
}

func main() {

	/*db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Connection to DB has failed, err: %s", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("DB Ping has failed.")
	}
	defer db.Close()*/

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", get).Methods(http.MethodGet)
	api.HandleFunc("/cats", cats).Methods(http.MethodGet)

	log.Fatalln(http.ListenAndServe(":8080", r))
}

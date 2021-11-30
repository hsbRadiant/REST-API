package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hsbRadiant/restapi/configuration"
)

var router = mux.NewRouter()

// FUNCTION TO CREATE A 'BOOKS' TABLE IF IT DOES NOT EXISTS :
func createBooksTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS books(
			id INT NOT NULL AUTO_INCREMENT, 
			title VARCHAR(255) NOT NULL,
			description VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func createAuthorsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS authors(
			id INT NOT NULL AUTO_INCREMENT, 
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func clearBooksTable() {
	db.Exec(`DELETE from books`)
	db.Exec(`ALTER TABLE books AUTO_INCREMENT = 1`)
}

func clearAuthorsTable() {
	db.Exec(`DELETE from authors`)
	db.Exec(`ALTER TABLE authors AUTO_INCREMENT = 1`)
}

func clearUsersTable() {
	db.Exec(`DELETE from users`)
}

func checkStatusEquality(t *testing.T, expectedCode, actualCode int) {
	if actualCode != expectedCode {
		t.Errorf("Expected code : %v, Actual code : %v", expectedCode, actualCode)
	}
}

func TestMain(m *testing.M) {
	db = configuration.NewDatabase()
	// router = mux.NewRouter() // Using it not able to execute ServeHttp function properly - I THINK.
	// routes.Setup(router)
	router = Router()

	createBooksTable(db)
	createAuthorsTable(db)

	fmt.Println("CHECK - Before running")
	exitCode := m.Run() // 'Run' runs the test case and returns an exit code to pass to 'os'.
	fmt.Println("CHECK - After running")

	clearBooksTable()
	clearAuthorsTable()
	// db.Close()
	os.Exit(exitCode)
}

func Router() *mux.Router {
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books/create", CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	router.HandleFunc("/authors", GetAuthors).Methods("GET")
	router.HandleFunc("/authors/{id}", GetAuthor).Methods("GET")
	router.HandleFunc("/authors/create", CreateAuthor).Methods("POST")
	router.HandleFunc("/authors/{id}", UpdateAuthor).Methods("PUT")
	router.HandleFunc("/authors/{id}", DeleteAuthor).Methods("DELETE")

	router.HandleFunc("/register", RegisterUser).Methods("POST")

	return router
}

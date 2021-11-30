// Package classification REST API (Go).
//
// A RESTful application / APIs built using Go.
//
// Terms Of Service:
//
// No TOS at this moment
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: API Support<support@swagger.io> http://www.swagger.io/support
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: JWT
//          in: header
// swagger:meta
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	configuration "github.com/hsbRadiant/restapi/configuration"
	_ "github.com/hsbRadiant/restapi/docs"

	// _ "github.com/hsbRadiant/restapi/docsSwagger"
	routes "github.com/hsbRadiant/restapi/routes"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// DEFINING DATABASE VARIABLE(S) :
var db *sql.DB

// THE SIGNING KEY FOR JWT :
var SigningKey = []byte("harskey")

// To check an error anywhere :
func CheckError(e error) {
	if e != nil {
		log.Fatal(e.Error())
		// panic(e.Error())
	}
}

// FUNCTION TO CREATE A 'BOOKS' TABLE IF IT DOES NOT EXISTS :
// func createBooksTable(database *sql.DB) {
func CreateBooksTable() {
	query := `CREATE TABLE IF NOT EXISTS books(
		id INT NOT NULL AUTO_INCREMENT, 
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		);`

	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	// Executing the query in a context using ExecContext - READ?
	_, err := db.ExecContext(cntx, query)
	CheckError(err)
	// return err
}

// FUNCTION TO CREATE AN 'AUTHORS' TABLE IF IT DOES NOT EXISTS :
// func createAuthorsTable(db *sql.DB) {
// func CreateAuthorsTable() {
// 	query := `CREATE TABLE IF NOT EXISTS authors(
// 		id INT NOT NULL AUTO_INCREMENT,
// 		name VARCHAR(255) NOT NULL,
// 		book_id INT,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		PRIMARY KEY (id),
// 		FOREIGN KEY (book_id) REFERENCES books(id)
// 		);`

// 	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancelFunc()

// 	// Executing the query in a context using ExecContext - READ?
// 	_, err := db.ExecContext(cntx, query)
// 	CheckError(err)
// }
func CreateAuthorsTable() {
	query := `CREATE TABLE IF NOT EXISTS authors(
		id INT NOT NULL AUTO_INCREMENT, 
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		);`

	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	// Executing the query in a context using ExecContext - READ?
	_, err := db.ExecContext(cntx, query)
	CheckError(err)
}

// func AlterAuthorsTable() {
// 	query := `ALTER TABLE authors
// 		ADD FOREIGN KEY (book_id) REFERENCES books(id);`

// 	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancelFunc()

// 	// Executing the query in a context using ExecContext - READ?
// 	_, err := db.ExecContext(cntx, query)
// 	CheckError(err)
// }

// FUNCTION TO CREATE AN 'AUTHORS' TABLE IF IT DOES NOT EXISTS :
// func createAuthorsTable(db *sql.DB) {
// func CreateBooksAuthorsJoinTable() {
// 	query := `CREATE TABLE IF NOT EXISTS booksauthors as (SELECT books.id as book_id,
// 		books.title as book_title, books.description as book_description, authors.name as author_name
// 		FROM books INNER JOIN authors WHERE books.id = authors.book_id)`

// 	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancelFunc()

// 	// Executing the query in a context using ExecContext - READ?
// 	_, err := db.ExecContext(cntx, query)
// 	CheckError(err)
// 	// return err
// }

// FUNCTION TO CREATE A 'USERS' TABLE IF IT DOES NOT EXISTS :
func CreateUsersTable() {
	query := `CREATE TABLE IF NOT EXISTS users(
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		PRIMARY KEY (email)
		);`

	cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	// Executing the query in a context using ExecContext - READ?
	_, err := db.ExecContext(cntx, query)
	CheckError(err)
	// return err
}

// @title Swagger REST API
// @version 1.0
// @description A Book - Author REST API
// @termsOfService http://swagger.io/terms
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// CONFIGURING DATABASE :
	db = configuration.DB // (var 'DB' in 'config.go' stores the value of type '*sql.DB')

	// Deferred closing of the connection :
	defer db.Close()

	// Checking the database connection :
	errPing := db.Ping()
	CheckError(errPing)
	fmt.Println("Successfully connected to the database") // Display a message if 'ping' connects successfully.

	// Creating books table in database :
	// createBooksTable(db)
	// createAuthorsTable(db)
	// createUsersTable(db)
	CreateBooksTable()
	CreateAuthorsTable()
	// AlterAuthorsTable()
	// CreateBooksAuthorsJoinTable()
	CreateUsersTable()

	// var validate *validator.Validate

	/*
		MUX.ROUTER - A 'mux.Router' matches the incoming requests against a list of registered routes and calls a handler for the route that mathces the URL.
		It implements the 'http.Handler' interface.
		Requests can be matched based on URL host, path, path prefix, schemes, header and query values, HTTP methods or custom mathcers.
	*/

	// INIT A ROUTER :
	router := mux.NewRouter()

	// MATHCING ROUTES TO ENDPOINTS :
	routes.Setup(router)

	// SWAGGER IMPLEMENTATION :
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// CORS handling :
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"})
	age := handlers.MaxAge(0)

	// RUN THE SERVER : (Call 'ListenAndServer' method of the http package)
	// log.Fatal(http.ListenAndServe(":3000", router)) // ListenAndServe - takes two params. : 1. address of the port  2. http handler - CHECK?
	// ListenAndServe with CORS handled :-
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, headers, origins, methods, age)(router)))
}

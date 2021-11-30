package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	models "github.com/hsbRadiant/restapi/models"
)

// handler := http.HandlerFunc(controllers.GetBook)
// handler.ServeHTTP(response, req)
// controllers.GetBook(response, req)
// res := response.Result()
// defer res.Body.Close()
// if res.StatusCode != http.StatusFound {

func addBooks(n int) {
	if n < 1 {
		n = 1
	}
	for i := 0; i < n; i++ {
		db.Exec("INSERT INTO books (title, description) values (?,?)", "Book "+strconv.Itoa(i+1), "Test Book")
	}
}

func TestEmptyBooksTable(t *testing.T) {
	clearBooksTable() // deleting the table first and then testing.

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder() // 'NewRecorder' returns an initialized 'ResponseRecorder'.
	router.ServeHTTP(response, req)

	// if response.Code != http.StatusNotFound {
	// 	t.Errorf("Expected Code - %v / Received Code - %v", http.StatusOK, response.Code)
	// }
	checkStatusEquality(t, http.StatusNotFound, response.Code)

	if response.Body.String() != "Empty table" {
		t.Errorf("Exepected an empty array but got %v", response.Body)
	}
}

func TestGetBooks(t *testing.T) {
	clearBooksTable()
	addBooks(2) // inserting 2 test books in the database first.

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// if response.Code != http.StatusOK {
	// 	t.Errorf("Exepected : %v, Actual : %v", http.StatusOK, response.Code)
	// }
	checkStatusEquality(t, http.StatusOK, response.Code)
}

func TestGetBook(t *testing.T) {
	clearBooksTable()
	addBooks(3)

	req, err := http.NewRequest("GET", "/books/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// if response.Code != http.StatusOK {
	// 	t.Errorf("Expected code - %v, Received code - %v", http.StatusOK, response.Code)
	// }
	checkStatusEquality(t, http.StatusOK, response.Code)

	var book models.Book
	json.Unmarshal(response.Body.Bytes(), &book)

	if book.ID != 2 && book.Title != "Book 2" && book.Description != "Test Book" {
		t.Errorf("Expected ID, Title, Description : %v, %v, %v \n Actual ID, Title, Description : %v, %v, %v", 2, "Book 2", "Test Book", book.ID, book.Title, book.Description)
	}
}

func TestGetbookNotFound(t *testing.T) {
	clearBooksTable()

	req, err := http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	checkStatusEquality(t, http.StatusNotFound, response.Code)

	var book models.Book
	json.Unmarshal(response.Body.Bytes(), &book)

	var emptybook models.Book
	if book.ID != emptybook.ID && book.Title != emptybook.Title && book.Description != emptybook.Description {
		t.Errorf("Expected an empty struct but received %v", book)
	}
}

func TestCreateBook(t *testing.T) {
	clearBooksTable()
	var payload = []byte(`{"title":"NEW BOOK 1","description":"Test Book"}`)

	req, err := http.NewRequest("POST", "/books/create", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	// fmt.Println("START", strings.SplitAfterN(response.Body.String(), "}}", 1)[0], "END")

	// if response.Code != http.StatusOK {
	// 	t.Errorf("Expected code - %v, Received code - %v", http.StatusOK, response.Code)
	// }
	checkStatusEquality(t, http.StatusOK, response.Code)
	var actualBook models.Book
	json.Unmarshal(response.Body.Bytes(), &actualBook)

	if actualBook.ID != 1 || actualBook.Title != "NEW BOOK 1" || actualBook.Description != "Test Book" {
		t.Errorf("Expected ID, Title, Description : %v, %v, %v \n Actual ID, Title, Description : %v, %v, %v", 1, "NEW BOOK 1", "Test Book", actualBook.ID, actualBook.Title, actualBook.Description)
	}
}

func TestUpdateBook(t *testing.T) {
	clearBooksTable()
	addBooks(1)

	// 1. Getting the book stored at id = 1 :
	req, err := http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	var originalBook models.Book
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	// 2. Updating the book with id = 1 :
	payload := `{"title":"Updated Book 1"}`

	req, err = http.NewRequest("PUT", "/books/1", bytes.NewBuffer([]byte(payload))) // The caller should not use buff ('[]byte(payload)') after it has been passed to the 'NewBuffer' function.
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)
	// fmt.Println(response.Body)

	// 3. Checking the book got updated by again getting the book with id = 1 :
	req, err = http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	// fmt.Println(response.Body)
	checkStatusEquality(t, http.StatusOK, response.Code)
	var updatedBook models.Book
	json.Unmarshal(response.Body.Bytes(), &updatedBook)

	if updatedBook.ID != originalBook.ID || updatedBook.Title != "Updated Book 1" || updatedBook.Description != originalBook.Description {
		t.Errorf(`Expected credentials : ID - %v, Title - %v, Description - %v
		 Actual credentials : ID - %v, Title - %v, Description - %v`, originalBook.ID, "Updated Book 1", updatedBook.Description, originalBook.ID, originalBook.Title, originalBook.Description)
	}
}

func TestDeleteBook(t *testing.T) {
	clearBooksTable()
	addBooks(1)

	// Testing a user exists at id 1 :
	req, err := http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	// Deleting the user :
	req, err = http.NewRequest("DELETE", "/books/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	// Retesting again to confirm the user is deleted :
	req, err = http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusNotFound, response.Code)
}

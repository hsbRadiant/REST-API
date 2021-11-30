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

func addAuthors(n int) {
	if n < 1 {
		n = 1
	}
	for i := 0; i < n; i++ {
		db.Exec("INSERT INTO authors (name) values (?)", "Author "+strconv.Itoa(i+1))
	}
}

func TestEmptyAuthorsTable(t *testing.T) {
	clearAuthorsTable() // deleting the table first and then testing.

	req, err := http.NewRequest("GET", "/authors", nil)
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
		t.Errorf("Expected an empty array but got %v", response.Body)
	}
}

func TestGetAuthors(t *testing.T) {
	clearAuthorsTable()
	addAuthors(3) // inserting 2 test authors in the database first.

	req, err := http.NewRequest("GET", "/authors", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	checkStatusEquality(t, http.StatusOK, response.Code)
}

func TestGetAuthor(t *testing.T) {
	clearAuthorsTable()
	addAuthors(3)

	req, err := http.NewRequest("GET", "/authors/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	checkStatusEquality(t, http.StatusOK, response.Code)

	var author models.Author
	json.Unmarshal(response.Body.Bytes(), &author)

	if author.ID != 2 && author.Name != "Author 2" {
		t.Errorf("Expected ID, Name : %v, %v \n Actual ID, Name : %v, %v", 2, "Author 2", author.ID, author.Name)
	}
}

func TestGetAuthorNotFound(t *testing.T) {
	clearAuthorsTable()

	req, err := http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	checkStatusEquality(t, http.StatusNotFound, response.Code)

	var author models.Author
	json.Unmarshal(response.Body.Bytes(), &author)

	var emptyAuthor models.Author
	if author.ID != emptyAuthor.ID && author.Name != emptyAuthor.Name {
		t.Errorf("Expected an empty struct but received %v", author)
	}
}

func TestCreateAuthor(t *testing.T) {
	clearAuthorsTable()
	var payload = []byte(`{"name":"New Author"}`)

	req, err := http.NewRequest("POST", "/authors/create", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	var createdAuthor models.Author
	json.Unmarshal(response.Body.Bytes(), &createdAuthor)

	if createdAuthor.ID != 1 || createdAuthor.Name != "New Author" {
		t.Errorf("Expected ID, Name : %v, %v\n Actual ID, Name : %v, %v", 1, "New Author", createdAuthor.ID, createdAuthor.Name)
	}
}

func TestUpdateAuthor(t *testing.T) {
	clearAuthorsTable()
	addAuthors(1)

	// 1. Getting the author stored at id = 1 :
	req, err := http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	var originalAuthor models.Author
	json.Unmarshal(response.Body.Bytes(), &originalAuthor)

	// 2. Updating the author with id = 1 :
	payload := `{"name":"Updated Author"}`

	req, err = http.NewRequest("PUT", "/authors/1", bytes.NewBuffer([]byte(payload))) // The caller should not use buff ('[]byte(payload)') after it has been passed to the 'NewBuffer' function.
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)
	// fmt.Println(response.Body)

	// 3. Checking the author got updated by again getting the author with id = 1 :
	req, err = http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	// fmt.Println(response.Body)
	checkStatusEquality(t, http.StatusOK, response.Code)

	var updatedAuthor models.Author
	json.Unmarshal(response.Body.Bytes(), &updatedAuthor)

	if updatedAuthor.ID != originalAuthor.ID || updatedAuthor.Name != "Updated Author" {
		t.Errorf(`Expected credentials : ID - %v, Name - %v
		 Actual credentials : ID - %v, Name - %v`, originalAuthor.ID, "Updated Author", updatedAuthor.ID, updatedAuthor.Name)
	}
}

func TestDeleteAuthor(t *testing.T) {
	clearAuthorsTable()
	addAuthors(1)

	// Testing a user exists at id 1 :
	req, err := http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	// Deleting the user :
	req, err = http.NewRequest("DELETE", "/authors/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)

	// Retesting again to confirm the user is deleted :
	req, err = http.NewRequest("GET", "/authors/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusNotFound, response.Code)
}

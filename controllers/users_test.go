package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessfulUserRegistration(t *testing.T) {
	clearUsersTable()

	var payload = []byte(`{"name":"Harsimran Singh Bedi","email":"hsbhsb@gmail.com","password":"123456123456"}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err.Error())
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	checkStatusEquality(t, http.StatusOK, response.Code)
}

func TestUnsuccessfulUserRegistration(t *testing.T) {
	clearUsersTable()

	var payload = []byte(`{"name":"Harsimran Singh Bedi","email":"hsbhsb@gmail.com","password":"123456"}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err.Error())
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code == http.StatusOK {
		t.Errorf("Got Wrong status code. Expected : %v, Received : %v", http.StatusBadRequest, response.Code)
	}

	if response.Body.String() == "Registration successful." {
		t.Errorf("There is some error while registering the user. Registration should not be succesful. Check the error.")
	}
}

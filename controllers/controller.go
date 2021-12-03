package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	configuration "github.com/hsbRadiant/restapi/configuration"
	models "github.com/hsbRadiant/restapi/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	db         *sql.DB = configuration.DB  // DATBASE VARIABLE
	SigningKey         = []byte("harskey") // THE SIGNING KEY FOR JWT // []byte(os.Getenv("signingKey"))
	// Books      []models.Book                          // INITIALIZE A BOOK
	// emailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	// emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// To check an error anywhere and log it :
func CheckError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func CheckErrorAsJSON(w http.ResponseWriter, message string, code int) {
	response, _ := json.Marshal(map[string]string{"error": message}) // creating a new instance of map type and passing it to be marshalled into JSON.
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

// func setContentType(rw *http.ResponseWriter) {
// 	(*rw).Header().Set("Content-Type", "application/json; charset=utf-8") // took reference from Error function definition for 'charset' property (READ about 'charset=utf-8')
// }

// Was initially taking a pointer type parameter - Understand any differene in taking a pointer type param. to a normal (value type) param.
func setContentType(rw http.ResponseWriter) {
	(rw).Header().Set("Content-Type", "application/json; charset=utf-8") // took reference from Error function definition for 'charset' property (READ about 'charset=utf-8')
}

// Validating the given email address for the correct form using a 'RegExp' :-
// func isValidEmail(email string) bool {
// 	// matchedBool,err := regexp.MatchString(emailRegex, email)
// 	// return matchedBool, nil
// 	return emailRegex.MatchString(email)
// }

type CustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

/*
 ~ The function ('AuthorizeClientRequest') takes in an endpoint (a fn. type with a 'repsonse type' and a 'request type' parameter).
 ~ It returns an 'http.HandlerFunc' type.
 ~ (Previously was returning 'http.Handler' type. It would had been correct had the handler been defined as :
   'router.Handle("/books", controllers.AuthorizeClientRequest(controllers.GetBooks)).Methods("GET")'
 ~ 'http.HandlerFunc(f)' allows to use any ordinary function (f) with appropriate signature to be used as an 'http.Handler'.
*/
// Creating a function to authorize requests first and then serve the endpoint or any resource :
func AuthorizeClientRequest(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	var f = func(w http.ResponseWriter, r *http.Request) {
		setContentType(w)
		cookie, err := r.Cookie("JWT")
		if err != nil {
			if err == http.ErrNoCookie {
				// setContentType(w)
				// w.WriteHeader(http.StatusUnauthorized)
				// (http.Error or Error functions write a text/plain type content to the header - SEE DEFINITION OF 'Error' fn)
				// http.Error(w, "No cookie sent with the request. ERROR - "+err.Error(), http.StatusUnauthorized)
				CheckErrorAsJSON(w, "(As JSON) No cookie sent with the request. "+err.Error(), http.StatusUnauthorized)
				return
			}
			// http.Error(w, err.Error(), http.StatusBadRequest)
			CheckErrorAsJSON(w, "(As JSON) "+err.Error(), http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value // JWT string stored as a value of the cookie sent is stored in 'tokenString'.

		// An instance of the claims struct :
		claims := &CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		})
		// Checking for errors :
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.Write([]byte("THE SIGNATURE IS INVALID"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(err.Error())
			return
		}

		// Is the token valid (CHECK the documentation of 'Token' type to find a 'Valid' field name in the 'Token' type struct. It gets populated when one parses / verifies a token.
		if _, ok := token.Claims.(*CustomClaims); !ok && !token.Valid {
			// if !token.Valid {
			// w.Write([]byte("Invalid token."))
			// w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Calling the endpoint function with args. :
		endpoint(w, r)
	}
	return http.HandlerFunc(f)
}

// DEFINING THE HANDLER FUNCTIONS :

// /////////////////////////////////////////////////////
// /////////////////////// USERS //////////////////////

// RegisterUser godoc
// @Summary To register a user
// @Description Register User
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Register User"
// @Success 200 {string} string "Registration successful"
// @Failure 400 {string} string "StatusBadRequest"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	// Will help in handling 'EOF' :
	if r.Body == http.NoBody {
		fmt.Fprintf(w, "Please provide the request body.")
		return
	}

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	// fmt.Println(user)

	// Calling email validation fn. (using 'regexp' package) :-
	// if !isValidEmail(user.Email) {
	// 	fmt.Fprintf(w, "Please provide a correct form of email address.")
	// 	return
	// }

	// Email validation using 'net/mail' package (with the help 'ParseAddress' function)
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(w, "Please provide a correct form of email address.")
		http.Error(w, "Please provide a correct form of email address. ERROR - "+err.Error(), http.StatusBadRequest)
		return
	}

	// e := user.Validate() // gives http panic - invalid memory address or nil pointer dereference.
	if user.Validate() != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte(user.Validate().Error()))
		http.Error(w, user.Validate().Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		user.Name = "Unknown"
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Fatal(err.Error())
	}

	encryptedPassword := string(bytes)

	// fmt.Println(user) // Notice, if any field is left empty then the zero value of the field gets stored in the db. (It satisfies NOT NULL constraint of the db but what is needed is the field should not have values zero values like this)
	// cntx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancelFunc()
	// db.ExecContext(cntx, "INSERT INTO users (name, email, password) VALUES (?,?,?)")

	stmt, err := db.Prepare("INSERT INTO users (name,email,password) VALUES (?,?,?)")
	CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, encryptedPassword)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("Could not execute the statement.\n" + err.Error()))
		http.Error(w, "Could not execute the statement. "+err.Error(), http.StatusInternalServerError)
		return
		// log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Registration successful.")
}

// UserLogin godoc
// @Summary To login a registered user
// @Description Login User
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Login request body parameters"
// @Success 200 {string} string "WELCOME! 'user.Email'. Login successful."
// @Failure 400 {string} string "StatusBadRequest"
// @Failure 401 {string} string "StatusUnathorized"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /login [post]
func UserLogin(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	if r.Body == http.NoBody {
		fmt.Fprintf(w, "Please provide login credentials in the request body.")
		return
	}

	// Firstly, authenticate the user / client :-
	// var user models.User
	var user map[string]string
	// Storing the user credentials in a new User type variable :
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	// res, err := db.Exec("SELECT email FROM users WHERE email = ?", user.Email) // Number of rows usually chnaged when updating, inserting, replacing, etc. queries. With SELECT will return number of rows = 0 always - I THINK.
	// CheckError(err)
	// if number, _ := res.RowsAffected(); number == 0 {
	// 	fmt.Fprintf(w, "The email %v does not exists. Please provide a correct one", user.Email)
	// 	return
	// }

	// Getting the expected password from the database :-
	var expectedUser models.User
	row := db.QueryRow("SELECT email, password from users WHERE email = ?", user["email"])

	err = row.Scan(&expectedUser.Email, &expectedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// fmt.Fprintf(w, "The email %q does not exists. Please provide the correct email address.", user.Email)
			http.Error(w, "The user (email address) does not exists", http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}

	// Checking if passwords match :
	if err := bcrypt.CompareHashAndPassword([]byte(expectedUser.Password), []byte(user["password"])); err != nil {
		// w.Write([]byte("UNAUTHENTICATED USER - Check the password."))
		// w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, "UNAUTHENTICATED USER (Incorrect password) - Check the password.", http.StatusUnauthorized)
		return
	}

	// Creating an instance of a claim (or payload)
	claims := &CustomClaims{
		Name:  user["email"],
		Email: user["email"],
		StandardClaims: jwt.StandardClaims{
			Issuer:    user["email"],
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 'JWT' takes in 'Unix' type time
		},
	}

	// Creating a JWT token using NewWithClaims func.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Takes in the signing method and the claims as parameters.

	// Generate the signed JWT string (the 3 part string) using the token and the key :
	tokenString, err := token.SignedString(SigningKey)

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Could not log in. ERROR - "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Using the token string generated to be set as a cookie to be sent to the client :
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 1), // Try to check difference between defining a variable for using here to directly using the value here.
		HttpOnly: true,                          // READ?
	})
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "WELCOME! %q. Login successful.", user.Email)
	json.NewEncoder(w).Encode("WELCOME! " + user["email"] + ". Login successful.")
	// w.Write([]byte("WELCOME! " + user.Email + ". Login successful."))
}

// UserLogout godoc
// @Summary To logout
// @Description UserLogout
// @Tags users
// @Success 200 {string} string "Logout successful."
// @Failure 400 {string} string "StatusBadRequest"
// @Router /logout [post]
func UserLogout(w http.ResponseWriter, r *http.Request) {
	// Setting the expiry time of a cookie in the past is the way to remove a cookie from the browser.
	cookie := http.Cookie{
		Name:     "JWT",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // The cookie expired one hour ago. (-'time.Hour' is same '(-1) * time.Hour')
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logout successful.")
}

// /////////////////////////////////////////////////////
// /////////////////////// BOOKS //////////////////////

// GetBooks godoc
// @Summary To get details of all books
// @Description Get Books
// @Tags books
// @Produce json
// @Success 200 {array} models.Book "All books"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /books [get]
func GetBooks(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	// Defining the slice of Book types here because if defined outside as a global variable then when appending a book, the global Books variable
	// will get appended which will contain resuls from previous get request also and will also have the results appended by the current get request.
	// This is the reason why the JSON result / response of the request in POSTMAN has everytime an increased slice and repeated of books.

	var Books []models.Book

	// With database :
	result, err := db.Query("SELECT id, title, description FROM books")
	// result, err := db.Query("SELECT book_id, book_title, book_description, author_name from booksauthors")
	CheckError(err)
	defer result.Close() // to free the result set of resources.

	if !result.Next() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Empty table"))
		return
	}
	for result.Next() { // 'for' is like the 'while' loop (same as 'while result.Next()' i.e. there is a next row to be scanned in the result variable).
		var book models.Book
		// var author models.Author
		err = result.Scan(&book.ID, &book.Title, &book.Description)
		// err = result.Scan(&book.ID, &book.Title, &book.Description, &book.Author.Name)
		// CheckError(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Books = append(Books, book)
	}
	// fmt.Println(Books)
	setContentType(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Books)
	// res, _ := json.Marshal(Books)
	// w.Write(res)
}

// GetBook godoc
// @Summary To get details of a book with it's id.
// @Description Get Book by ID
// @Tags books
// @Produces json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book "Book stored at the given ID"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /books/{id} [get]
func GetBook(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// setContentType(&w)
	params := mux.Vars(r) // 'mux.Vars' creates a map ('string' type key - value pair) of route variables / parameters.

	// With database :
	row := db.QueryRow("SELECT id, title, description FROM books WHERE id = ?", params["id"]) // ('db.Query' takes args. as any placeolder parameters)

	var book models.Book

	err := row.Scan(&book.ID, &book.Title, &book.Description)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// w.WriteHeader(http.StatusNotFound)
			// fmt.Fprintf(w, "Book with id %v not found", params["id"])
			http.Error(w, "Book with id "+params["id"]+" not found", http.StatusNotFound)
			return
		default:
			// w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err.Error())
			// return
		}
	}
	setContentType(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

	// // Using reflect package's 'DeepEqual' method to check whether the book is empty (i.e. the values of the fields is empty) or not
	// if (reflect.DeepEqual(book, Book{})) {
	// 	json.NewEncoder(w).Encode(fmt.Sprintf("Book with id %v not found", params["id"]))
	// } else {
	// 	json.NewEncoder(w).Encode(book)
	// }
}

// CreateBook godoc
// @Summary To create a new book
// @Description Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "New book request body parameters"
// @Success 200 {object} models.Book "New book created"
// @Failure 500 {string} string "InternalServerError"
// @Failure default {string} string "DefaultError"
// @Router /books/create [post]
func CreateBook(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	setContentType(w)
	var book models.Book
	// fmt.Fprintf(w, "REQUEST BODY : %v\n", r.Body)
	err := json.NewDecoder(r.Body).Decode(&book)
	// fmt.Fprintln(w, book)
	if err != nil {
		CheckErrorAsJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("CHECK")

	// if book.Title == "" || book.Author.Name == "" {
	// 	json.NewEncoder(w).Encode("Title / Author name cannot be empty")
	// 	return
	// }
	// fmt.Println("CHECK")

	stmt, e := db.Prepare("INSERT INTO books (title, description) VALUES (?,?)")
	if e != nil {
		CheckErrorAsJSON(w, e.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	res, e := stmt.Exec(book.Title, book.Description)
	if e != nil {
		CheckErrorAsJSON(w, "Not able to insert in the Database. Error = "+e.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		book.ID = int(id)
		book.Author = models.Author{}
	}

	// stmt, e = db.Prepare("INSERT INTO authors (name, book_id) VALUES (?,?)")
	// if e != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// defer stmt.Close()
	// _, e = stmt.Exec(book.Author.Name, book.ID)
	// if e != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, "Not able to insert in the Database. Error = "+e.Error())
	// 	return
	// }

	// stmt, e = db.Prepare("INSERT INTO booksauthors (book_id, book_title, book_description, author_name) VALUES (?,?,?,?)")
	// if e != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// defer stmt.Close()
	// _, e = stmt.Exec(book.ID, book.Title, book.Description, book.Author.Name)
	// if e != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, "Not able to insert in the Database. Error = "+e.Error())
	// 	return
	// }

	// if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
	// 	id, _ := res.LastInsertId()
	// 	book.ID = int(id)
	// 	if book.Author == nil {
	// 		book.Author = &models.Author{}
	// 	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
	// fmt.Fprintf(w, "New book inserted into the database")
}

// UpdateBook godoc
// @Summary To update a book with the given id
// @Description Update Book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Update book request body parameters"
// @Success 200 {string} string "Updated book details successfully"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "InternalServerError"
// @Failure default {string} string "DefaultError"
// @Router /books/{id}/edit [put]
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	setContentType(w)
	params := mux.Vars(r)

	var originalBook models.Book
	row := db.QueryRow("SELECT title, description FROM books WHERE id = ?", params["id"])

	err := row.Scan(&originalBook.Title, &originalBook.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			// w.WriteHeader(http.StatusNotFound)
			// fmt.Fprintf(w, "Book with id %v not found", params["id"])
			http.Error(w, "Book with id "+params["id"]+" not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(err.Error()))
		log.Fatal(err)
	}

	// Check if a request is made without providing any request body data :
	if r.Body == http.NoBody {
		fmt.Fprintf(w, "Please provide a request body.")
		return
	}

	// To decode the request body values into a new instance of a book (updated book):
	var updatedBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	CheckError(err)

	if updatedBook.Description == "" {
		updatedBook.Description = originalBook.Description
	}

	// 'RowsAffected' method is generally used with 'update', 'insert, 'replace' statements since it will tell the rows actually changed or affected.

	// Using a prepared statement instead :
	stmt, err := db.Prepare("UPDATE books SET title = ? , description = ? , updated_at = ? WHERE id = ?")
	CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(updatedBook.Title, updatedBook.Description, time.Now(), params["id"])
	// // For checking if 'id' passed as parameter exists :
	// if number, _ := res.RowsAffected(); number == 0 {
	// 	fmt.Fprintf(w, "The id %v does not exists.", params["id"])
	// 	return
	// }
	// CheckError(err)
	if err != nil {
		http.Error(w, "Could not execute database query. ERROR - "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if updatedBook.Author.Name != "" {
	// 	stmt, err = db.Prepare("UPDATE authors SET name = ? , updated_at = ? WHERE book_id = ?")
	// 	CheckError(err)
	// 	defer stmt.Close()

	// 	_, err = stmt.Exec(updatedBook.Author.Name, time.Now(), params["id"])
	// 	if err != nil {
	// 		http.Error(w, "Could not execute database query. ERROR - "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated data successfully.")
}

// func PatchBook(w http.ResponseWriter, r *http.Request) {
// 	// w.Header().Set("Content-Type", "application/json")
// 	setContentType(&w)
// 	params := mux.Vars(r)
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		w.WriteHeader(400)
// 		w.Write([]byte("ID cannot be converted to integer"))
// 		return
// 	}

// 	pB := &Books

// 	for _, book := range *pB { // Dereferencing // Ranging over the dereferenced value (i.e. the Books slice)
// 		if book.ID == id {
// 			updatedBook := book
// 			fmt.Println(updatedBook)
// 			_ = json.NewDecoder(r.Body).Decode(&updatedBook) // why Decode prefers a pointer / address value - CHECK?
// 			json.NewEncoder(w).Encode(updatedBook)
// 			return
// 		}
// 	}
// 	json.NewEncoder(w).Encode(fmt.Sprintf("Cannot update (patch) the book as book with ID %v is not found", id))
// }

// DeleteBook godoc
// @Summary To delete a book with the given id
// @Description Delete Book
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {string} string "Deleted book successfully"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "InternalServerError"
// @Router /books/{id} [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	params := mux.Vars(r)

	// res, err := db.Exec("UPDATE authors SET book_id = NULL where book_id = ?", params["id"])
	// if err != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// if number, _ := res.RowsAffected(); number == 0 {
	// 	// fmt.Fprintf(w, "Book with id %v not found.", params["id"])
	// 	json.NewEncoder(w).Encode("Book with ID " + params["id"] + " not found")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// Do not need to return rows on executing delete statement so uing 'db.Exec' method.
	res, err := db.Exec("DELETE FROM books WHERE id = ?", params["id"])
	if err != nil {
		CheckErrorAsJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// CheckError(err)
	if number, _ := res.RowsAffected(); number == 0 {
		// fmt.Fprintf(w, "Book with id %v not found.", params["id"])
		json.NewEncoder(w).Encode("Book with ID " + params["id"] + " not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted data successfully.")
}

// /////////////////////////////////////////////////////
// /////////////////////// AUTHORS //////////////////////

// GetAuthors godoc
// @Summary To get details of all authors
// @Description Get Authors
// @Tags authors
// @Produce json
// @Success 200 {array} models.Author "All authors"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /authors [get]
func GetAuthors(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type","application/json")

	// If written here, then when running the api with swagger UI then will get this line - 'Can't parse JSON. Raw result :' printed in the starting.
	// setContentType(&w)
	var Authors []models.Author

	rows, err := db.Query("SELECT id, name FROM authors")
	CheckError(err)
	defer rows.Close()
	if !(rows.Next()) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Empty table"))
		return
	}

	for rows.Next() {
		var author models.Author
		err = rows.Scan(&author.ID, &author.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Authors = append(Authors, author)
		// fmt.Println(Authors)
	}
	setContentType(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Authors)
}

// CreateAuthor godoc
// @Summary To create a new author
// @Description Create a new author
// @Tags authors
// @Accept json
// @Produce json
// @Param author body models.Author true "New author request body parameters"
// @Success 200 {object} models.Author "New author created"
// @Failure 500 {string} string "InternalServerError"
// @Failure default {string} string "DefaultError"
// @Router /authors/create [post]
func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	var author models.Author
	err := json.NewDecoder(r.Body).Decode(&author)
	// CheckErrorAsJSON(w, http.StatusInternalServerError, err.Error())
	// If not declared inside an if statement then will give an 'invalid memory / nil pointer dereference' error while testing the function.
	// This is because, if the 'Decode' method does not produces any error then 'err = nil' and this value will be passed to 'CheckErrorAsJSON' function like - CheckErrorAsJSON(w, http.InternalServerError, nil).
	// So, that's why the error is 'invalid memory add / nil pointer dereference'. - I THINK
	if err != nil {
		CheckErrorAsJSON(w, err.Error(), http.StatusInternalServerError)
	}
	// if err := author.Validate(); err != nil {
	// 	CheckErrorAsJSON(w, http.StatusInternalServerError, "Validation error "+err.Error())
	// 	fmt.Fprintln(w, err.Error())
	// 	return
	// }

	stmt, err := db.Prepare("INSERT INTO authors (name) VALUES (?)")
	if err != nil {
		CheckErrorAsJSON(w, err.Error(), http.StatusInternalServerError)
	}
	defer stmt.Close()

	res, err := stmt.Exec(author.Name)
	if err != nil {
		CheckErrorAsJSON(w, "Not able to insert a book in the database. ERROR - "+err.Error(), http.StatusInternalServerError)
	}

	if number, _ := res.RowsAffected(); number != 0 {
		// fmt.Fprintf(w, "New book inserted into the database")
		id, _ := res.LastInsertId()
		author.ID = int(id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(author)
	}
}

// GetAuthor godoc
// @Summary To get details of an author with it's id.
// @Description Get Author by ID
// @Tags authors
// @Produces json
// @Param id path int true "Author ID"
// @Success 200 {object} models.Author "Author`` stored at the given ID"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /authors/{id} [get]
func GetAuthor(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	params := mux.Vars(r)
	row := db.QueryRow("SELECT id, name FROM authors WHERE id = ?", params["id"])

	var author models.Author
	err := row.Scan(&author.ID, &author.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book with id "+params["id"]+" not found.", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

// UpdateAuthor godoc
// @Summary To update an author with the given id
// @Description Update Author
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param author body models.Author true "Update author request body parameters"
// @Success 200 {string} string "Updated author details successfully"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "InternalServerError"
// @Failure default {string} string "DefaultError"
// @Router /authors/{id}/edit [put]
func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Check if a request is made without providing any request body data :
	if r.ContentLength <= 0 {
		fmt.Fprintf(w, "Please provide a request body.")
		return
	}
	setContentType(w)
	var updatedAuthor models.Author
	err := json.NewDecoder(r.Body).Decode(&updatedAuthor)
	// CheckError(err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, e := db.Exec("UPDATE authors SET name = ? , updated_at = ? WHERE id = ?", updatedAuthor.Name, time.Now(), params["id"])
	// CheckError(e) // Will not check for if 'id' value passed in the query does not exists in the database.
	if e != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if no, _ := res.RowsAffected(); no == 0 {
		fmt.Fprintf(w, "The id %v does not exists.", params["id"])
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated data successfully.")
}

// DeleteAuthor godoc
// @Summary To delete an author with the given id
// @Description Delete Author
// @Tags authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {string} string "Deleted author successfully"
// @Failure 404 {string} string "StatusNotFound"
// @Failure 500 {string} string "InternalServerError"
// @Router /authors/{id} [delete]
func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Do not need to return rows on executing delete statement so uing 'db.Exec' method.
	res, err := db.Exec("DELETE FROM authors WHERE id = ?", params["id"])
	CheckError(err)
	if number, _ := res.RowsAffected(); number == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Book with id %v not found.", params["id"])
		return
	}
	setContentType(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted data successfully")
}

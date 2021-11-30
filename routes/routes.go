package routes

import (
	"github.com/gorilla/mux"
	"github.com/hsbRadiant/restapi/controllers"
)

// func f(rw http.ResponseWriter, r *http.Request) {
// 	rw.WriteHeader(http.StatusCreated)
// }

// func understand() {
// 	// unserstand 'http.Handle', 'http.HandlerFunc' with the following example :-
// 	// http.Handle("/handlerFunction", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 	// 	rw.WriteHeader(http.StatusCreated)
// 	// }))

// 	// http.Handle("/handlerFunction", f)                   // Error - CHECK it out
// 	http.Handle("/handlerFunction", http.HandlerFunc(f)) // Not an error.
// 	a := http.HandlerFunc(f)
// 	var (
// 		w http.ResponseWriter
// 		r *http.Request
// 	)
// 	a.ServeHTTP(w, r)                         //
// 	mux.NewRouter().NewRoute().HandlerFunc(f) // NewRoute method registers an empty route; HandlerFunc(f) sets a handler function (f) or allows 'f' to be used as a handler function.
// 	// 'mux.NewRouter' is helping to register a route to a matcher (handler function) - I THINK.
// 	// Without mux.NewRouter this could have been done like http.Handle()
// }

func Setup(router *mux.Router) {
	// MATHCING ROUTES TO HANDLER FUNCTIONS / ENDPOINTS :

	// (ROUTES FOR USER)
	// Route for registering a new user :
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	// Route for logging in an existing user :
	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	// Route for logging out a user :
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// (ROUTES FOR BOOK(S))
	// Route for getting all books :
	// swagger:route GET /books books getBooks
	//
	//
	//
	// Get all books.
	//
	// Produces:
	// - application/json
	//
	// Responses:
	//  default: error
	//  200: models.Book
	//  422: validationError
	router.HandleFunc("/books", controllers.AuthorizeClientRequest(controllers.GetBooks)).Methods("GET") // combining more than one matcher in a single route.

	// Route for getting a specific book (with the help of id) :
	// swagger:route GET /books/{id} books getBook
	// ---
	// description: Get a book.
	// produces:
	//   - application/json
	// responses:
	//   200:
	//     description: Get details of a book.
	//     schema:
	//	     $ref: "#/definitions/models.Book"
	//     default:
	//       description: Error response
	//       schema:
	//	       $ref: "#/definitions/error"
	router.HandleFunc("/books/{id}", controllers.AuthorizeClientRequest(controllers.GetBook)).Methods("GET")
	// router.HandleFunc("/books/{id}", controllers.GetBook).Methods("GET")

	// Route for creating a new book :
	router.HandleFunc("/books/create", controllers.CreateBook).Methods("POST")
	// router.HandleFunc("/books/create", controllers.AuthorizeClientRequest(controllers.CreateBook)).Methods("POST")
	// Route for updating a book :
	router.HandleFunc("/books/{id}/edit", controllers.UpdateBook).Methods("PUT")
	// router.HandleFunc("/books/create", controllers.AuthorizeClientRequest(controllers.UpdateBook)).Methods("PUT")
	// Route for updating (patch) a book :
	// router.HandleFunc("/books/{id}/edit", controllers.PatchBook).Methods("PATCH")
	// Route for deleting a book :
	router.HandleFunc("/books/{id}", controllers.DeleteBook).Methods("DELETE")
	// router.HandleFunc("/books/create", controllers.AuthorizeClientRequest(controllers.DeleteBook)).Methods("DELETE")

	// (ROUTES FOR AUTHOR)
	// Route for getting all authors :
	router.HandleFunc("/authors", controllers.GetAuthors).Methods("GET")
	// Route for getting a specific author (with the help of id) :
	router.HandleFunc("/authors/{id}", controllers.GetAuthor).Methods("GET")
	// Route for creating a new author :
	router.HandleFunc("/authors/create", controllers.CreateAuthor).Methods("POST")
	// Route for updating an author :
	router.HandleFunc("/authors/{id}/edit", controllers.UpdateAuthor).Methods("PUT")
	// Route for updating (patch) an author :
	// router.HandleFunc("/authors/{id}/edit", patchAuthor).Methods("PATCH")
	// Route for deleting an author :
	router.HandleFunc("/authors/{id}", controllers.DeleteAuthor).Methods("DELETE")
}

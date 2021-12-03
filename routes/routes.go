package routes

import (
	"github.com/gorilla/mux"
	"github.com/hsbRadiant/restapi/controllers"
)

// MATHCING ROUTES TO HANDLER FUNCTIONS / ENDPOINTS :
func Setup(router *mux.Router) {

	// (ROUTES FOR USER)
	// Route for registering a new user :
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	// Route for logging in an existing user :
	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	// Route for logging out a user :
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// (ROUTES FOR BOOK) (AUTHORIZED) :
	// Route for getting all books :
	router.HandleFunc("/books", controllers.AuthorizeClientRequest(controllers.GetBooks)).Methods("GET") // combining more than one matcher in a single route.
	// Route for getting a specific book (with the help of id) :
	router.HandleFunc("/books/{id}", controllers.AuthorizeClientRequest(controllers.GetBook)).Methods("GET")
	// Route for creating a new book :
	router.HandleFunc("/books/create", controllers.AuthorizeClientRequest(controllers.CreateBook)).Methods("POST")
	// Route for updating a book :
	router.HandleFunc("/books/{id}/edit}", controllers.AuthorizeClientRequest(controllers.UpdateBook)).Methods("PUT")
	// Route for deleting a book :
	router.HandleFunc("/books/{id}", controllers.AuthorizeClientRequest(controllers.DeleteBook)).Methods("DELETE")

	// (ROUTES FOR BOOK - UNAUTHORIZED) :
	// router.HandleFunc("/books", controllers.GetBooks).Methods("GET") // combining more than one matcher in a single route.
	// router.HandleFunc("/books/{id}", controllers.GetBook).Methods("GET")
	// router.HandleFunc("/books/create", controllers.CreateBook).Methods("POST")
	// router.HandleFunc("/books/{id}/edit", controllers.UpdateBook).Methods("PUT")
	// Route for updating (patch) a book :
	// router.HandleFunc("/books/{id}/edit", controllers.PatchBook).Methods("PATCH")
	// router.HandleFunc("/books/{id}", controllers.DeleteBook).Methods("DELETE")

	// (ROUTES FOR AUTHOR)
	// Route for getting all authors :
	router.HandleFunc("/authors", controllers.GetAuthors).Methods("GET")
	// Route for getting a specific author (with the help of id) :
	router.HandleFunc("/authors/{id}", controllers.GetAuthor).Methods("GET")
	// Route for creating a new author :
	router.HandleFunc("/authors/create", controllers.CreateAuthor).Methods("POST")
	// Route for updating an author :
	router.HandleFunc("/authors/{id}/edit", controllers.UpdateAuthor).Methods("PUT")
	// Route for deleting an author :
	router.HandleFunc("/authors/{id}", controllers.DeleteAuthor).Methods("DELETE")
	// Route for updating (patch) an author :
	// router.HandleFunc("/authors/{id}/edit", patchAuthor).Methods("PATCH")

}

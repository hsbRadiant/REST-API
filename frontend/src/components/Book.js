import { Link } from "react-router-dom";

const Book = ({ book }) => {
    const deleteBook = async () => {
        const response = await fetch("http://localhost:8080/books/" + book.id, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include"
        })
        if (response.ok) {
            alert("Deleted book successfully.")
            window.location.reload(false) // to refresh the page.
        } else {
            alert(response.text())
        }
    }

    return (
        <>
            <th scope="row">{book.id}</th>
            <td>{book.title}</td>
            <td>{book.description}</td>
            <td>TEST AUTHOR</td>
            {/* <td>{book.Author.Name}</td> */}
            <td><Link to={"/books/" + book.id + "/edit"}><button>Update</button></Link></td>
            <td><Link to={"/books"}><button onClick={deleteBook}>Delete</button></Link></td>
            {/* Link to will only establish a link so on clicking the link the next page opened is being defined b link to. Don't confuse it with <Route> types. */}
        </>
    )
}

export default Book;
import { Link } from "react-router-dom";

const Author = ({ author }) => {
    const deleteAuthor = async () => {
        const response = await fetch("http://localhost:8080/authors/" + author.id, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include"
        })
        if (response.ok) {
            alert("Deleted author successfully.")
            window.location.reload(false) // to refresh the page.
        } else {
            alert(response.text())
        }
    }

    return (
        <>
            <th scope="row">{author.id}</th>
            <td>{author.name}</td>
            <td><Link to={"/authors/" + author.id + "/edit"}><button>Update</button></Link></td>
            <td><Link to={"/authors"}><button onClick={deleteAuthor}>Delete</button></Link></td>
            {/* Link to will only establish a link so on clicking the link the next page opened is being defined b link to. Don't confuse it with <Route> types. */}
        </>
    )
}

export default Author;
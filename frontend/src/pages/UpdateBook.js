import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";

const UpdateBook = () => {
    let params = useParams(); // 'useParams' will help to get the request parameters. READ?
    // Acc. to official docs -> 'useParams' will return an object of key/value pairs of the URL parameters. Use it to access 'match.params' of the current route.

    const [title, setTitle] = useState('')
    const [description, setDesciption] = useState('')
    const [authorName, setAuthorName] = useState('')

    const navigate = useNavigate();

    const submit = async (e) => {
        e.preventDefault();
        console.log(title, description, authorName)

        const response = await fetch("http://localhost:8080/books/" + params.id + "/edit", {
            method: 'PUT',
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({ title, description, authorName })
        });
        if (title === "" || authorName === "") {
            return alert("Please provide a title / authorName for the book")
        }
        else if (response.ok) {
            // const content = await response.json()
            alert("Updated book details succesfully.")
            navigate('/books')
            // TO DO - Any confirmation that book is updated.
        } else {
            return alert("NOT OK" + response.text)
        }

    }
    return (
        <>
            <form onSubmit={submit}>
                <center><h1 className="h3 mb-3 fw-normal">UPDATE BOOK DETAILS</h1></center>
                <input type="text" className="form-control" id="floatingInput" placeholder="Title of the book" onChange={e => setTitle(e.target.value)} />
                <input type="text" className="form-control" id="floatingInput" placeholder="Description of the book" onChange={e => setDesciption(e.target.value)} />
                <input type="text" className="form-control" id="floatingPassword" placeholder="Author" onChange={e => setAuthorName(e.target.value)} />

                < button className="w-100 btn btn-lg btn-primary" type="submit">UPDATE</button>
            </form>
        </>
    );
};

export default UpdateBook;
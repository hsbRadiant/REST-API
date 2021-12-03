import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const CreateBook = () => {

    const [title, setTitle] = useState('')
    const [description, setDesciption] = useState('')
    const [authorName, setAuthorName] = useState('')

    const navigate = useNavigate();

    const submit = async (e) => {
        e.preventDefault();

        // console.log(title, description, authorName)

        const response = await fetch("http://localhost:8080/books/create", {
            method: 'POST',
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({ title, description, authorName })
        });

        if (title === "" || authorName === "") {
            return alert("Title / Author name of the book should not be empty. Please provide a title / authorName for the book.")
        }
        else if (response.ok) {
            const content = await response.json()
            alert("BOOK \"" + content.title + "\" CREATED SUCCESFULLY.")
            navigate('/books')
            // TO DO - Any confirmation that book is updated.
        } else {
            return alert(response.body)
        }

    }
    return (
        <>
            <form onSubmit={submit}>
                <center><h1 className="h3 mb-3 fw-normal">CREATE A BOOK</h1></center>
                <input type="text" className="form-control" id="floatingInput" placeholder="Title of the book" onChange={e => setTitle(e.target.value)} />
                <input type="text" className="form-control" id="floatingInput" placeholder="Description of the book" onChange={e => setDesciption(e.target.value)} />
                <input type="text" className="form-control" id="floatingPassword" placeholder="Author" onChange={e => setAuthorName(e.target.value)} />

                < button className="w-100 btn btn-lg btn-primary" type="submit">SAVE</button>
            </form>
        </>
    );
};

export default CreateBook;
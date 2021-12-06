import React, { useRef, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";

const UpdateAuthor = () => {
    let params = useParams(); // 'useParams' will help to get the request parameters. READ?
    // Acc. to official docs -> 'useParams' will return an object of key/value pairs of the URL parameters. Use it to access 'match.params' of the current route.

    const [name, setName] = useState('')
    const nameRef = useRef()

    const navigate = useNavigate();

    const submit = async (e) => {
        e.preventDefault();
        console.log(name)

        const response = await fetch("http://localhost:8080/authors/" + params.id + "/edit", {
            method: 'PUT',
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({ name })
        });
        if (name === "") {
            return alert("Please provide a name for the author")
        }
        else if (response.ok) {
            // const content = await response.json()
            alert("Updated author details successfully.")
            navigate('/authors')
            // TO DO - Any confirmation that book is updated.
        } else {
            nameRef.current.value = null
            return alert("NOT OK" + response.text)
        }

    }
    return (
        <>
            <form onSubmit={submit}>
                <center><h1 className="h3 mb-3 fw-normal">UPDATE AUTHOR DETAILS</h1></center>
                <input type="text" className="form-control" id="floatingInput" placeholder="Name" ref={nameRef} onChange={e => setName(e.target.value)} />

                < button className="w-100 btn btn-lg btn-primary" type="submit">UPDATE</button>
            </form>
        </>
    );
};

export default UpdateAuthor;
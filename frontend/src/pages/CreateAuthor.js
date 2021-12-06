import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const CreateAuthor = () => {
    const [name, setName] = useState('')
    const navigate = useNavigate();

    const submit = async (e) => {
        e.preventDefault();

        const response = await fetch("http://localhost:8080/authors/create", {
            method: 'POST',
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({ name })
        });

        if (name === "") {
            return alert("Author name should not be empty. Please provide a name for the author.")
        }
        else if (response.ok) {
            const content = await response.json()
            alert("Author \"" + content.name + "\" created successfully.")
            navigate('/authors')
        } else {
            return alert(response.text)
        }
    }
    return (
        <>
            <form onSubmit={submit}>
                <center><h1 className="h3 mb-3 fw-normal">CREATE AN AUTHOR</h1></center>
                <input type="text" className="form-control" id="floatingInput" placeholder="Name of the author" onChange={e => setName(e.target.value)} />
                < button className="w-100 btn btn-lg btn-primary" type="submit">SAVE</button>
            </form>
        </>
    );
};

export default CreateAuthor;
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Author from "../components/Author"


const Authors = () => {
    const [authors, setAuthors] = useState([])
    useEffect(async function fetchAuthors() {
        const response = await fetch("http://localhost:8080/authors", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json"
            },
            credentials: "include"
        })
        if (response.ok) {
            setAuthors(await response.json())
        } else {
            alert(response.text())
        }
        fetchAuthors();
    }, []);


    // const deleteAuthor = async (id) => {
    //     console.log("AUTHOR ID -", id)

    //     const response = await fetch("http://localhost:8080/authors/" + id, {
    //         method: "DELETE",
    //         headers: { "Content-Type": "application/json" },
    //         credentials: "include"
    //     })
    //     if (response.ok) {
    //         alert("AUTHOR DELETED SUCCESSFULLY.")
    //         return
    //     } else {
    //         alert(response.text())
    //         return
    //     }
    // }

    return (
        <>
            <div>
                <center>
                    <Link to="/authors/create">< button type="button" className="btn btn-info">CREATE NEW AUTHOR</button></Link>
                    <p> </p>
                    <h1> All authors </h1>
                </center>

                <table className="table">
                    <thead>
                        <tr>
                            <th scope="col">ID</th>
                            <th scope="col">NAME</th>
                            <th scope="col">UPDATE AUTHOR</th>
                            <th scope="col">DELETE AUTHOR</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            authors && authors.map(author => { return <tr><Author author={author} /></tr> })
                        }
                    </tbody>
                </table>
            </div>
        </>
    )
}
export default Authors;
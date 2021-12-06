import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Book from "../components/Book"

const Books = () => {
    let display;
    const [status, setStatus] = useState(0)
    const [books, setBooks] = useState([])
    // As 'Books' will not have an form so using useeffect to display all books :
    useEffect(() => {
        // 'useEffect' does not accepts aync functions.
        (
            async () => {
                // TO DO - a logic to render different html if login or otherwise :
                const response = await fetch("http://localhost:8080/books", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                        "Accept": "application/json"
                    },
                    credentials: "include"
                });
                if (response.ok) {
                    // console.log(response.status)
                    setStatus(response.status)
                    const content = await response.json();
                    setBooks(content)
                } else if (response.status === 401) {
                    setStatus(response.status)
                } else {
                    // console.log(response.status)
                    setStatus(response.status)
                    // console.log(status)
                }
            }
        )();
    }, []);

    if (status === 200) {
        display = (
            <div>
                <center>
                    <Link to="/books/create">< button type="button" className="btn btn-info">CREATE NEW BOOK</button></Link>
                    <span> </span><Link to="/authors">< button type="button" className="btn btn-info">ALL AUTHORS</button></Link>
                </center>

                <h1>All books</h1>
                <table className="table">
                    <thead>
                        <tr>
                            <th scope="col">ID</th>
                            <th scope="col">TITLE</th>
                            <th scope="col">DESCRIPTION</th>
                            <th scope="col">AUTHOR</th>
                            <th scope="col">UPDATE BOOK</th>
                            <th scope="col">DELETE BOOK</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            books && books.map(book => { return <tr><Book book={book} /></tr> })
                            // var obj = JSON.parse(book) // Throws a cross-orogin error whenever a cookie is parsed. READ ?
                        }
                    </tbody>
                </table>
            </div>
        )
    } else if (status === 401) {
        display = (
            <div>
                <center>
                    <h2>NOT LOGGED IN !!! PLEASE LOGIN TO PROCEED</h2>
                    <Link to="/login"><button className="btn btn-info"><h3>Login</h3></button></Link>
                </center>
            </div>
        )
    }
    return (<>{display}</>)
};

export default Books;
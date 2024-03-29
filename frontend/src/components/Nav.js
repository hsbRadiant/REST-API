import React, { useState } from "react";
import { Link } from "react-router-dom";

const Nav = ({ loggedEmail }) => {
    // const [text, setText] = useState("LogOut")

    const logout = async () => {
        await fetch('http://localhost:8080/logout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include'
        });
        // setText("Login")
    }

    let menu;
    if (loggedEmail !== "") {
        menu = (
            <ul className="navbar-nav me-auto mb-2 mb-md-0">
                <li className="nav-item">
                    <Link to="/login" className="nav-link active" aria-current="page" onClick={logout}>LogOut</Link>
                </li>
                <li className="nav-item">
                    <Link to="/books" className="nav-link active" aria-current="page">Books</Link>
                </li>
            </ul>
        )
    } else {
        menu = (
            <ul className="navbar-nav me-auto mb-2 mb-md-0">
                {/* <li className="nav-item">
                    <Link to="/login" className="nav-link active" aria-current="page" onClick={logout}>LogOut</Link>
                </li> */}
                <li className="nav-item">
                    <Link to="/books" className="nav-link active" aria-current="page">Books</Link>
                </li>
                <li className="nav-item">
                    <Link to="/login" className="nav-link active" aria-current="page">Login</Link>
                </li>
                <li className="nav-item">
                    <Link to="/register" className="nav-link active" aria-current="page">Register</Link>
                </li>
            </ul>
        )
    }

    return (
        <>
            <nav className="navbar navbar-expand-md navbar-dark bg-dark mb-4">
                <div className="container-fluid">
                    <Link to="/" className="navbar-brand">Home</Link>
                    {/* <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
          </button> */}
                    <div className="collapse navbar-collapse" id="navbarCollapse">
                        {menu}
                    </div>
                </div>
            </nav>
            {/* <p>{isLoggedIn}</p> */}
        </>
    );
};

export default Nav;
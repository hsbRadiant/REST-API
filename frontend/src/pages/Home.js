import React from "react";
import { Link } from "react-router-dom";

const Home = ({ email }) => {
    let display, displayLoggedIn;
    display = (
        <>
            < h1 > WELCOME TO BOOKS SHELF</h1>
            <p>REACT + GO REST APIS</p>
            <h3></h3>
            <Link to="/login"><button className="btn btn-info"><h4>Login to view all books</h4></button></Link>
            <p> </p>
        </>
    )
    displayLoggedIn = (
        <>
            < h1 > WELCOME "{email}" TO BOOKS SHELF</h1>
            <p>REACT + GO REST APIS</p>
            <h3></h3>
            <Link to="/books"><button className="btn btn-info"><h4>View all books</h4></button></Link>
            <p> </p>
        </>
    )
    return (
        <div>
            <center>
                {email !== "" ? displayLoggedIn : display}

                {/* <p>REACT + GO REST APIS</p>
                <h3></h3>
                <Link to="/login"><button className="btn btn-info"><h4>Login to view books</h4></button></Link>
                <p> </p> */}
                <Link to="/register"><button className="btn btn-info"><h4>New User?- Click on Register.</h4></button></Link>
            </center >
        </div >
    );
};

export default Home;
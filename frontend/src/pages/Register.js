import React, { useState } from "react";
import { useNavigate } from "react-router";

const Register = () => {
    // For every input need to get a variable. First is the variable and second is a natural function that changes the variable.
    const [name, setName] = useState('') // (Empty string intially / intial state - I THINK)
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const [error, setError] = useState('')
    // const [Redirect, setRedirect] = useState(false) // intially the redirect is false so to remain on the registration page until register button is clicked. 
    let navigate = useNavigate()
    // Create a 'submit' function - It accepts an 'event' :
    const submit = async (e) => { // In 'TypeScript' need to define the type of the event too.
        // On submitting a form, it usually refreshes the page. To prevent it using 'PreventDefault()' method.
        e.preventDefault();

        // Sending request to the server using fetch :
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                "Accept": "text/plain"
            },
            body: JSON.stringify({ name, email, password })
        }).catch(err => {
            JSON.stringify({ err })
            setError(err.text())
            // console.log(err.text())
        });
        // const content = await response.json();
        // console.log(content); // Registration successful.
        // setRedirect(true) // after the fetch has been done then set 'redirect' to 'true'.
        if (response.ok) {
            navigate("/login")
        } else {
            // const content = JSON.stringify(response)
            alert("ERROR")
            // window.location.reload(false)
            // navigate("/register")
        }
    };

    return (
        <form onSubmit={submit}>
            <center><h1 className="h3 mb-3 fw-normal">PLEASE REGISTER</h1></center>
            {/* Using the function (setName) (to set a name - I THINK) whenever the input is changed */}
            <input type="name" className="form-control" id="floatingInput" placeholder="Full Name" onChange={e => setName(e.target.value)} />


            <input type="email" className="form-control" id="floatingInput" placeholder="name@example.com" onChange={e => setEmail(e.target.value)} />
            {/* <label for="floatingInput">Email address</label> */}

            <input type="password" className="form-control" id="floatingPassword" placeholder="Password" onChange={e => setPassword(e.target.value)} />
            {/* <label for="floatingPassword">Password</label> */}

            <button className="w-100 btn btn-lg btn-primary" type="submit">Register</button>
        </form>
    );
};

export default Register;
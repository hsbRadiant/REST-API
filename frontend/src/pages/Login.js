import React, { useRef, useState } from 'react'
import { useNavigate } from 'react-router';


const Login = ({ setLoggedIn }) => {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const emailRef = useRef();
    const passRef = useRef();

    let navigate = useNavigate();

    const submit = async (e) => {
        e.preventDefault();

        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
            credentials: 'include' // to include the cookie to be sent when a client logs in.
        });
        const content = await response.text()
        console.log(content)
        if (response.ok) {
            setLoggedIn("true", email)
            alert("WELCOME! \"" + email + ". Login successful.\n View books.")
            navigate('/books')
        } else if (email === "" || password === "") {
            alert("Email / Password field cannot be empty. RETRY.")
        } else {
            alert("Login unsuccessful (Email / Password not correct) : " + content + "TRY AGAIN.")
            emailRef.current.value = null
            passRef.current.value = null
        }
    };
    return (
        <>
            <form onSubmit={submit}>
                <center><h1 className="h3 mb-3 fw-normal">PLEASE SIGNIN</h1></center>
                <input type="email" className="form-control" id="floatingInput" placeholder="name@example.com" ref={emailRef} onChange={e => setEmail(e.target.value)} />
                {/* <label for="floatingInput">Email address</label> */}

                <input type="password" className="form-control" id="floatingPassword" placeholder="Password" ref={passRef} onChange={e => setPassword(e.target.value)} />
                {/* <label for="floatingPassword">Password</label> */}

                < button className="w-100 btn btn-lg btn-primary" type="submit">Sign in</button>
            </form>
        </>
    );
};

export default Login;
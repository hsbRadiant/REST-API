import React, { useState } from 'react';
import './App.css';
import Home from "./pages/Home"
import Books from './pages/Books';
import Authors from './pages/Authors';

import CreateAuthor from './pages/CreateAuthor'
import UpdateAuthor from './pages/UpdateAuthor'
import CreateBook from './pages/CreateBook'
import UpdateBook from './pages/UpdateBook'

import Login from "./pages/Login"
import Register from "./pages/Register"

import Nav from "./components/Nav"
import { BrowserRouter, Route, Routes } from 'react-router-dom';


function App() {
  // To check whether logged in or not :
  const [isLoggedIn, setIsLoggedIn] = useState("false")
  const [email, setEmail] = useState("")
  const callbackFunction = (loggedIn, email) => {
    setIsLoggedIn(loggedIn);
    setEmail(email);
    // console.log("CHECK -", loggedIn, email)
  }

  return (
    <div className="App">
      <BrowserRouter>
        <Nav loggedEmail={email} />

        <Routes>
          <Route path="/" element={<Home email={email} />}></Route>
          <Route path="/books" element={<Books />}></Route>
          <Route path="/authors" element={<Authors />}></Route>
        </Routes>
        <main className="form-signin">
          <Routes>
            <Route path="/authors/create" element={<CreateAuthor />}></Route>
            <Route path="/authors/:id/edit" element={<UpdateAuthor />}></Route>
            <Route path="/books/create" element={<CreateBook />}></Route>
            <Route path="/books/:id/edit" element={<UpdateBook />}></Route>
            <Route path="/register" element={<Register></Register>}></Route>
            <Route path="/login" element={<Login setLoggedIn={callbackFunction} ></Login>}></Route>
          </Routes>
        </main>

      </BrowserRouter >
    </div >
  );
}

export default App;

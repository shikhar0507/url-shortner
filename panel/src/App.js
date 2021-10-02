import React, { useState,useEffect } from 'react';
import {BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  useHistory,
  Link
} from 'react-router-dom';
import Navbar from './navbar';
import {Login,Signup} from './Auth/'
import Campaign from './Campaign';
import Home from './Home';
import './App.scss'

const App = () => {
  const [isAuthenticated,setAuth] = useState(false);
  const history = useHistory()

  useEffect(()=>{
    fetch("http://localhost:8080/auth",{
      method:'GET',
      body:null,
      credentials:"include",
    }).then(res=>{
      if(res.ok) {
        return res.json()
      }
    }).then(response=>{
      setAuth(response.Authenticated)
    }).catch(err=>{
      console.error(err)
    })
  },[])

  const handleLogout = () => {
    console.log("logout")

    fetch("http://localhost:8080/logout",{
      body:null,
      headers:{
        'Content-Type':'application/json'
      },
      credentials:"include",
      method:"DELETE"
    }).then(res=>{
      if(res.ok) return res.json()
    }).then(()=>{
      setAuth(false)
      history.push("/")
    }).catch(err=>{
      console.error(err)
    })
    
  }

  const onLoginSuccess = (his)  => {
    setAuth(true)
  }

  return (
    <div>
    <Router>
      <Navbar isAuthenticated={isAuthenticated} handleLogout={handleLogout}></Navbar>
      <div className="app">
        <Switch>
          <Route path="/login">
            {isAuthenticated ? <Home></Home> : <Login handleLogin={onLoginSuccess}></Login>}
          </Route>
          <Route path="/signup">
          {isAuthenticated ? <Home></Home> : <Signup></Signup>}
          </Route>
          <Route path="/campaign">
            {isAuthenticated ? <Campaign></Campaign> : <Redirect to="/login"></Redirect> }
          </Route>
          <Route path="/">
            <Home auth={isAuthenticated}></Home>
          </Route>
        </Switch>
      </div>
    </Router>
    </div>
  )
}



const CampaginLink = () => {
  return (
    <Link to="/campaign">
      <button className="button is-primary">
      Create campaign
      </button>
    </Link>
  )
}


// export default App;
export {App,CampaginLink}


import React from 'react';
import {
    Link,
    withRouter,
} from "react-router-dom";
import { useAuth } from '../App';


 function Navbar(props) {
    const {user,setUser} = useAuth()
    const logout = () => {
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
        setUser(false)
      }).catch(err=>{
        console.error(err)
      })
      
    }

    return (
        <nav className="navbar is-fixed-top is-dark" role="navigation" aria-label="main navigation">
        <div className="navbar-brand">
          <a className="navbar-item" href="./">
              <span className="app-name">URL Shortner</span>
          </a>
        </div>
        <div className="navbar-menu is-active">
        <div className="navbar-start">
            <div className="navbar-item">
                    <Link to="/campaign">
                        <button className="button is-primary">
                        Create campaign
                        </button>

                    </Link>
             
            </div>
        </div>
          <div className="navbar-end">
            <div className="navbar-item">
              <div className="buttons">
                { user?
                        <a className="button is-primary" onClick={logout}>
                        Logout
                        </a>
                 : 
                <div>
                    <Link to="/signup">
                        <div className="button is-primary">
                            Sign up
                        </div>
                    </Link>
                    <Link to="/login">
                        <div className="button is-light ml-2">
                            Login
                        </div>
                    </Link>
                </div>
                }
              </div>
            </div>
          </div>
        </div>
      </nav>
    )
}


export default withRouter(Navbar);
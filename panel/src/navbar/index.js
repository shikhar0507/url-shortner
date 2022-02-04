import React from 'react';
import {
    Link,
    withRouter,
} from "react-router-dom";
import { useAuth } from '../App';
import "./index.scss";

 function Navbar(props) {
     console.log("navbar",props)
     
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
           <div className="navbar-item">
             <button className="button menu-button" onClick={()=>props.isOpen ? props.setOpen(false):  props.setOpen(true)}>
	       <span className="icon is-small">
                 <i className="fas fa-bars"></i>
               </span>
	     </button>
        </div>
          <a className="navbar-item" href="./">
	    
	    <span className="app-name">URL Shortner</span>
          </a>
            {user &&
            <div className="navbar-item logout-item">
               <Link to="/logout" onClick={logout} className="is-hidden-desktop-only is-hidden-widescreen-only">
                                <div className="button is-light ml-2">
                            Logout
                                </div>
               </Link>
            
            </div>
            }
        </div>
        <div className="navbar-menu is-active is-hidden-mobile is-hidden-tablet-only">
        <div className="navbar-start">
        </div>
          <div className="navbar-end">
            <div className="navbar-item">
              <div className="buttons">
                { !user?
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
                  :    <Link to="/logout" onClick={logout}>
                                <div className="button is-light ml-2">
                            Logout
                                </div>
               </Link>
            
          
                }
              </div>
            </div>
          </div>
        </div>
      </nav>
    )
}


export default withRouter(Navbar);

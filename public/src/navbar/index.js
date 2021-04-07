import React from 'react';
import {
    Link,
    withRouter,
} from "react-router-dom";

class Navbar extends React.Component {
    constructor(props) {
        super(props)
        console.log(this.props)
        // this.state = {
        //     auth : this.props.authState
        // }
        this.logout = this.logout.bind(this)

    }
    
    logout() {
       
        this.props.handleLogout(this.props.history)
        
    }

    render() {
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
                {this.props.authState?
                        <a className="button is-primary" onClick={this.logout}>
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
}

export default withRouter(Navbar);
import React from 'react';
import {BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  Link
} from 'react-router-dom';
import Navbar from './navbar';
import Auth from './Auth';
import Campaign from './Campaign';
import Home from './Home';
import './App.scss';
class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      auth : isLoggedIn()
    }
    this.handleLogout = this.handleLogout.bind(this)
    this.handleLogin = this.handleLogin.bind(this)

  } 
  handleLogout(his) {
    console.log("logout")
    localStorage.removeItem("cookie")
    this.setState({auth:false})
    his.push('/')

  }
  handleLogin(his) {
    console.log('login')
    localStorage.setItem('cookie','true')
    his.push('/')
    this.setState({auth:true})
  }
  render(){
    console.log("sending auth state as",this.state.auth)
    return (
      <div>
      <Router>
        <Navbar authState={this.state.auth} handleLogout={this.handleLogout}></Navbar>
        <div className="app">
          <Switch>
            <Route path="/login">
              <Auth type="login" key="login" handleLogin={this.handleLogin}></Auth>
            </Route>
            <Route path="/signup">
              <Auth type="signup" key="signup"></Auth>
            </Route>
            <Route path="/campaign">
              {this.state.auth ? <Campaign></Campaign> : <Redirect to="/login"></Redirect> }
            </Route>
            <Route path="/">
              <Home auth={this.state.auth}></Home>
            </Route>
          </Switch>
        </div>
      </Router>
      </div>
    )
  }
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

const isLoggedIn = () => {
  return localStorage.getItem("cookie")
  const cookie = document.cookie;
  const split = cookie.split(";");
  let auth = false
  split.forEach(item=>{
    const pair = item.split("=")
    const key = pair[0];
    const value = pair[1];
    if(key === "sessionId" && value) {
      auth = true;
    }
  })
  return true
}

// export default App;
export {App,CampaginLink}


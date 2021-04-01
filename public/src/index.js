import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
import Auth from './Auth';
import './App.scss';
import reportWebVitals from './reportWebVitals';
import Navbar from './navbar';

const Footer = () => {
  return (
    <footer className="footer is-light">
      <div className="content has-text-centered">
        <p>
          <strong>Made with React,bulma,GO,postgresql</strong>
        </p>
      </div>
    </footer>
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
    // this.pro

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
            <Route path="/campaign">campaign</Route>
            <Route path="/">Home</Route>
          </Switch>
        </div>
      </Router>
      </div>
    )
  }
} 

ReactDOM.render(
  <React.StrictMode>
      <App></App>
      <Footer></Footer>
  </React.StrictMode>,
  document.getElementById('root')
);



// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

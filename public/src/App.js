import './App.scss';
import Navbar from './navbar';
import React from 'react';

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoggedIn: false,
    }
    this.onLogin = this.onLogin.bind(this)
    this.onSignup = this.onSignup.bind(this)
    this.onLogout = this.onLogout.bind(this)

  }
  componentDidMount() {
    console.log("component mounted")
  }
  
  onLogin(bool) {
      console.log("user logged in",bool)
      localStorage.setItem("cookie","true")
      this.setState({isLoggedIn:true})
      
      // window.location.href = window.location.origin+"/"
  }
  onSignup(bool) {
    console.log("user signed up",bool)
    console.log(this.history)
    
    // window.history.replaceState(null,null,window.location.origin+"/login")
    // window.location.href = window.location.origin+"/login"
  }
  onLogout() {
    console.log("user logged out")
    window.location.href = window.location.origin+"/"
    this.setState({isLoggedIn:false})
  }

  render(){
    console.log(this.state.isLoggedIn)
    return(
      <div>
        <Router>
        <Navbar isLoggedIn={this.state.isLoggedIn} onLogout={this.onLogout}></Navbar>
        <div className="app">
          <Switch>
            <Route path="/login">
              <Login onLogin={this.onLogin} type="login" key="login"></Login>
            </Route>
            <Route path="/signup">
              <Login onSignup={this.onSignup} type="signup" key="signup"></Login>
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






export default App;

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
import './App.scss'


class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      auth : false
    }
    this.handleLogout = this.handleLogout.bind(this)
    this.handleLogin = this.handleLogin.bind(this)

  } 
  componentDidMount() {
    fetch("http://localhost:8080/auth",{
      method:'GET',
      body:null,
      credentials:"include",
    }).then(res=>{
      if(res.ok) {
        return res.json()
      }
    }).then(response=>{
      this.setState({auth:response.Authenticated})
    }).catch(err=>{
      console.error(err)
    })
  }
  handleLogout(his) {
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
      this.setState({auth:false})
      his.push('/')
    }).catch(err=>{
      console.error(err)
    })
    
  }

  handleLogin(his) {

    console.log('login')
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
              {this.state.auth ? <Home></Home> : <Auth type="login" key="login" handleLogin={this.handleLogin}></Auth>}
            </Route>
            <Route path="/signup">
            {this.state.auth ? <Home></Home> : <Auth type="signup" key="signup" handleLogin={this.handleLogin}></Auth>}
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


// export default App;
export {App,CampaginLink}


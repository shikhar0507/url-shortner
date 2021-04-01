import React from 'react';
import {withRouter} from 'react-router-dom';

 class Auth extends React.Component {
    constructor(props) {
      super(props)
      this.handleSubmit = this.handleSubmit.bind(this)
      this.setUsername = this.setUsername.bind(this)
      this.setPassword = this.setPassword.bind(this)
      console.log(props)
      this.state = {
        psswd:"",
        username:"",
        error:"",
        isActive:false,
        authType : this.props.type
      }
    }
    componentDidMount() {
    //   this.setState({authType:this.props.type})
    }
    setUsername(e) {
      this.setState({username:e.target.value})
    }
    setPassword(e) {
      this.setState({psswd:e.target.value})
    }
  
    handleSubmit(e) {
      if(!this.state.username) {
        return this.setState({error:'Please enter a username'})
      }
      if(!this.state.psswd) {
        return this.setState({error:'Please enter a password'})
      }
      this.setState({isActive:true})
      if(this.state.authType === "login") {
        //   localStorage.setItem('cookie','true')
          this.props.handleLogin(this.props.history)
    } 
    else {
        this.props.history.push('/login')
    }
      return
      fetch(`http://localhost:8080/${this.state.authType}`,{
        method:'POST',
        body:JSON.stringify({
         username:this.state.username,
         psswd:this.state.psswd 
        }),
        headers: {
          'Content-Type':'application/json'
        }
      }).then(res=>{
        return res.json()
      }).then(response=>{
        if(response.status >= 226) {
          this.setState({error:response.message,isActive:false})
          return
        }    
        this.state.authType === "login" ?  this.props.onLogin(true) : this.props.onSignup(true)
      }).catch(error=>{
        this.setState({error:error.message,isActive:false})
      })
  
    }
  
    render() {
      return (
      <div className="login-form">
      <h1 className="title is-3">{this.state.authType  === "login" ? "Login" : "Signup"}</h1>
      <div className="field">
        <label className="label">Username</label>
         <div className="control">
           <input className="input is-success" type="text" placeholder="Enter username" onChange={this.setUsername} value={this.state.username}></input>
         </div>
      </div>
      <div className="field">
        <label className="label">Password</label>
         <div className="control">
           <input className="input" type="password" placeholder="Enter password" onChange={this.setPassword} value={this.state.psswd}></input>
         </div>
      </div>
      {this.state.error ?  <div className="has-text-danger pt-2 pb-2">{this.state.error}</div> : ''}
     
      <button className={"button is-primary "+(this.state.isActive ? 'is-loading'  : '')} onClick={this.handleSubmit}>login</button>
    </div>
    )
    }
  }
  
export default withRouter(Auth)
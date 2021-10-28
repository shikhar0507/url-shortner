import React from 'react';
import { useState } from 'react';
import {useHistory} from 'react-router-dom';
import { useAuth } from '../App';

const Login = (props) => {
  const {setUser} = useAuth()
  const [username,setUsername] = useState('')
  const [password,setPassword] = useState('')
  const [error,setError] = useState('')
  const [active,setActive] = useState(false)

  const handleSubmit = (e) => {
    if(!username) {
      return setError('Please enter a username')
    }
    if(!password) {
      return setError('Please enter a password')
    }
    setActive(true)
    setError('')
    
    fetch(`http://localhost:8080/login-user`,{
      method:'POST',
      body:JSON.stringify({
       username:username,
       psswd:password 
      }),
      credentials: 'include',
      headers: {
        'Content-Type':'application/json'
      }
    }).then(res=>{
      return res.json()
    }).then(response=>{
      if(response.status >= 226) {
        setError(response.message)
        setActive(false)
        return
      }    
      setUser(true)
    }).catch(error=>{
      setError(error.message)
      setActive(false)
    })

  }

  return (
    <div className="login-form">
      <h1 className="title is-3">Login</h1>
      <div className="field">
        <label className="label">Username</label>
        <div className="control">
          <input className="input is-success" type="text" placeholder="Enter username" onChange={(e)=>{setUsername(e.target.value)}} value={username}></input>
        </div>
      </div>
      <div className="field">
        <label className="label">Password</label>
        <div className="control">
          <input className="input" type="password" placeholder="Enter password" onChange={(e)=>{setPassword(e.target.value)}} value={password}></input>
        </div>
      </div>
      {error ?  <div className="has-text-danger pt-2 pb-2">{error}</div> : ''}
    
      <button className={"button is-primary "+(active ? 'is-loading'  : '')} onClick={handleSubmit}>Login</button>
  </div>
  )
}
const Signup = () => {
  const history = useHistory()
  const [username,setUsername] = useState('')
  const [password,setPassword] = useState('')
  const [error,setError] = useState('')
  const [active,setActive] = useState(false)

  const handleSubmit = (e) => {
    if(!username) {
      return setError('Please enter a username')
    }
    if(!password) {
      return setError('Please enter a password')
    }
    setActive(true)
    setError('')
    
    fetch(`http://localhost:8080/signup-user`,{
      method:'POST',
      body:JSON.stringify({
       username:username,
       psswd:password 
      }),
      credentials: 'include',
      headers: {
        'Content-Type':'application/json'
      }
    }).then(res=>{
      return res.json()
    }).then(response=>{
      if(response.status >= 226) {
        setError(response.message)
        setActive(false)
        return
      }    

      history.push('/login')
        
    }).catch(error=>{
      setError(error.message)
      setActive(false)
    })

  }

  return (
    <div className="login-form">
      <h1 className="title is-3">Signup</h1>
      <div className="field">
        <label className="label">Username</label>
        <div className="control">
          <input className="input is-success" type="text" placeholder="Enter username" onChange={(e)=>{setUsername(e.target.value)}} value={username}></input>
        </div>
      </div>
      <div className="field">
        <label className="label">Password</label>
        <div className="control">
          <input className="input" type="password" placeholder="Enter password" onChange={(e)=>{setPassword(e.target.value)}} value={password}></input>
        </div>
      </div>
      {error ?  <div className="has-text-danger pt-2 pb-2">{error}</div> : ''}
    
      <button className={"button is-primary "+(active ? 'is-loading'  : '')} onClick={handleSubmit}>Signup</button>
  </div>
  )
}

export {Login,Signup}

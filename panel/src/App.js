import React, { useState,useEffect,useContext,createContext } from 'react';
import {BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  useHistory,
  Link
} from 'react-router-dom';
import Navbar from './navbar';
import {Login,Signup} from './Auth/'
import Home from './Home';
import './App.scss'


const authContext = createContext(null)

const useAuth = () => {
  return useContext(authContext)
}

const ProvideAuth = ({children}) => {
  const auth = useProvideAuth();
  console.log(auth)
  return (
    <authContext.Provider value={auth}>
      {children}
    </authContext.Provider>
  )
}

const useProvideAuth = () =>{
  const [user,setUser] = useState(false);
  const [isLoading,setIsLoading] = useState(true)
  const fetchAuth = async () => {
    try {
      const authResponse = await fetch("http://localhost:8080/auth",{
        method:'GET',
        body:null,
        credentials:"include",
      })
      const authState = await authResponse.json()
      return authState
    }catch(e){
     return e
    }
  }
  fetchAuth().then(authState=>{
    console.log(authState)
    setUser(authState.Authenticated)
    setIsLoading(false)
  }).catch(console.error)
  console.log(user)
  return {
    user,
    setUser,
    isLoading
  }
}

const App = () => {
  return (
  <ProvideAuth>
    <Router>
      <Navbar></Navbar>
      <div className="app">
        <Switch>
          <PublicRoute path="/login">
            <Login></Login>
          </PublicRoute>
          <PublicRoute path="/signup">
            <Signup></Signup>
          </PublicRoute>
          <PrivateRoute path="/">
           <Home></Home>
          </PrivateRoute>
        </Switch>
      </div>
    </Router>
</ProvideAuth>
  )
}


const PrivateRoute = ({children,...rest}) => {
 const {user,isLoading} = useAuth()
 return (

   <Route {...rest}>
    {isLoading ? (<div>Loading...</div>) : user ? (children) : (<Redirect to="/login"></Redirect>)}
   </Route>
 )
}

const PublicRoute = ({children,...rest}) => {
  const {user} = useAuth()
  return (
    <Route {...rest}>
      {user ? (<Redirect to="/"></Redirect>) : (children)}
    </Route>
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
export {App,CampaginLink,useAuth}


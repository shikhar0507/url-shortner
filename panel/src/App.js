import React, { useState,useEffect,useContext,createContext } from 'react';
import {BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';
import Navbar from './navbar';
import Sidebar from './sidebar';
import {Login,Signup} from './Auth/'
import Home from './Home';
import './App.scss'
import LinkCreate from './linkcreate';
import Campaign from './Campaign';
import {Links}  from './links';

const authContext = createContext(null)

const useAuth = () => {
  return useContext(authContext)
}

const ProvideAuth = ({children}) => {
  const auth = useProvideAuth();
  return (
    <authContext.Provider value={auth}>
      {auth.isLoading ? (<div>Loading...</div>) : (children)}
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
  useEffect(()=>{
    console.log('sending auth request...')
    fetchAuth().then(authState=>{
        setUser(authState.Authenticated)
        setIsLoading(false)
    }).catch(console.error)
  },[])

  return {
    user,
    setUser,
    isLoading
  }
}

const App = () => {
    console.log('app')
    const [openMenu,setOpenMenu] = useState(false)


    useEffect(()=>{
      	console.log("menu changed",openMenu)

    },[openMenu])
  return (
  <ProvideAuth>
  
    <Router>
	  <Navbar setOpen={setOpenMenu} isOpen={openMenu}></Navbar>
	  <Sidebar isOpen={openMenu} setOpen={setOpenMenu}></Sidebar>
	  <div className="app">
	  {openMenu &&  <div className="backdrop-root" onClick={()=>setOpenMenu(false)}> </div> }
       <div className="content-column container">
          <Switch>
            <PublicRoute path="/login">
              <Login></Login>
            </PublicRoute>
            <PublicRoute path="/signup">
              <Signup></Signup>
            </PublicRoute>
            <PublicRoute path="/create-link">
              <LinkCreate></LinkCreate>
            </PublicRoute>
            <PrivateRoute path="/links">
              <Links></Links>
            </PrivateRoute>
            <PrivateRoute path="/campaigns">
              <Campaign></Campaign>
            </PrivateRoute>
            <PrivateRoute path="/">
              <Home></Home>
            </PrivateRoute>
          </Switch>
        </div>
      </div>
    </Router>
</ProvideAuth>
  )
}


const PrivateRoute = ({children,...rest}) => {
 const {user,isLoading} = useAuth()
 console.log(user,isLoading,children)
 return (

   <Route {...rest}>
    {isLoading ? (<div>Loading...</div>) : user ? (children) : (<Redirect to="/login"></Redirect>)}
   </Route>
 )
}

const PublicRoute = ({children,...rest}) => {
  const {user,isLoading} = useAuth()
  const {path} = {...rest}
  return (
    <Route {...rest}>

      {isLoading ? (<div>Loading...</div>) : path === "/create-link" ? (children) : user ? (<Redirect to="/"></Redirect>) : (children)}
    </Route>
  )
}



// export default App;
export {App,useAuth}


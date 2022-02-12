import React, { useState,useEffect,useContext,createContext } from 'react';
import {BrowserRouter as Router,
  Navigate,
  Route,
  Routes,
  useLocation
} from 'react-router-dom';
import Navbar from './navbar';
import Sidebar from './sidebar';
import {Login,Signup} from './Auth/'
import Home from './Home';
import './App.scss'
import LinkCreate from './linkcreate';
import Campaign from './Campaign';
import {LinkDetail, Links, LinkTable}  from './links';

const authContext = createContext(null)

const useAuth = () => {
  return useContext(authContext)
}

const ProvideAuth = ({children}) => {
  const auth = useProvideAuth();
  console.log("provide auth",auth)
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
        console.log("setting user status",authState.Authenticated)
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
          <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/signup" element={<Signup />} />
              <Route path="/create-link" element={<LinkCreate />} />

              <Route path="links" element={<RequireAuth> <Links /></RequireAuth>}>
                <Route index element={<LinkTable />} />
                <Route path=":linkId" element={<LinkDetail />} />
              </Route>
              <Route path="/campaigns" element={<RequireAuth> <Campaign /></RequireAuth>} />
              <Route path="/" element={<RequireAuth> <Home /></RequireAuth>} />
              <Route
      path="*"
      element={
        <main style={{ padding: "1rem" }}>
          <p>There's nothing here!</p>
        </main>
      }
    />
          </Routes>
        </div>
      </div>
    </Router>
</ProvideAuth>
  )
}

const RequireAuth = ({children}) => {
  const {user,isLoading} = useAuth();
  const location = useLocation();
  console.log("navigate",user,isLoading)
  return isLoading ? (<div>Loading...</div>) :  user  ? children : <Navigate to="/login" replace state={{path:location.pathname}}></Navigate>
}



// export default App;
export {App,useAuth}


import React,{useState} from 'react';
import {CampaginLink} from '../App';
require('./index.scss');


const Home = (props) => {

    return (
        <div className="home">
                <div className="home-url">
                    <URLShortner></URLShortner>
                    {!props.isAuthenticated ? <CampaginLink></CampaginLink> : ''}
                </div>
                {props.isAuthenticated ? 
                    <div className="dashboard">
                       <nav className="level">
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Total clicks</p>
                                   <p className="title"></p>
                               </div>
                           </div>
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Most popular campaign</p>
                                   <p className="title"></p>
                               </div>
                           </div>
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Most used device</p>
                                   
                               </div>
                           </div>
                       </nav>
                    </div>
                : ''}
            </div>
    )
    
}

const URLShortner = () => {
    const [inputUrl,setInputUrl] = useState(null)
    const handleUrl = (e) => {
        setInputUrl(e.target.value)
    }
    const [shortenUrl,setShortenUrl] = useState(null)
    const [error,setError] = useState('');
    const [active,setActive] = useState(false);
    
    const createLink = () => {
        if(!inputUrl) {
            setError({error:'Enter url'})
            return
        }

        setActive(true)
        setError('')
        fetch("http://localhost:8080/links/",{
            method:'POST',
            headers:{
                'Content-Type':'application/json'
            },
            credentials:"include",
            body:JSON.stringify({longUrl:inputUrl,tag:'Work',notFoundUrl:'https://google.com',password:'xanadu'})
        }).then(res=>{
            return res.json()
        }).then(response=>{
            setShortenUrl(response.LongUrl)
        }).catch(error=>{
            setError(error.message)
        })
    } 
    return(
        <div className="url-card has-text-centered">
                <div className="is-size-4 has-text-weight-semibold">Shorten link</div>
                <div className="field mt-2">
                    <div className="control">
                        <input className="input" placeholder="Enter url" onChange={handleUrl} required></input>
                        <button className={"button is-primary ml-2"+(active ?'is-loading' :'')} onClick={createLink}>Submit</button>
                    </div>
                    {error ? <div className="error has-danger-text mt-1">{error}</div> :''}
                </div>
                <div className="result mt-2 is-success">
                    <a className="title has-text-success is-5" href={shortenUrl}>
                        {shortenUrl}
                    </a>
                </div>
            </div>
    )
}




export default Home
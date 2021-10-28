import React,{useState} from 'react';
import {CampaginLink, useAuth} from '../App';
require('./index.scss');


const Home = (props) => {
    const auth = useAuth()
    return (
        <div className="home">
                <div className="home-url">
                    <URLShortner></URLShortner>
                    {!auth.user ? <CampaginLink></CampaginLink> : ''}
                </div>
            </div>
    )
    
}


const URLShortner = () => {
    const [linkAttrs,setLinkAttrs] = useState({
        longUrl:'',
        androidDeepLink:'',
        iosDeepLink:'',
        expiration:{
            time:new Date(),
            expirationUrl:'https://www.youtube.com/watch?v=yydZbVoCbn0'
        },
        notFoundUrl:'',
        password:'',
        tag:''
    })
    const handleInputUrl = (e) => {
        linkAttrs.longUrl = e.target.value
        setLinkAttrs(linkAttrs)
    }
    const handleAndroidLink = (e) => {
        linkAttrs.androidLink = e.target.value
        setLinkAttrs(linkAttrs)
    }
    const handleIosLink = (e) => {
        linkAttrs.iosLink = e.target.value
        setLinkAttrs(linkAttrs)
    }
    const handle404= (e) => {
        linkAttrs.notFoundUrl = e.target.value
        setLinkAttrs(linkAttrs)
    }
    const handlePassword = (e) => {
        linkAttrs.password = e.target.value
        setLinkAttrs(linkAttrs)
    }
    const handleExpirationDate = (e) => {
        console.log("expiration",new Date(e.target.value))
        linkAttrs.expiration.time = new Date(e.target.value)
        setLinkAttrs(linkAttrs)
    }
    const handleTag = (e) => {
        linkAttrs.tag = e.target.value;
        setLinkAttrs(linkAttrs)
    }


    const [shortenUrl,setShortenUrl] = useState(null)
    const [error,setError] = useState('');
    const [active,setActive] = useState(false);
    
    const createLink = () => {
        console.log("asd")
        if(!linkAttrs.longUrl) {
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
            body:JSON.stringify(linkAttrs)
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
                        <input className="input" placeholder="Enter url" onChange={handleInputUrl} required></input>
                        <button className={"button is-primary ml-2"+(active ?'is-loading' :'')} onClick={createLink}>Submit</button>
                    </div>
                    {error ? <div className="error has-danger-text mt-1">{error}</div> :''}
                </div>
                <div className="mt-2">
                    <div className="columns">
                        <div className="column">
                            <div className="field">
                                <label className="label">Tag</label>
                                <div className="control">
                                    <input className="input" type="" placeholder="link tag" onChange={handleTag}></input>
                                </div>
                            </div>
                        </div>
                        <div className="column">
                            <div className="field">
                                <label className="label">Expiration</label>
                                <div className="control">
                                    <input className="input" type="datetime-local" placeholder="Expiration Time" onChange={handleExpirationDate}></input>
                                </div>
                                <p class="help">Accessing this link after the expiration will redirect user to <i>404 page</i> to our website</p>
                            </div>
                        </div>
                        <div className="column">
                            <div className="field">
                                <label className="label">Not found url</label>
                                <div className="control">
                                    <input className="input" type="text" placeholder="Fallback url" onChange={handle404}></input>
                                </div>
                                <p class="help">If requested url is unavailable , then user will be redirected to this page</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="result mt-2 is-success">
                    <a className="title has-text-success is-5" href={shortenUrl} target="_blank">
                        {shortenUrl}
                    </a>
                </div>
            </div>
    )
}




export default Home

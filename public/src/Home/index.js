import React from 'react';
import {CampaginLink} from '../App';
import './index.scss';
class Home extends React.Component{
    constructor(props) {
        super(props)
    }
    render() {
        return (
            <div className="home">
                <div className={"home-url " +(this.props.auth ? 'dashboard' :'')}>
                    <URLShortner></URLShortner>
                    {!this.props.auth ? <CampaginLink></CampaginLink> : ''}
                </div>
            </div>
        )
    }
}


class URLShortner extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            active:false,
            shortenUrl:"",
            url:""
        }
        this.handleUrl = this.handleUrl.bind(this)
        this.shortenUrl = this.shortenUrl.bind(this)
    }
    handleUrl(e) {
        this.setState({url:e.target.value})
    }
    shortenUrl() {
        if(!this.state.url) {
            this.setState({error:'Enter url'})
            return
        }
        if(!isValidURL(this.state.url)) {
            this.setState({error:"Doesn't look like a correct url"})
            return
        }
        this.setState({active:true,error:""})
        fetch("https://httpbin.org/post",{
            method:'POST',
            headers:{
                'Content-Type':'application/json'
            },
            body:JSON.stringify({url:this.state.url})
        }).then(res=>{
            return res.json()
        }).then(response=>{
            response.url = 'https://short/asd'
            this.setState({shortenUrl:response.url,url:''})
        }).catch(error=>{
            this.setState({error:error.message})
        })
    } 
    render() {
        return(
            <div className="url-card has-text-centered">
                <div className="field">
                    <div className="control">
                        <input className="input" placeholder="Enter url" onChange={this.handleUrl} required></input>
                        <button className={"button is-primary ml-2"+(this.state.active ?'is-loading' :'')} onClick={this.shortenUrl}>Submit</button>
                    </div>
                    {this.state.error ? <div className="error has-danger-text mt-1">{this.state.error}</div> :''}
                </div>
                <div className="result mt-2 is-success">
                    <a className="title has-text-success is-5" href={this.state.shortenUrl}>
                        {this.state.shortenUrl}
                    </a>
                </div>
            </div>
        )
    }
}

const isValidURL = (str) => {
    var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
      '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
      '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
      '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
      '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
      '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
    return !!pattern.test(str);
}
export default Home
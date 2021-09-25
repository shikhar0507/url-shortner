import React from 'react';
import {CampaginLink} from '../App';
import './index.scss';
import {useState} from 'react'

const Home = (props) => {
    console.log(props)
    const data = {
        cardData: {
            totalClicks : 300,
            mostClikedCampaign: 'Campaign 3',
            deviceStats : [{
                Mobile:10
            },{
                Desktop:280
            },
            {
                Tablet:10
            }]
        },
        data: [{
            cname:'campaign 1',
            sname:'google',
            mname:'google ads',
            clicks:50,
        },{
            cname:'campaign 2',
            sname:'facebook',
            mname:'fb ads',
            clicks:40,
        },{
            cname:'campaign 3',
            sname:'gmail',
            mname:'email',
            clicks:210,
        }]
    }


    return (
        <div className="home">
                <div className="home-url">
                    <URLShortner></URLShortner>
                    {!props.auth ? <CampaginLink></CampaginLink> : ''}
                </div>
                {props.auth ? 
                    <div className="dashboard">
                       <nav className="level">
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Total clicks</p>
                                   <p className="title">{data.cardData.totalClicks}</p>
                               </div>
                           </div>
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Most popular campaign</p>
                                   <p className="title">{data.cardData.mostClikedCampaign}</p>
                               </div>
                           </div>
                           <div className="level-item has-text-centered">
                               <div>
                                   <p className="heading">Most used device</p>
                                    {/* <BarGraph data={data.cardData.deviceStats}></BarGraph> */}
                               </div>
                           </div>
                       </nav>
                    </div>
                : ''}
            </div>
    )
    
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
        // if(!isValidURL(this.state.url)) {
        //     this.setState({error:"Doesn't look like a correct url"})
        //     return
        // }
        this.setState({active:true,error:""})
        fetch("http://localhost:8080/shorten",{
            method:'POST',
            headers:{
                'Content-Type':'application/json'
            },
            credentials:"include",
            body:JSON.stringify({url:this.state.url})
        }).then(res=>{
            return res.json()
        }).then(response=>{
            console.log(response.url)
            this.setState({shortenUrl:response.url,url:''})
        }).catch(error=>{
            this.setState({error:error.message})
        })
    } 
    render() {
        return(
            <div className="url-card has-text-centered">
                <div className="is-size-4 has-text-weight-semibold">Shorten link</div>
                <div className="field mt-2">
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
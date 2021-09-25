import React from 'react';
import {withRouter} from 'react-router-dom'

class Campaign extends React.Component {
    constructor(props) {
        super(props) 
        this.state = {
            campaign:"",
            source:"",
            medium:"",
            error:"",
            active:false,
            url:"",
            // previewLink : ""
        }
        this.setCampaign = this.setCampaign.bind(this)
        this.setSource = this.setSource.bind(this)
        this.setMedium = this.setMedium.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
        this.setUrl = this.setUrl.bind(this)
    }
    setCampaign(e) {
        this.setState({campaign:e.target.value})
    }
    setSource(e) {
        this.setState({source:e.target.value})
    }
    setMedium(e) {
        this.setState({medium:e.target.value})
    }
    setUrl(e) {
        
        this.setState({url:e.target.value})
    }
    handleSubmit() {
        if(!this.state.campaign) {
            this.setState({error:'Campaign name is required'})
            return
        }
        this.setState({error:"",active:true})
        fetch("http://localhost:8080/campaign",{
            body:JSON.stringify({
                campaign:this.state.campaign,
                source:this.state.source,
                medium:this.state.medium,
                url:this.state.url
            }),
            credentials:"include",
            headers:{
                'Content-Type':'application/json'
            },
            method:"POST"
        }).then(res=>{
            if (res.ok) return res.json()
        }).then(()=>{
            this.props.history.push('/')
        }).catch(err=>{
            console.error(err)
        })
    }

    render() {
        const pl = this.state.url+(this.state.campaign ? "?&campaign="+this.state.campaign :'')+(this.state.source ? "&source="+this.state.source : '') +(this.state.medium ? "&medium="+this.state.medium : '')
        return (
                <div className="login-form">
                    <div className="field">
                        <label className="label">Url</label>
                        <div className="control">
                            <input className="input" type="text" placeholder="Enter url" onChange={this.setUrl} value={this.state.url} required></input>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Campaign name</label>
                        <div className="control">
                            <input className="input is-success" type="text" placeholder="Enter name" onChange={this.setCampaign} value={this.state.campaign}></input>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Source</label>
                        <div className="control">
                            <input className="input" type="text" placeholder="Enter source" onChange={this.setSource} value={this.state.source}></input>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Medium</label>
                        <div className="control">
                            <input className="input" type="text" placeholder="Enter medium" onChange={this.setMedium} value={this.state.medium}></input>
                        </div>
                    </div>
                    <div className="live-preview">
                        <a href={pl}>{pl}</a>
                    </div>
                    {this.state.error ? <div className="error-field has-text-danger pt-2 pb-2">{this.state.error}</div> : ''}
                    <button className={"button is-primary " + (this.state.active ? 'is-loading' :'')} onClick={this.handleSubmit}>Create campaign</button>
                </div>
            )
    }
}


export default withRouter(Campaign);

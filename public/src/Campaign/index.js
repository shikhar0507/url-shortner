import React from 'react';
import './index.scss';
import {withRouter} from 'react-router-dom'
class Campaign extends React.Component {
    constructor(props) {
        super(props) 
        this.state = {
            campaign:"",
            source:"",
            medium:"",
            error:"",
            active:false
        }
        this.setCampaign = this.setCampaign.bind(this)
        this.setSource = this.setSource.bind(this)
        this.setMedium = this.setMedium.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)

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

    handleSubmit() {
        if(!this.state.campaign) {
            this.setState({error:'Campaign name is required'})
            return
        }
        this.setState({error:"",active:true})
        this.props.history.push('/')
    }

    render() {

        return (
                <div className="login-form">
                    <div className="field">
                        <label className="label">Campaign name</label>
                        <div className="control">
                            <input className="input is-success" type="text" placeholder="Enter name" onChange={this.setCampaign} value={this.state.username} required></input>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Source</label>
                        <div className="control">
                            <input className="input" type="text" placeholder="Enter source" onChange={this.setSource} value={this.state.username}></input>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Medium</label>
                        <div className="control">
                            <input className="input" type="text" placeholder="Enter medium" onChange={this.setMedium} value={this.state.username}></input>
                        </div>
                    </div>
                    {this.state.error ? <div className="error-field has-text-danger pt-2 pb-2">{this.state.error}</div> : ''}
                    <button className={"button is-primary " + (this.state.active ? 'is-loading' :'')} onClick={this.handleSubmit}>Create campaign</button>
                </div>
            )
    }
}


export default withRouter(Campaign);

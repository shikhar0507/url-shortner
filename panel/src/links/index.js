import axios from "axios"
import { Fragment, useEffect, useState } from "react"
import './index.scss';
import {Link,useParams,Switch,Route} from 'react-router-dom';
import { useRouteMatch,useLocation } from "react-router";
const Links = () => {

    const {path} = useRouteMatch();

    return (
        
        <Switch>
            <Route path={`${path}/:linkId`}>
                <LinkDetail></LinkDetail>
            </Route>
            <Route exact path={path}>
                <LinkTable></LinkTable>
            </Route>
           
        </Switch>
    )

}

const LinkTable = () => {
    const [topBrowser, setTopBrowser] = useState("-")
    const [topDevice, setTopDevice] = useState("-")
    const [topOs, setTopOs] = useState("-")
    const [totalClicks, setTotalClicks] = useState("")
    const [links, setLinks] = useState([])
    const getLinks = (type) => {
        if (!type) return;
        axios.get('http://localhost:8080/links/?sort=' + type, {
            withCredentials: true
        }).then(response => {
            setTopBrowser(response.data.top_browser)
            setTopDevice(response.data.top_device)
            setTopOs(response.data.top_os)
            setTotalClicks(response.data.total_clicks)
            setLinks(response.data.result)
        }).catch(err => {
            console.error(err)
        });
    }
    useEffect(() => {

        getLinks('click_desc')

    }, []);
    return (


        <section>
            <div className="columns">
                <div className="column">
                    <StatCard title="Top Browser" value={topBrowser} />
                </div>
                <div className="column">
                    <StatCard title="Top Device" value={topDevice} />

                </div>
                <div className="column">
                    <StatCard title="Top Os" value={topOs} />
                </div>
                <div className="column">
                    <StatCard title="Total Clicks" value={totalClicks} />
                </div>
            </div>
            <div className="filters mb-4">
                <div className="is-flex">
                    <div className="filter">
                        <div className="is-size-6 has-text-weight-bold has-text-grey has-text-centered">Time created</div>
                        <div className="select is-rounded">
                            <select onChange={(e) => getLinks(e.target.value)}>
                                <option value="Default" disabled selected>Default</option>
                                <option value="time_latest">Latest</option>
                                <option value="time_oldest">Oldest</option>
                            </select>
                        </div>
                    </div>
                    <div className="filter ml-4">
                        <div className="is-size-6 has-text-weight-bold has-text-grey has-text-centered">Popularity</div>
                        <div className="select is-rounded">
                            <select onChange={(e) => getLinks(e.target.value)}>
                                <option value="clicks_desc">Highest</option>
                                <option value="clicks_asc">Lowest</option>
                            </select>
                        </div>
                    </div>
                </div>
            </div>
            <table className="table links-table">
                <thead>
                    <tr>
                        <th><abbr title="link">Link Id</abbr></th>
                        <th>Name</th>
                        <th><abbr title="Tag">Tag</abbr></th>
                        <th><abbr title="Top browser">Top browser</abbr></th>
                        <th><abbr title="Top os">Top os</abbr></th>
                        <th><abbr title="Top device">Top device</abbr></th>
                        <th><abbr title="total_clicks">Total clicks</abbr></th>
                    </tr>
                </thead>
                <tfoot>
                    {links.map(link => {
                        return <tr key={link.id}>
                            <td>
                                <Link to={{pathname:"/links/" + link.id,state:{
                                    totalClicks :link.total_clicks
                                }}}  >
                                    {link.id}
                                </Link>
                            </td>
                            <td className="ellipsis link-name">{link.name}</td>
                            <td className="ellipsis link-tag">
                                {link.tag ? <span className="tag is-link">{link.tag}</span> : '-'}
                            </td>
                            <td>{link.browser}</td>
                            <td>{link.os}</td>
                            <td>{link.device_type}</td>
                            <td>{link.total_clicks}</td>
                        </tr>
                    })}
                </tfoot>
            </table>
        </section>
    )

}

const LinkDetail = ({totalClicks})=> {
    const {linkId} = useParams();
    const location = useLocation()
    console.log(location.state)
    const [topBrowser,setTopBrowser] = useState("-")
    const [topDevice,setTopDevice] = useState("-")
    const [topOs,setTopOs] = useState("-")
    // const [totalClicks,setTotalClicks] = useState("")
    const [linkDetails,setLinkDetail] = useState({})

    useEffect(()=>{
        axios.get(`http://localhost:8080/links/${linkId}/`,{
            withCredentials:true
        }).then(resp=>{
            const result = resp.data;
            setTopBrowser(result.top_browser[0].name)
            setTopOs(result.top_os[0].name)
            setTopDevice(result.top_device[0].name)
            setLinkDetail({
                ...linkDetails,
               ...result
            })
            console.log(result)
        }).catch(console.error)
    },[])
    return(
        <section>
            <div className="columns">
                        <div className="column">
                            <StatCard title="Top Browser" value={topBrowser}/>
                        </div>
                        <div className="column">
                        <StatCard title="Top Device" value={topDevice}/>

                        </div>
                        <div className="column">
                            <StatCard title="Top Os" value={topOs}/>
                        </div>
                        <div className="column">
                            <StatCard title="Total Clicks" value={location.state.totalClicks}/>
                        </div>
            </div>
            <div className="">{linkDetails.link_name}</div>
            <div className="is-size-4 has-text-weight-bold">Link : http://localhost:8080/{linkId}</div>
            <span className="tag is-primary">{linkDetails.link_tag}</span>
            <p>{linkDetails.link_description}</p>
            <div className="is-size-5 mt-2">
                <a href={linkDetails.long_url}>{linkDetails.long_url}</a>
            </div>
            <div className="buttons are-small">
                <button className="button is-info is-outlined">
                     <span className="icon">
                        <i className="fas fa-edit"></i>
                    </span>
                    <span>Edit</span>
                </button>
                <button className="button is-danger is-outlined">
                    <span className="icon">
                        <i className="fas fa-trash"></i>
                    </span>
                    <span>Delete</span>
                </button>
            </div>
        </section>
    )
}

const StatCard  = ({title,value}) => {
    return (
        <div className="card">
                <div className="card-content pb-4">
                <div className="media-content">
                    <p className="title is-4 has-text-dark">{title}</p>
                </div>
                    <div className="content mt-2">
                    <div className="is-size-5  has-text-info"> {value}</div>
                    </div>
                </div>
        </div>
    )
}

export  {Links,LinkDetail};
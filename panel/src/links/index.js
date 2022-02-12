import axios from "axios"
import {useEffect, useState } from "react"
import './index.scss';
import {Link,useParams,Outlet, useNavigate, useLocation} from 'react-router-dom';


const Links = () => {


    return (
        <div>
            {/* <LinkTable /> */}
            <Outlet />
        </div>

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

        getLinks('time_latest')

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
                                <Link to={"/links/" + link.id} state={link.total_clicks}>
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

const LinkDetail = ()=> {
    const {linkId} = useParams();
    const [linkDetails,setLinkDetail] = useState({
        top_browser:[],
        top_os:[],
        top_device:[],
        top_referrer:[],
        top_countries:[]
    })
    const navigate = useNavigate() 
    const [deleteLoad,setDeleteLoad] = useState(false)
    const {state} = useLocation()

    const deleteUrl = () => {
        setDeleteLoad(true)
        axios.delete(`http://localhost:8080/links/${linkId}`,{
            withCredentials:true
        }).then(resp=>{
            console.log(resp.data)
            navigate("/links")
        }).catch(err=>{
            setDeleteLoad(false)
        })
    };
    
    useEffect(()=>{
        axios.get(`http://localhost:8080/links/${linkId}/`,{
            withCredentials:true
        }).then(resp=>{
            const result = resp.data;
     
            setLinkDetail({
                ...linkDetails,
               ...result
            })
            console.log(result)
        }).catch(console.error)
    },[])
    return(
        <section>
            <div className="is-pulled-left">{linkDetails.link_name}</div>
                <div className="is-pulled-right has-text-centered">
                    <div>Total Clicks</div>
                    <span className="is-size-3 has-text-success has-text-weight-bold">{state || 0}</span>
                </div>
            <div className="pt-5">
                <a className="is-size-4 has-text-weight-bold" href={"http://localhost:8080/"+linkId} target="_blank">Link : http://localhost:8080/{linkId}</a>
            </div>
            <div>
                <span className="tag is-primary">{linkDetails.link_tag}</span>
                <span className="tag is-info ml-2">Created on : {linkDetails.created_on}</span>
            </div>
            <p className="mt-3">{linkDetails.link_description}</p>
            <p className="is-size-5 mt-2">
                <a href={linkDetails.long_url}>{linkDetails.long_url}</a>
            </p>
            <div className="buttons are-small mt-2">
                <button className="button is-info is-outlined">
                     <span className="icon">
                        <i className="fas fa-edit"></i>
                    </span>
                    <span>Edit</span>
                </button>
                <button className={"button is-danger is-outlined "+(deleteLoad ? 'is-loading' : '')} onClick={deleteUrl}>
                    <span className="icon">
                        <i className="fas fa-trash"></i>
                    </span>
                    <span>Delete</span>
                </button>
            </div>
            <div className="columns mt-4 is-multiline">
                <RankCard title='Top Browser' data={linkDetails.top_browser}/>
                <RankCard title='Top OS' data={linkDetails.top_os}/>
                <RankCard title='Top Device' data={linkDetails.top_device}/>
                <RankCard title='Top Referrer' data={linkDetails.top_referrer}/>
                <RankCard title='Top Country' data={linkDetails.top_countries}/>

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

const RankCard = ({title,data}) => {
    let total = 0
    data.forEach(item=>{
        total += item.value
    })
    return (
        <div className="column is-4">
            <div className="card">
            <header className="card-header">
                <p className="card-header-title">
                {title}
                </p>
                <p className="is-pulled-right pt-4 pb-3 pl-4 pr-4 has-text-weight-semibold">
                Clicks
                </p>
            </header>
            <div className="card-content">

                    <ul>
                        {data.map((item,idx)=>{
                            return <li key={{title}+item.name+idx} className="mt-1">
                                <span className="has-text-black">{item.name}</span>
                                <span className="is-pulled-right has-text-black"><span className="has-text-dark mr-2">{Number(100*(Number(item.value)/total)).toFixed(2)}%</span> {item.value}</span>
                            </li>
                        })}
                    </ul>
            </div>
        </div>
        </div>
    )
}

export  {Links,LinkDetail,LinkTable};
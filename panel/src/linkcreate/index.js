import {useEffect,useState,Fragment } from 'react';
import {countries} from './countries';
import {Link} from 'react-router-dom';
import './index.scss';
const QRCode = require('qrcode.react')
const LinkCreate = () => {
  return (
      <URLShortner></URLShortner>
  )
}

// const country_names = require('./countries')

const URLShortner = () => {

  const [linkAttrs, setLinkAttrs] = useState({
    longUrl: '',
    android_deep_link: '',
    play_store_link:'',
    ios_deep_link: '',
    notFoundUrl: '',
    password: '',
    tag: '',
    name: '',
    description: '',
    campaign : {
      name:'',
      source:'',
      medium:'',
      term:'',
      content:'',
      id:''
    },
    expiration: '',
    httpStatus : 308,
    qrCode:false,
    country_block:[],
    country_redirect:[],
    mobile_url:'',
    desktop_url:'',
  });
  const [shortenUrl,setShortenUrl] = useState(false)


  const [activeView,setActiveVIew] = useState('Basic')



  return (
    shortenUrl ? <Success shortUrl={shortenUrl} qrCode={linkAttrs.qrCode}></Success> :
    <div className="url-card">
      <div className="columns is-multiline">
        <div className="column is-2">
          <aside className={"menu"}>
            <ul className="menu-list">
              {['Basic','Android link','Ios link','Campaign','Security','Expiration','HTTP Status','Redirection','Country Block'].map((item,k)=>{
                return <li key={k}><a onClick={()=>{setActiveVIew(item)}} className={(activeView === item ? 'is-active':'')}>{item}</a></li>
              })}
            </ul>
          </aside>
        </div>
        <div className="column is-10">
          <div className="is-size-4 has-text-weight-semibold">Shorten link</div>
          {activeView === "Basic" && <GeneralView setShortenUrl={setShortenUrl} linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs} />}
          {activeView === "Android link" && <AndroidLinking linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "Ios link" && <IosLinking linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "Campaign" && <Campaign  linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs} />}
          {activeView === "Security" && <Security linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "Expiration" && <Expiration linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "HTTP Status" && <HttpRedirect linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "Redirection" && <DeviceRedirect linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
          {activeView === "Country Block" && <CountryBlock linkAttrs={linkAttrs} setLinkAttrs={setLinkAttrs}/>}
        </div>
      </div>
    </div>
  );
};

const Campaign = ({ linkAttrs, setLinkAttrs }) => {

  const setName = (e) => {
    const cm = linkAttrs.campaign
    cm.name = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
  }
  const setMedium = (e) => {
    const cm = linkAttrs.campaign
    cm.medium = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
    
  }
  const setSource = (e) => {
    const cm = linkAttrs.campaign
    cm.source = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
  }

  const setTerm = (e) => {
    const cm = linkAttrs.campaign
    cm.term = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
  }
  const setContent = (e) => {
    const cm = linkAttrs.campaign
    cm.content = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
  }
  const setId = (e) => {
    const cm = linkAttrs.campaign
    cm.id = e.target.value
    setLinkAttrs({
      ...linkAttrs,
      ...cm
    })
  }

  return (
    <section>
      <div className="columns is-multiline">
        <div className="column is-3">
          <div className="field">
            <label className="label is-pulled-left">Campaign Name</label>
            <div className="control">
              <input className="input" type="text" onChange={setName} value={linkAttrs.campaign.name}></input>
              <p className="help">utm_campaign</p>
            </div>
          </div>
        </div>
        <div className="column is-3">
          <div className="field">
            <label className="label is-pulled-left">Campaign Source</label>
            <div className="control">
              <input className="input" type="text" onChange={setSource} value={linkAttrs.campaign.source}></input>
              <p className="help">utm_source</p>

            </div>
          </div>
        </div>
        <div className="column is-3">
          <div className="field">
            <label className="label is-pulled-left">Campaign Medium</label>
            <div className="control">
              <input className="input" type="text" onChange={setMedium} value={linkAttrs.campaign.medium}></input>
              <p className="help">utm_meidum</p>

            </div>
          </div>
        </div> 
        <div className="column is-3">
          <div className="field">
            <label className="label is-pulled-left">Campaign Term</label>
            <div className="control">
              <input className="input" type="text" onChange={setTerm} value={linkAttrs.campaign.term}></input>
              <p className="help">utm_term</p>

            </div>
          </div>
        </div> 
        <div className="column is-6">
          <div className="field">
            <label className="label is-pulled-left">Campaign Content</label>
            <div className="control">
              <input className="input" type="text" onChange={setContent} value={linkAttrs.campaign.content}></input>
              <p className="help">utm_content</p>

            </div>
          </div>
        </div> 
        <div className="column is-2">
          <div className="field">
            <label className="label is-pulled-left">Campaign Id</label>
            <div className="control">
              <input className="input" type="text" onChange={setId} value={linkAttrs.campaign.id}></input>
              <p className="help">utm_id</p>

            </div>
          </div>
        </div> 
      </div>
    </section>
  );

};

const HttpRedirect = ({ linkAttrs, setLinkAttrs }) => {
  const handleStatusSelect = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      httpStatus:Number(e.target.value)
    });
  };
  return (
    <section>
      <span className="is-size-5 has-text-weight-semibold is-pulled-left form-section-label column is-full no-padding">Http status</span>
      <div className="field">
        <div className="control is-flex is-flex-direction-column http-control">
          <label className="radio">
            <input type="radio" name="foobar" value="301" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 301 ? true :false} />
            <span className="is-size-5 ml-2">301 <span className="has-text-weight-semibold is-size-6">Moved Permanently</span></span>
            
          </label>
          <label className="radio">
            <input type="radio" name="foobar"  value="302" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 302 ? true :false}/>
            <span className="is-size-5 ml-2">302 <span className="has-text-weight-semibold is-size-6">Found</span></span>

          </label>
          <label className="radio">
            <input type="radio" name="foobar" value="303" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 303 ? true :false}/>
            <span className="is-size-5 ml-2">303 <span className="has-text-weight-semibold is-size-6">See Other</span></span>
          </label>
          <label className="radio">
            <input type="radio" name="foobar" value="304" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 304 ? true :false}/>
            <span className="is-size-5 ml-2">304 <span className="has-text-weight-semibold is-size-6">Not Modified</span></span>
          </label>
          <label className="radio">
            <input type="radio" name="foobar" value="305" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 305 ? true :false}/>
            <span className="is-size-5 ml-2">305 <span className="has-text-weight-semibold is-size-6">Use Proxy</span></span>
          </label>
          <label className="radio">
            <input type="radio" name="foobar"value="306" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 306 ? true :false}/>
            <span className="is-size-5 ml-2">306 <span className="has-text-weight-semibold is-size-6">Switch Proxy</span></span>
          </label>
          <label className="radio">
            <input type="radio" name="foobar" value="307" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 307 ? true :false}/>
            <span className="is-size-5 ml-2">307 <span className="has-text-weight-semibold is-size-6">Temporary Redirect </span></span>
          </label>
          <label className="radio">
            <input type="radio" name="foobar" value="308" onChange={handleStatusSelect} checked={linkAttrs.httpStatus === 308 ? true :false}/>
            <span className="is-size-5 ml-2">308 <span className="has-text-weight-semibold is-size-6"> Permanent Redirect</span></span>
          </label>
        </div>
      </div>

    </section>
  );
};

const CountrySearch = (props) => {
  // const [countriesList,setCountries] = useState(countries)
  // const [selected,setSelected] = useState([])
  const handleInput = (e) =>{
    const filtered = countries.filter(i=>{
      const idx = props.linkAttrs.country_block.findIndex(j=>j.name===i.name)
      return  i.name.toLowerCase().includes(e.target.value.toLowerCase()) && idx == -1
    })
    props.setCountries(filtered)
  }
  const handleSelect = (country) => {
    props.selectedCountries.push(country)
    const newSelected = [...props.selectedCountries]
    props.setSelected(newSelected);
    props.countriesList.splice(props.countriesList.findIndex(i=>i.code===country.code),1)
    props.setCountries(props.countriesList)
  }
  
  return (
    <Fragment>
       <div className="field">
            <div className="control">
              <input className="input country-search" type="text" placeholder="Seach country" onChange={handleInput} />
            </div>

          </div>
          <div className="mt-2 block">
            <div className="select is-multiple">
              <select multiple size="8" className="country-list">
                {props.countriesList.map(i => {
                  return <option value={i.code} key={i.code} onClick={()=>{handleSelect(i)}}>{i.name}</option>
                })}
              </select>
            </div>
          </div>
    </Fragment>
  )
}

const CountryBlock = ({linkAttrs,setLinkAttrs}) => {
  const [selected,setSelected] = useState(linkAttrs.country_block)
  const filtered = countries.filter(i=>{
    const idx = linkAttrs.country_block.findIndex(j=>j.name===i.name)
    if(idx == -1) {
      return i
    }
  })
  const [countriesList,setCountries] = useState(filtered)


  useEffect(()=>{
    console.log(selected)
    setLinkAttrs({
        ...linkAttrs,
        country_block:[...selected]
      })
  },[selected])

  
  
  const removeSelected = (country,idx) => {
    selected.splice(idx,1)
    setSelected([...selected]);
    countriesList.push(country);
    setCountries([...countriesList])
    
  }
  return (
    <section>
      <div className="columns">
        <div className="column">
          <CountrySearch linkAttrs={linkAttrs} countriesList={countriesList} selectedCountries={selected} setCountries={setCountries} setSelected={setSelected}></CountrySearch>
        </div>
        <div className="column">
          <div className="field is-grouped is-grouped-multiline">
                  {selected.map((i,j)=>{
                    return <div className="control" key={i.code+j}>
                              <div className="tags has-addons">
                                <a className="tag is-link">{i.name}</a>
                                <a className="tag is-delete" onClick={()=>removeSelected(i,j)}></a>
                              </div>
                        </div>
                  })}
          </div>
        </div>
      </div>
    </section>
  )
}


const DeviceRedirect = ({ linkAttrs, setLinkAttrs }) => {
  
  return (
    <section>
      <div className="tabs">
        <ul>
          <li className="is-active"><a>By Device</a></li>
        </ul>
      </div>
      <div className="tab-area">
         <div> 
          <div className="field">
          <label className="label">Mobile/Tablet</label>
          <div className="control">
            <input className="input" type="text" placeholder="Mobile url" value={linkAttrs.mobile_url} onChange={(e)=> setLinkAttrs({
      ...linkAttrs,
      mobile_url:e.target.value
    })}/>
          </div>
        </div>
        <div className="field">
          <label className="label">Desktop</label>
          <div className="control">
            <input className="input" type="text" placeholder="Desktop url" value={linkAttrs.desktop_url} onChange={(e)=> setLinkAttrs({
      ...linkAttrs,
      desktop_url:e.target.value
    })}/>
          </div>
        </div>
        
        </div> 
        
      </div>
    </section>
  );
};

const AndroidLinking = ({ linkAttrs, setLinkAttrs }) => {
  const handleAndroidLink = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      android_deep_link:e.target.value
    });
  };
  const playSotreLink = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      play_store_link:e.target.value
    });
  };
  return (
    <section>
      <div className="field mt-2">
        <label className="label">Play store url</label>
        <div className="control">
          <input className="input" type="text" placeholder="Play store url" value={linkAttrs.play_store_link} onChange={playSotreLink}></input>
        </div>
      </div>
      <div className="field">
        <label className="label">Deep link</label>
        <div className="control">
          <input className="input" type="text" placeholder="Android deep link" value={linkAttrs.android_deep_link} onChange={handleAndroidLink}></input>
        </div>
      </div>

    </section>
  );
};

const IosLinking = ({ linkAttrs, setLinkAttrs }) => {
  const handleIosLink = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      ios_deep_link:e.target.value
    });
  };


  return (
    <section>

      <div className="field">
        <label className="label">App store url</label>
        <div className="control">
          <input className="input" type="text" placeholder="Ios deep link" value={linkAttrs.ios_deep_link} onChange={handleIosLink}></input>
        </div>
      </div>

    </section>
  );
};

const Expiration = ({ linkAttrs, setLinkAttrs }) => {


  const handleExpirationDate = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      expiration:e.target.value
    });
  };

  return (
    <section>
      <div className="">Link is only valid until expiration time. If fallback url is provided , then after expiration user will be redirected to it</div>
      	      		
        <div className="field mt-2">
          <label className="label">Expiration time</label>
          <div className="control">
            <input className="input" type="datetime-local" value={linkAttrs.expiration} onChange={handleExpirationDate}></input>
          </div>
        </div>
     
    </section>
  );
};


const Security = ({ linkAttrs, setLinkAttrs }) => {
  const handlePassword = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      password:e.target.value
    });
  };



  return (
    <section>
        <p>Protect your link by adding a password.</p>
        <div className="field mt-2">
          <label className="label">Password</label>
          <div className="control">
            <input className="input" type="password" placeholder="Enter password" onChange={handlePassword} value={linkAttrs.password}></input>
          </div>
          <p className="help">When users will open the click, they will need to provide the above password</p>
        </div>
    </section>
  );
};

const GeneralView = ({linkAttrs, setLinkAttrs,setShortenUrl })=>{

  const [linkImages, setLinkImage] = useState([]);
  const [active, setActive] = useState(false);
  const [preview,setPreview] = useState(null)
  useEffect(() => {
    console.log(linkAttrs.longUrl)
    if (!linkAttrs.longUrl) return
    fetch("http://localhost:8080/links/opengraph", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: "include",
      body: JSON.stringify({ url: linkAttrs.longUrl })
    }).then(res => {
      return res.json();
    }).then(response => {
      setLinkAttrs({
        ...linkAttrs,
        description: response.Description,
        name: response.Title,
        tag: response.SiteName
      });
      setLinkImage(response.Image)
      if(!response.Title && !response.SiteName && !response.Description ) {
        setPreview(false)
        return
      }
      setPreview(true)
    }).catch(error => {
      setPreview(false)
    });
  }, [linkAttrs.longUrl]);

  const [error, setError] = useState('');

  const handle404 = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      notFoundUrl:e.target.value
    })
  };
  
  const handleTag = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      tag:e.target.value
    })
  };
  const handleName = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      name:e.target.value
    })
  };
  const handleDescription = (e) => {
    setLinkAttrs({
      ...linkAttrs,
      description:e.target.value
    })
  };
  
  const handleInputUrl = (e) => {
    e.preventDefault();
    setLinkAttrs({
      ...linkAttrs,
      longUrl:e.target.value
    })
  };

  const createLink = () => {
    if (!linkAttrs.longUrl) {
      setError({ error: 'Enter url' });
      return;
    }
    setActive(true);
    setError('');
    fetch("http://localhost:8080/links/", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: "include",
      body: JSON.stringify({
        longUrl: linkAttrs.longUrl,
        name: linkAttrs.name,
        description: linkAttrs.description,
        psswd: linkAttrs.password,
        tag: linkAttrs.tag,
        expiration: linkAttrs.expiration,
        campaign:linkAttrs.campaign,
        android_deep_link:linkAttrs.android_deep_link,
        ios_deep_link:linkAttrs.ios_deep_link,
        http_status:linkAttrs.httpStatus,
        play_store_link:linkAttrs.play_store_link,
        not_found_url:linkAttrs.notFoundUrl,
        qr_code:linkAttrs.qrCode,
        country_block:linkAttrs.country_block.map(i=>i.code),
        country_redirect:[{
          country_code:'IN',
          country_url:linkAttrs.longUrlby
        }],
        mobile_url:linkAttrs.mobile_url,
        desktop_url:linkAttrs.desktop_url,
      })
    }).then(res => {
      return res.json();
    }).then(response => {
      setShortenUrl(response.shortUrl);
    }).catch(error => {
      setError(error.message);
    });
  };
  return (
    <section>
      <div className="columns mb-0">
        <div className="column">
            <label className="label">Main URL</label>
            <div className="field mt-2 has-addons">

              <div className="control" style={{width:100+"%"}}>

                <input className="input" placeholder="Enter url" onChange={handleInputUrl} required value={linkAttrs.longUrl}></input>
              </div>
              <div className="control">
                <button className={"button is-primary ml-2" + (active ? 'is-loading' : '')} onClick={createLink}>Submit</button>
              </div>
              {error ? <div className="error has-danger-text mt-1">{error}</div> : ''}
            </div>
            <div className="field">
                <div className="control">
                  <label className="checkbox has-text-left">
                    <input type="checkbox" className="mr-1" checked={linkAttrs.qrCode} onChange={()=>{setLinkAttrs({...linkAttrs,qrCode:linkAttrs.qrCode ? false:true})}}/>
                    Generate Qr 
                  </label>
                </div>
              </div>
        </div>
      <div className="column">
        <div className="field">
          <label className="label">Fallback URL</label>
          <div className="control">
            <input className="input" type="text" placeholder="Fallback url" onChange={handle404} value={linkAttrs.notFoundUrl}></input>
          </div>
          <p className="help is-size-7">If Main URL is unavailable, link will redirect to this page</p>
        </div>
      </div>
      </div>
    
      <div className="is-size-5 has-text-weight-semibold mb-4 pt-2">Link details</div>

      <div className="columns">
        <div className="column is-one-fifth">
          <div className="field">
            <label className="label is-pulled-left">Tag</label>
            <div className="control">
              <input className="input" type="text" placeholder="Enter tag name" onChange={handleTag} value={linkAttrs.tag}></input>
            </div>
          </div>
        </div>
        <div className="column is-one-fifth">
          <div className="field">
            <label className="label is-pulled-left">Name</label>
            <div className="control">
              <input className="input" type="text" placeholder="Enter name for your link" onChange={handleName} value={linkAttrs.name}></input>
            </div>
          </div>
        </div>      
        <div className="column is-three-fifth">

          <div className="field">
            <label className="label is-pulled-left">Description</label>
            <div className="control">
              <input className="input" type="text" placeholder="Enter description"  onChange={handleDescription} value={linkAttrs.description} />
            </div>
          </div>
        </div>
      </div>
      {preview == true ? <Fragment>
        <div className="is-size-5 has-text-weight-bold">Social media preview</div>
      <div className="columns">

        <div className="column is-half">
          <div className="is-size-6 has-text-weight-bold">Facebook</div>

          <div className="card social facebook-social mt-2">
            <div className="card-image">
              {linkImages && linkImages.length ?   <figure className="image is-4by3">
                <img src={linkImages[0].Url} alt="Placeholder image" />
              </figure> : ""}
            
            </div>
            <div className="card-content pt-1 pb-2 pr-2 pl-2">
                <div className="media mb-0">
                
                  <div className="media-content">
                    <p className="subtitle is-6">{linkAttrs.tag}</p>
                    <p className="title is-6 has-text-weight-bold">{linkAttrs.name}</p>
                  </div>
                </div>

                <div className="content is-size-6 pt-1 social-description">
                  {linkAttrs.description}
                </div>
              </div>
          </div>
        
        </div>
        <div className="column is-half">
            <div className="is-size-6 has-text-weight-bold">Twitter</div>
            <div className="card social twitter-social mt-2">
              <div className="card-image">
                {linkImages && linkImages.length ?   <figure className="image is-4by3">
                  <img src={linkImages[0].Url} alt="Placeholder image" />
                </figure> : ""}
              
              </div>
              <div className="card-content pt-1 pb-2 pr-2 pl-2">
                  <div className="media mb-0">
                  
                    <div className="media-content">
                      <p className="subtitle is-7 has-text-grey">{linkAttrs.tag}</p>
                      <p className="title is-6 has-text-weight-normal">{linkAttrs.name}</p>
                    </div>
                  </div>

                  <div className="content is-size-6 pt-1 social-description has-text-grey">
                    {linkAttrs.description}
                  </div>
                </div>
            </div>        
        </div>
      </div>
      </Fragment> : preview == false ?  <p className="is-danger">No Open Graph Meta data found. Refer to <a href="ogp.me/">OGP Protocal to know more</a></p> : ""}
      
    </section>
  );

};

const Success =  ({shortUrl,qrCode}) => {

  const [downloadLink,setDownloadLink] = useState("")
  const downloadQr = () => {
    const canvas = document.getElementById('download-qr');
    console.log(canvas)
    setDownloadLink(canvas.toDataURL("image/jpeg"))
  }
  console.log(qrCode)
  return (
    <section>
        <div className="has-text-centered">
            <div className="mt-4">
                <div className="is-size-3 has-text-success has-text-weight-bold">Link created successfully</div>
                <a className="is-size-4 has-text-weight-bold" href={shortUrl}>{shortUrl}</a>
            </div>
            {qrCode &&
            <Fragment>
              <div className="qr-code-container mt-4">
                <QRCode value={shortUrl} size={200} id="download-qr"></QRCode>
              </div>
              <div>
                <a className="button is-info mt-4 mb-4" download={shortUrl+".jpg"} onClick={downloadQr} href={downloadLink}>Download QR code</a>
              </div>
            </Fragment>
            }
              <Link to="/links">
                <a className="button">Go to links</a>
              </Link>
          </div>
    </section>
  )
}

export default LinkCreate;

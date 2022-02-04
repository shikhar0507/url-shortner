import { useContext, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../App";
import './sidebar.scss';
const Sidebar = (props) => {
    console.log("sidebar props",props.isOpen)
    const {user} = useAuth()
    const isOpen = props.isOpen
    const [selected,setSelected] = useState(user ? 0 : 1);
    const list = [
      {path:'/',name:'Home',key:0},
      {path:'/create-link',name:'Create link',key:1},
      {path:'/links',name:'Links',key:2},
      {path:'/campaign',name:'Campaign',key:3}
    ];
    return (
            <aside className={"menu sidebar "+(isOpen ? "open" : "")}>
        <ul className="menu-list">
          {list.map(item=>{
              return (<SideBarLink item={item}  active={selected} key={item.key} activateTab={(k)=>{setSelected(k);props.setOpen(false)}}></SideBarLink>)  
          })}
        </ul>
      </aside>
    )
}

const SideBarLink = (props) => {

  const {item,active} = props
  return (
    
      <li data-key={item.key}><Link to={item.path} className={active == Number(item.key) ? 'is-active' : ''} onClick={()=>props.activateTab(item.key)}>{item.name}</Link></li>
    
  ) 
}

export default Sidebar

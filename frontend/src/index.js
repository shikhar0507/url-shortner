import React from 'react';
import ReactDOM from 'react-dom';

import reportWebVitals from './reportWebVitals';
import {App} from './App'
const Footer = () => {
  return (
    <footer className="footer is-light">
      <div className="content has-text-centered">
        <p>
          <strong>Made with React,bulma,GO,postgresql</strong>
        </p>
      </div>
    </footer>
  )
}




ReactDOM.render(
  <React.StrictMode>
      <App></App>
      <Footer></Footer>
  </React.StrictMode>,
  document.getElementById('root')
);



// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

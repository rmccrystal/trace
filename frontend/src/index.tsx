import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './app';

import "@blueprintjs/core/lib/css/blueprint.css";
import "./tailwind.css";

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
);

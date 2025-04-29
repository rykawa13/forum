// index.js (или index.tsx)
import React from 'react';
import ReactDOM from 'react-dom/client'; //Изменилось в React 18
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root')); //React 18
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

//Старый синтаксис (до React 18):
//ReactDOM.render(<App />, document.getElementById('root'));
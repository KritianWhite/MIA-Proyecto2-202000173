// import React from 'react';
// import ReactDOM from 'react-dom';
// import Index from './Pages';

// const App = () => {
//   return (
//     <div>
//       <Index/>
//     </div>
//   );
// };

// ReactDOM.render(<App />, document.getElementById('root'));

import React from "react";
import ReactDOM from "react-dom/client";
import reportWebVitals from "./reportWebVitals.js";
import Index from "./Pages/Home.js";
import Principal from "./Pages/Editores.js";
import Login from "./Pages/Login.js";
import Routes from "./Services/Routes.js";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    {/* <Index /> */}
    {/* <Principal /> */}
    {/* <Login /> */}
    <Routes />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

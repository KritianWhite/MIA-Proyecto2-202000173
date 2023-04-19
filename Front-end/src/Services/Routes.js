import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import Home from "../Pages/Home.js";
import Editores from "../Pages/Editores.js";
import Login from "../Pages/Login.js";
import Reportes from "../Pages/Reportes.js";

function Rutas() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/editores" element={<Editores />} />
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<Reportes />} />
      </Routes>
    </Router>
  );
}

export default Rutas;
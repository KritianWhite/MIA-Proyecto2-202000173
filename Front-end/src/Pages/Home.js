import React from "react"
import { useNavigate } from 'react-router-dom';

import "../Styles/Home.css"

function Index() {

  //TODO: Primero creamos la variable navigate para navegar entre las paginas
  const navigate = useNavigate();
  //* Luego creamos la funcion que nos permitira navegar a la pagina de editores
  const irEditores = () => {
    navigate("/editores");
  }
  //* Luego creamos la funcion que nos permitira navegar a la pagina de login
  const irLogin = () => {
    navigate("/login");
  }
  //* Luego creamos la funcion que nos permitira navegar a la pagina de reportes
  const irReportes = () => {
    navigate("/reportes");
  }

  return (
    <>
      <div class="home">
        <div class="container">
          <div class="container-home">
            <button class="custom-btn btn-home" onClick={irEditores}>Editores</button>
            <button class="custom-btn btn-home" onClick={irLogin}>Login</button>
            <button class="custom-btn btn-home" onClick={irReportes}>Reportes</button>
          </div>
        </div>
      </div>
    </>
  );
}

export default Index;

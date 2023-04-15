import React, { useState } from "react";
import "../Styles/Login.css"; // importa tu archivo de estilos CSS aquí
import axios from "axios";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [idParticion, setIdParticion] = useState("");

  const handleIdParticionChange = (event) => {
    const newText = event.target.value; // Obtener el nuevo valor del textarea
    setIdParticion(newText); // Actualizar el estado local con el nuevo valor del textarea
  };

  const handleEmailChange = (event) => {
    const newText = event.target.value; // Obtener el nuevo valor del textarea
    setEmail(newText); // Actualizar el estado local con el nuevo valor del textarea
  };

  const handlePasswordChange = (event) => {
    const newText = event.target.value; // Obtener el nuevo valor del textarea
    setPassword(newText); // Actualizar el estado local con el nuevo valor del textarea
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    // Lógica para enviar los datos de inicio de sesión al servidor y actualizar el objeto entrada
  };
  var Command = "login >user=" + email + " >pwd=" + password + " >id=" + idParticion;
  //console.log(comand);

  const entrada = {
    comando: "",
    idParticion: "",
    idU: 0,
    idG: 0,
    idMoun: " ",
    nombreU: " ",
    login: false,
  };

  const updateUsuario = (user) => {
    entrada.comando = "";
    entrada.idParticion = user.id_mount;
    entrada.idU = user.id_u;
    entrada.idG = user.id_g;
    entrada.idMoun = user.id_mount;
    entrada.nombreU = user.nombre_u;
    entrada.login = user.login;

    sessionStorage.setItem("entrada", JSON.stringify(entrada));
    console.log("Entrada update: ", JSON.parse(sessionStorage.getItem("entrada")));
  }

  const sendEntrada = () => {
    entrada.comando = Command;
    console.log(entrada);
    // Realizar la solicitud POST con Axios
    axios
      .post("http://localhost:8080/Exec", entrada)
      .then((response) => {
        // Manejar la respuesta del servidor si es necesario
        setResponse(response.data.res);
        console.log(response.data.usuario);
        updateUsuario(response.data.usuario);
      })
      .catch((error) => {
        // Manejar errores si los hay
        console.error("Error al enviar JSON a Go:", error);
      });
  };

  return (
    <>
      <div className="login-container">
        <form onSubmit={handleSubmit}>
          <h1>Iniciar sesión</h1>
          <label htmlFor="id">ID Particion:</label>
          <input
            type="text"
            id="idParticion"
            value={idParticion}
            onChange={handleIdParticionChange}
            required
          />
          <label htmlFor="email">Usuario:</label>
          <input
            type="text"
            id="email"
            value={email}
            onChange={handleEmailChange}
            required
          />
          <label htmlFor="password">Contraseña:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={handlePasswordChange}
            required
          />
          <button type="submit" onClick={sendEntrada}>Iniciar sesión</button>
          <button type="submit">Cerrar Sesion</button>
        </form>
      </div>
    </>
  );
}

export default Login;

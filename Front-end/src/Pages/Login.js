import React, { useState } from "react";
import { useNavigate } from 'react-router-dom';
import axios from "axios";

import Login_ from "./Src/Object.js";
import "../Styles/Login.css";

function Login() {
    /* crear un estado local para guardar el valor de los inputs */
    const [idParticion, setIdParticion] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

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

    const limpiarVariables = () => {
        setIdParticion("");
        setEmail("");
        setPassword("");
    };

    //* Ahora creamos la función para enviar los datos de inicio de sesión al servidor
    //* Y actualizar el objeto user para poder realizar el logout

    var logeo = new Login_();
    const navigate = useNavigate();

    const iniciarSesion = (event) => {
        event.preventDefault();
        // Lógica para enviar los datos de inicio de sesión al servidor y actualizar el objeto entrada
        let comando = ["login >id=" + idParticion + " >user=" + email + " >pwd=" + password];
        let datos = logeo.entrada;
        datos.comandos = comando;
        axios.post("http://3.145.14.213:8080/Exec", datos)
            .then((respuesta) => {
                logeo.updateUsuario(respuesta.data.usuario) //* Actualizamos el usuario
                console.log(respuesta)
                if (respuesta.data.usuario.login) {
                    if (respuesta.data.res === "0)YA HAY UNA SESION ACTIVA\n") {
                        limpiarVariables();
                        return alert("Ya existe una sesión activa");
                    } else {
                        limpiarVariables();
                        return navigate("/Editores")
                    }
                } else {
                    limpiarVariables();
                    return alert("Usuario/Particion/Contraseña incorrectos");
                }
            }).catch((err) => {
                console.error("Error:", err);
                return alert("Error al recibir la petición del servidor.");
            });

    };


    return (
        <>
            <div className="login-container">
                <form>
                    <h1>Iniciar sesión</h1>
                    <label htmlFor="id">ID Particion:</label>
                    <input
                        type="text"
                        id="idParticion"
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
                    <button onClick={iniciarSesion} type="submit">Iniciar sesión</button>
                    {/* <button onClick={logout} type="submit">Cerrar Sesion</button> */}
                </form>
            </div>
        </>
    );
}

export default Login;

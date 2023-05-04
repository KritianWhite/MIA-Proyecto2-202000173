import React from "react";
import { useState } from "react";
import { useNavigate } from 'react-router-dom';
import axios from "axios";

import Editor from "../Components/Editor.js";
import Login from "./Src/Object.js";
import "../Styles/Editores.css";

export default function Principal() {

    //* Primeramente mostramos el contenido del archivo en el text area Entrada
    const [archivo, setArchivo] = useState("");
    const [response, setResponse] = useState("");

    const leerArchivo = (event) => {
        const archivo = event.target.files[0];
        const reader = new FileReader();
        reader.onload = (event) => {
            setArchivo(event.target.result);
        };
        reader.readAsText(archivo);
    };

    //* Ahora capturamos cualquier cambio del text area Entrada
    const changeText = (text) => {
        setArchivo(text);
    }

    //TODO: Ahora realizamos las peticiones al servidor
    //* Ahora metemos el contenido del text area Entrada en lista para almacenarlo en el objeto
    let comandos = [];
    comandos = archivo.split("\n");
    //* Eliminamos las posiciones vacias "" del array
    comandos = comandos.map((item) => item.trim()).map((str) => (str !== "" ? str : null)).filter((str) => str !== null);
    //console.log(comandos);

    //* Ahora hacemos la peticion al servidor y recibimos la respuesta
    var logeo = new Login();

    const enviar_Exec = () => {
        const datos = logeo.entrada;
        datos.comandos = comandos;
        console.log("Enviar Json: ", datos)
        axios.post("http://3.145.14.213:8080/Exec", datos)
            .then((respuesta) => {
                setResponse(respuesta.data.res)
                let usuario = respuesta.data.usuario;
                logeo.updateUsuario(usuario) //* Actualizamos el objeto a SessionStorage
                console.log("Respuesta del servidor por Exec:", respuesta.data);
            }).catch((err) => {
                console.error("Error:", err);
                alert("Error al enviar la peticion al servidor.");
            });
    }

    //* Cerrar sesión
    const logout = () => {
        //event.preventDefault();
        const datos = logeo.entrada;
        datos.comandos = ["logout"];

        if (logeo.getLogin() === false) {
            return alert("No hay una sesión activa")
        } else {
            axios.post("http://3.145.14.213:8080/Exec", datos)
                .then((respuesta) => {
                    logeo.updateUsuario(respuesta.data.usuario) //* Actualizamos el usuario
                    if (respuesta.data.usuario.login === true) {
                        return alert("No se pudo cerrar sesión");
                    }
                    return alert("Se cerró sesión correctamente");
                }).catch((err) => {
                    console.error("Error:", err);
                    return alert("Error al recibir la petición del servidor.");
                });

        }

    }


    const navigate = useNavigate();
    const irLogin = () => {
        navigate("/login")
    }

    const botonStyle1 = {
        position: "absolute",
        top: "2%",
        left: "2%",
    };

    const botonStyle2 = {
        position: "absolute",
        top: "2%",
        right: "2%",
    };

    return (
        <>
            <div class="Editores">
                <div class="principal">
                    <div class="editores">
                        <button class="btn btn-light" style={botonStyle1} onClick={irLogin}>ㅤㅤLoginㅤㅤ</button>
                        <button class="btn btn-danger" style={botonStyle2} onClick={logout}>Cerrar sesión</button>
                        <div class="container">
                            <div class="code-editor">
                                <Editor Tittle={"Entrada"}
                                    contenido={archivo}
                                    handlerChange={changeText}
                                />
                                <input
                                    onChange={leerArchivo}
                                    id="input-archivo"
                                    class="form-control form-control-lg"
                                    type="file"
                                />
                            </div>

                            <div class="console">
                                <Editor Tittle={"Salida"}
                                    readOnly={true}
                                    contenido={response}
                                />
                                <button class="custom-btn btn-11" onClick={enviar_Exec}> Ejecutar</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}

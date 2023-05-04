import React from "react";
import { useState } from "react";
import Graphviz from "graphviz-react";
import axios from "axios";

import Login from "./Src/Object.js";
import "../Styles/Reportes.css"


export default function Reportes() {
    const [valor, setValor] = useState("");
    const [dot, setDot] = useState("digraph {\nnode00[label = \"SALE\"];\nnode00 -> node000\nnode003[label = \"EN\"];\nnode00 -> node001\nnode000[label = \"(\"];\nnode00 -> node003\nnode001[label = \"VACAS\"];\nnode00 -> node002\nnode002[label = \")\"];\n}");

    const changeText = (text) => {
        setValor(text.target.value);
    }


    const logeo = new Login();
    const enviar_Exec = () => {
        const datos = logeo.entrada;
        datos.comandos = [valor];
        axios.post("http://3.145.14.213:8080/Exec", datos)
            .then((respuesta) => {
                //console.log(respuesta.data.dot)
                if(respuesta.data.dot === undefined || respuesta.data.dot === null || respuesta.data.dot === ""){
                    return alert("Error al renderizar el reporte.");
                }
                setDot(respuesta.data.dot)
            }).catch((err) => {
                console.error("Error:", err);
                return alert("Error al recibir la petición del servidor.");
            });
    }


    return (
        <>

            <div className="reportes">
                <h1 class="titulo-report">REPORTES</h1>
                <div className="container-reportes">
                    <div className="box">
                        <input id="input-rep" class="form-control form-control-lg" type="text" placeholder="Ingresar comando Rep" aria-label=".form-control-lg example"
                            onChange={changeText}
                        />
                        <button id="ejecutar-rep" type="button" class="btn btn-light" onClick={enviar_Exec}>Ejecutar</button>
                        <div className="report-render">
                            <Graphviz dot={dot} options={{ width: "100%", zoom: true, fit:true }} />
                        </div>
                    </div>
                </div>

            </div>

        </>
    );

}
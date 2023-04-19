import React from "react";
import { useState } from "react";

import Editor from "../Components/Editor.js";
import Login from "./Src/Object.js";
import "../Styles/Editores.css";

export default function Principal() {

    //* Primeramente mostramos el contenido del archivo en el text area Entrada
    const [archivo, setArchivo] = useState("");
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

    

    return (
        <>
            <div class="Editores">
                <div class="principal">
                    <div class="editores">
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
                                />
                                <button class="custom-btn btn-11">
                                    Ejecutar<div class="dot"></div>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}

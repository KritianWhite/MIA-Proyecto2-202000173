import React from "react";
import "../Styles/Reportes.css"


export default function Reportes() {

    return (
        <>

            <div className="reportes">
                <h1 class="titulo-report">REPORTES</h1>
                <div className="container-reportes">
                    <div className="box">
                        <input id="input-rep" class="form-control form-control-lg" type="text" placeholder="Ingresar comando Rep" aria-label=".form-control-lg example" />
                        <button id="ejecutar-rep" type="button" class="btn btn-light">Ejecutar</button>
                        <div className="report-render">

                        </div>
                    </div>
                </div>

            </div>

        </>
    );

}
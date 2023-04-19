import React, { useState } from "react";
import "../Styles/Login.css";

function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

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

    return (
        <>
            <div className="login-container">
                <form onSubmit={handleSubmit}>
                    <h1>Iniciar sesión</h1>
                    <label htmlFor="id">ID Particion:</label>
                    <input
                        type="text"
                        id="idParticion"
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
                    <button type="submit">Iniciar sesión</button>
                    {/* <button type="submit">Cerrar Sesion</button> */}
                </form>
            </div>
        </>
    );
}

export default Login;

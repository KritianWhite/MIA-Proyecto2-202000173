import React, { useState } from "react";
import "../Styles/Login.css"; // importa tu archivo de estilos CSS aquí

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleEmailChange = (event) => {
    setEmail(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    // Lógica para enviar los datos de inicio de sesión al servidor
  };

  return (
    <>
      <div className="login-container">
        <form onSubmit={handleSubmit}>
          <h1>Iniciar sesión</h1>
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
        </form>
      </div>
    </>
  );
}

export default Login;

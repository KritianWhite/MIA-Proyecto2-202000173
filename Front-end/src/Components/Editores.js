import React, { useRef, useState } from "react";
import axios from "axios";
import "../Styles/Editores.css";

const Consoles2 = () => {
  /**Extraemos el contenido de la entrada de texto */
  const [text, setText] = useState(""); // Estado local para almacenar el valor del textarea
  const [response, setResponse] = useState(""); // Estado local para almacenar el valor del textarea

  const handleTextAreaChange = (event) => {
    // Función de manejo de cambios para capturar el texto del textarea
    const newText = event.target.value; // Obtener el nuevo valor del textarea
    setText(newText); // Actualizar el estado local con el nuevo valor del textarea
  };

  const textAreaRef = useRef(null);
  const handleFileInputChange = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target.result;
      textAreaRef.current.value = content;
    };
    reader.readAsText(file);
  };
  console.log(text);

  /**Creamos la estructura json */
  var Commands = [];
  Commands = text.split("\n");
  Commands = Commands.map((item) => item.trim())
    .map((str) => (str !== "" ? str : null))
    .filter((str) => str !== null);
  //console.log(Commands);
  const entrada = {
    comando: "",
    idParticion: "",
    idU: 0,
    idG: 0,
    idMoun: " ",
    nombreU: " ",
    login: false,
  };
  console.log("Entrada inicial: ", entrada);

  const updateUsuario = (user) => {
    entrada.comando = "";
    entrada.idParticion = user.id_mount;
    entrada.idU = user.id_u;
    entrada.idG = user.id_g;
    entrada.idMoun = user.id_mount;
    entrada.nombreU = user.nombre_u;
    entrada.login = user.login;

    sessionStorage.setItem("entrada", JSON.stringify(entrada));
    console.log(
      "Entrada Updated: ",
      JSON.parse(sessionStorage.getItem("entrada"))
    );
  };

  const enviarJsonAGo_Lista = () => {
    console.log(Commands)
    entrada.comando = Commands;
    console.log("Enviar entrada: ", entrada);
    // Realizar la solicitud POST con Axios
    axios
      .post("http://localhost:8080/Exec", entrada)
      .then((response) => {
        // Manejar la respuesta del servidor si es necesario
        setResponse(response.data.res);
        //console.log("Respuesta usuarios: ", response.data.usuario);
        //const user = response.data.usuario;
        //console.log("Respuesta User: ", user)
        updateUsuario(response.data.usuario);
      })
      .catch((error) => {
        // Manejar errores si los hay
        console.error("Error al enviar JSON a Backend:", error);
      });
  };

  return (
    <>
      <div className="editores">
        <div className="container">
          <div className="code-editor">
            <h2>Editor de código</h2>
            <textarea
              value={text}
              onChange={handleTextAreaChange}
              ref={textAreaRef}
              rows={20}
              cols={80}
              className="code-textarea"
            />
            <input
              onChange={handleFileInputChange}
              class="form-control form-control-lg"
              id="formFileLg"
              type="file"
            />
          </div>

          <div className="console">
            <h2>Consola</h2>
            <textarea
              rows={20}
              cols={80}
              className="console-textarea"
              readOnly // Para que no se pueda editar el textarea
              value={response}
            />
            <div class="d-grid gap-2">
              <button
                className="btn-ejecutar"
                type="button"
                class="btn btn-light"
                onClick={enviarJsonAGo_Lista}
              >
                Ejecutar
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Consoles2;

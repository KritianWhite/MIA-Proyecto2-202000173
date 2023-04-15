document
  .getElementById("input-archivo")
  .addEventListener("change", cargarArchivo);

function cargarArchivo(evento) {
  const archivo = evento.target.files[0];
  const lector = new FileReader();

  lector.onload = function (e) {
    const contenido = e.target.result;
    document.getElementById("code-textarea").value = contenido;
  };

  lector.readAsText(archivo);
}

class Login {
  entrada = {
    comando: "",
    idU: 0,
    idG: 0,
    idMount: " ",
    nombreU: " ",
    login: false,
  };

  constructor() {
    let user = JSON.parse(sessionStorage.getItem("user"));
    if (user == null) {
      this.entrada.idU = 0;
      this.entrada.idG = 0;
      this.entrada.idMount = " ";
      this.entrada.nombreU = " ";
      this.entrada.login = false;
    } else {
      this.entrada.idU = user.idU;
      this.entrada.idG = user.idG;
      this.entrada.idMount = user.idMount;
      this.entrada.nombreU = user.nombreU;
      this.entrada.login = user.login;
    }
    sessionStorage.setItem("entrada", JSON.stringify(this.entrada));
  }

  getUsuario() {
    return JSON.parse(sessionStorage.getItem("user"));
  }

  updateUsuario(user) {
    this.entrada.comando = "";
    this.entrada.idU = user.id_u;
    this.entrada.idG = user.id_g;
    this.entrada.idMount = user.id_mount;
    this.entrada.nombreU = user.nombre_u;
    this.entrada.login = user.login;

    sessionStorage.setItem("user", JSON.stringify(this.entrada));
    console.log("Entrada: ", JSON.parse(sessionStorage.getItem("user")));
  }
}

var logeo = new Login();
let contenido = "";
var textareaArchivo = document.getElementById("code-textarea");
var inputRep = document.getElementById("input-rep")
var textareaArchivo2 = document.getElementById("console-textarea");
console.log(
  "Login: ",
  JSON.parse(sessionStorage.getItem("user")),
  sessionStorage
);

function enviarData() {
  contenido = textareaArchivo.value;
  //textareaArchivo.value = "";
  console.log(contenido);
  const listaComandos = contenido
    .split("\n")
    .map((s) => s.trim())
    .map((str) => (str !== "" ? str : null))
    .filter((str) => str !== null);

  console.log(listaComandos);
  let user = logeo.getUsuario();
  if (!user) {
    user = logeo.entrada;
  }
  userSend = {
    comandos: listaComandos,
    idU: user.id_u,
    idG: user.id_g,
    idMount: user.id_mount,
    nombreU: user.nombre_u,
    login: user.login,
  };

  fetch("http://localhost:8080/Exec", {
    method: "POST",
    body: JSON.stringify(userSend),
    headers: {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*", // Required for CORS support to work
    },
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Error:", err);
      alert("Ocurrió un error, ver la consola");
    })
    .then((response) => {
      console.log("Respuesta del servidor por Exec:", response);
      let respuesta = response.res;
      let usuario = response.usuario;
      textareaArchivo2.value = respuesta;
      logeo.updateUsuario(usuario);
      console.log(respuesta, "respuesta user: ", response);
    });
}

function enviarRep() {
  contenido = inputRep.value;
  //textareaArchivo.value = "";
  console.log(contenido);
  const listaComandos = contenido
    .split("\n")
    .map((s) => s.trim())
    .map((str) => (str !== "" ? str : null))
    .filter((str) => str !== null);

  console.log(listaComandos);
  let user = logeo.getUsuario();
  if (!user) {
    user = logeo.entrada;
  }
  userSend = {
    comandos: listaComandos,
    idU: user.id_u,
    idG: user.id_g,
    idMount: user.id_mount,
    nombreU: user.nombre_u,
    login: user.login,
  };

  fetch("http://localhost:8080/Exec", {
    method: "POST",
    body: JSON.stringify(userSend),
    headers: {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*", // Required for CORS support to work
    },
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Error:", err);
      alert("Ocurrió un error, ver la consola");
    })
    .then((response) => {
      console.log("Respuesta del servidor por Exec:", response);
      let respuesta = response.res;
      let usuario = response.usuario;
      textareaArchivo2.value = respuesta;
      logeo.updateUsuario(usuario);
      console.log(respuesta, "respuesta user: ", response);
    });
}
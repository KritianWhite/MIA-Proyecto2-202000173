export default class Login {
  entrada = {
    comandos: [],
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
    sessionStorage.setItem("user", JSON.stringify(this.entrada));
  }

  getUsuario() {
    return this.entrada;
  }

  updateUsuario(user) {
    this.entrada.comandos = "";
    this.entrada.idU = user.id_u;
    this.entrada.idG = user.id_g;
    this.entrada.idMount = user.id_mount;
    this.entrada.nombreU = user.nombre_u;
    this.entrada.login = user.login;

    sessionStorage.setItem("user", JSON.stringify(this.entrada));
    //console.log("Entrada: ", JSON.parse(sessionStorage.getItem("user")));
  }
}

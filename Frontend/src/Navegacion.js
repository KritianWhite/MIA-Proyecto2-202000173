
document.getElementById("btn-principal").onclick = Editores
document.getElementById("btn-login").onclick = Loginn
document.getElementById("btn-ingresar").onclick = Ingresar
document.getElementById("btn-login-editores").onclick = Editores
document.getElementById("btn-editores-home").onclick = Home
document.getElementById("btn-editor-login").onclick = Reportes
document.getElementById("btn-rep").onclick = Reportes
document.getElementById("btnxx"),onclick = Editores



function Editores(){
  document.getElementById("home").style.display = "none"
  document.getElementById("Editores").style.display = "block"
  document.getElementById("login").style.display = "none"
  document.getElementById("reportes").style.display = "none"
}

function Loginn(){
  document.getElementById("home").style.display = "none"
  document.getElementById("login").style.display = "block"
  document.getElementById("Editores").style.display = "none"
  document.getElementById("reportes").style.display = "none"
}

function Ingresar(){
  document.getElementById("home").style.display = "none"
  document.getElementById("login").style.display = "none"
  document.getElementById("Editores").style.display = "block"
  document.getElementById("reportes").style.display = "none"
}

function Home(){
  document.getElementById("home").style.display = "block"
  document.getElementById("login").style.display = "none"
  document.getElementById("Editores").style.display = "none"
  document.getElementById("reportes").style.display = "none"
}

function Reportes(){
  document.getElementById("home").style.display = "none"
  document.getElementById("login").style.display = "none"
  document.getElementById("Editores").style.display = "none"
  document.getElementById("reportes").style.display = "block"
}
package main

import (
	"MIA-Proyecto2-202000173/Files_System"
	"MIA-Proyecto2-202000173/Structs"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Inicio")
	router := mux.NewRouter()
	enableCORS(router)

	router.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		res := Structs.Inicio{
			Res: "Simulador de Disco Duro Web Corriendo",
			U:   Files_System.UsuarioL,
		}
		json.NewEncoder(writer).Encode(res)
	}).Methods("GET")

	router.HandleFunc("/Entrada", func(writer http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body) // response body is []byte
		var entrada Structs.Entrada
		if err := json.Unmarshal(body, &entrada); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("Error al recibir el comando")
			fmt.Println(err)
		}
		//fmt.Println(entrada.Command)

		reco := recover()
		if reco != nil {
			json.NewEncoder(writer).Encode(Structs.Inicio{Res: "Error en la entrada"})
		}

		Files_System.UsuarioL = Structs.Usuario{
			IdU:     entrada.IdU,
			IdG:     entrada.IdG,
			IdMount: entrada.IdMount,
			NombreU: entrada.NombreU,
			Login:   entrada.Login,
		}
		r := Files_System.Lector(entrada.Command)
		r.U = Files_System.UsuarioL

		json.NewEncoder(writer).Encode(r)
	}).Methods("GET", "POST")

	router.HandleFunc("/Exec", func(writer http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body) // response body is []byte
		var entrada Structs.Exec
		if err := json.Unmarshal(body, &entrada); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("Error al recibir el comando")
			fmt.Println(err)
		}
		reco := recover()
		if reco != nil {
			json.NewEncoder(writer).Encode(Structs.Inicio{Res: "Error en la entrada"})
		}

		fmt.Println("--->",entrada)

		Files_System.UsuarioL = Structs.Usuario{
			IdU:     entrada.IdU,
			IdG:     entrada.IdG,
			IdMount: entrada.IdMount,
			NombreU: entrada.NombreU,
			Login:   entrada.Login,
		}
		respuesta := ""
		dot := ""
		for _, s := range entrada.Commands {
			fmt.Println(s)
			if s == "pause" {
				fmt.Println("Presione Enter para continuar")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				fmt.Println()
				entrada.I++
				continue
			}
			eje := Files_System.Lector(s)
			if eje.Res != "" {
				respuesta += strconv.Itoa(entrada.I) + ")" + eje.Res + "\n"
				dot = eje.Dot
			}
			entrada.I++
		}
		r := Structs.Resp{Res: respuesta, Dot: dot}
		r.U = Files_System.UsuarioL

		json.NewEncoder(writer).Encode(r)
	}).Methods("GET", "POST")

	router.HandleFunc("/ListRep", func(writer http.ResponseWriter, req *http.Request) {
		r := make([]string, 2)
		reportes, er := os.ReadDir("Reportes/")
		if er != nil {
			fmt.Println(er)
		}
		for _, reporte := range reportes {
			i := find(reporte.Name(), ".")
			if i >= len(reporte.Name()) {
				r = append(r, reporte.Name())
			}

		}
		json.NewEncoder(writer).Encode(r)
	}).Methods("GET")

	router.PathPrefix("/Reportes/").Handler(http.StripPrefix("/Reportes/", http.FileServer(http.Dir("./Reportes/"))))

	direccion := ":8080" // Como cadena, no como entero; porque representa una direcciÃ³n
	fmt.Println("Servidor listo escuchando en " + direccion)
	log.Fatal(http.ListenAndServe(direccion, router))
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// Just put some headers to allow CORS...
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,Access-Control-Allow-Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}

func find(cadena string, substring string) int {
	i := strings.Index(cadena, substring)
	if i == -1 {
		i = len(cadena)
	}
	return i
}
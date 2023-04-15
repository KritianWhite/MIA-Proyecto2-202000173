package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

var UsuarioL Structs.Usuario
var Mlist MountList

func Lector(comando string) Structs.Resp {
	//var c Structs.Comando
	entradaI := comando
	entradaL := strings.ToLower(comando)

	if len(entradaI) > 0 {
		fmt.Println(entradaI)
		if strncmp(entradaL, "#") {
			return Structs.Resp{Res: ""}
		}else if strncmp(entradaL, "mkdisk") {
			//c = Structs.Comando{Name: "mkdisk"}
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    Pdisk = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			//propiedadesTemp := make([]Structs.Propiedad, len(propiedades))
			for i := 0; i < len(propiedades); i++ {
				//fmt.Println(propiedades[i])
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">size":
							s, _ := strconv.Atoi(valor_propiedad_comando[1]) 
							Sdisk = s //*Asigna el valor de size
						case ">fit":
							Fdisk = strings.ToLower(valor_propiedad_comando[1]) //*Asigna el valor de fit
						case ">unit":
							Udisk = strings.ToLower(valor_propiedad_comando[1]) //*Asigna el valor de unit
						case ">path":
							Pdisk = valor_propiedad_comando[1] //*Asigna el valor de path
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
					//propiedadesTemp[i] = Structs.Propiedad{Name: valor_propiedad_comando[0], Val: valor_propiedad_comando[1]}
				}
			}
			//c.Propiedades = propiedadesTemp
			//fmt.Println(c)
			return MkDisk()


		}else if strncmp(entradaL, "rmdisk") {
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    Pdisk = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">path":
							Pdisk = valor_propiedad_comando[1] //*Asigna el valor de path
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			//fmt.Println(Pdisk)
			return RmDisk()
		}else if strncmp(entradaL, "fdisk") {
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    Ppart = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">size":
							s, _ := strconv.Atoi(valor_propiedad_comando[1])
							Spart = s //*Asigna el valor de size
						case ">unit":
							str := string(valor_propiedad_comando[1][0])
							str = strings.ToLower(str)
							Upart = str[0] //*Asigna el valor de unit 
						case ">path":
							Ppart = valor_propiedad_comando[1] //*Asigna el valor de path
						case ">type":
							str := string(valor_propiedad_comando[1][0])
							str = strings.ToLower(str)
							Tpart = str[0] //*Asigna el valor de type
						case ">fit":
							str := string(valor_propiedad_comando[1][0])
							str = strings.ToLower(str)
							Fpart = str[0] //*Asigna el valor de fit
						case ">name":
							Namepart = valor_propiedad_comando[1] //*Asigna el valor de name	
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return fdisk()
		}else if strncmp(entradaL, "mount") {
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    Pmontar = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">path":
							Pmontar = valor_propiedad_comando[1] //*Asigna el valor de path
						case ">name":
							Namemontar = valor_propiedad_comando[1] //*Asigna el valor de name
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return mount()
		}else if strncmp(entradaL, "mkfs") {
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">id":
							IdMontar = valor_propiedad_comando[1] //*Asigna el valor de path
						case ">type":
							Tmontar = valor_propiedad_comando[1] //*Asigna el valor de name
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return mkfs()
		}else if strncmp(entradaL, "login"){
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">user":
							original := valor_propiedad_comando[1]
							NameUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de user
						case ">pwd":
							original := valor_propiedad_comando[1]
							PassUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de pwd
						case ">id":
							IdUsuario = valor_propiedad_comando[1] //*Asigna el valor del id
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return login()
		}else if strncmp(entradaL, "logout"){
			return logout()	
		}else if strncmp(entradaL, "mkgrp"){
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">name=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    NameUsuario = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">name":
							original := valor_propiedad_comando[1]
							NameUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de user 
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return mkgrp()
		}else if strncmp(entradaL, "rmgrp"){
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">name=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    NameUsuario = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">name":
							original := valor_propiedad_comando[1]
							NameUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de user 
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return rmgrp()
		}else if strncmp(entradaL, "mkusr"){
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">user":
							original := valor_propiedad_comando[1]
							NameUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de user 
						case ">pwd":
							original := valor_propiedad_comando[1]
							PassUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de pwd
						case ">grp":
							original := valor_propiedad_comando[1]
							GrupoUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de grp
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return mkusr()
		}else if strncmp(entradaL, "rmusr"){
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">user":
							original := valor_propiedad_comando[1]
							NameUsuario = strings.Replace(original, `"`,"",-1) //*Asigna el valor de user 
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return rmusr()
		}else if strncmp(entradaL, "mkfile"){
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    PathArchivos = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">path":
							original := valor_propiedad_comando[1]
							PathArchivos = strings.Replace(original, `"`,"",-1) //*Asigna el valor de path
						case ">size":
							s, _ := strconv.Atoi(valor_propiedad_comando[1]) 
							SArchivos = s //*Asigna el valor de size
						case ">cont":
							original := valor_propiedad_comando[1]
							ContArchivos = strings.Replace(original, `"`,"",-1) //*Asigna el valor de cont
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			if(strings.Contains(entradaI, ">r")){
				RArchivos = true //*Asigna el valor de r
			}
			//fmt.Println(PathArchivos, RArchivos, SArchivos, ContArchivos)
			return mkfile()
		}else if strncmp(entradaL, "mkdir"){
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    PathArchivos = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">path":
							original := valor_propiedad_comando[1]
							PathArchivos = strings.Replace(original, `"`,"",-1) //*Asigna el valor de path
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			if(strings.Contains(entradaI, ">r")){
				RArchivos = true //*Asigna el valor de r
			}
			return mkdir()
		}else if strncmp(entradaL, "rep"){
			if strings.Contains(entradaI, "\"") {
                re := regexp.MustCompile(`"([^"]+)"`)
                result := re.FindStringSubmatch(entradaI)
                if len(result) > 1 {
                    var pathPure = ">path=\"" +result[1]+ "\""
                    entradaI = strings.Replace(entradaI, pathPure, "", -1)
                    Prep = pathPure
                    //fmt.Println(result[1])
                    //fmt.Println(entradaI)
                }
            }
			propiedades := strings.Split(string(entradaI), " ")
			for i := 0; i < len(propiedades); i++ {
				if strings.Contains(propiedades[i], "=") {
					valor_propiedad_comando := strings.Split(propiedades[i], "=")
					//fmt.Println(valor_propiedad_comando)
					switch valor_propiedad_comando[0] {
						case ">path":
							original := valor_propiedad_comando[1]
							Prep = strings.Replace(original, `"`,"",-1) //*Asigna el valor de path
						case ">name":
							original := valor_propiedad_comando[1]
							Namerep = strings.Replace(original, `"`,"",-1) //*Asigna el valor de name
						case ">id":
							original := valor_propiedad_comando[1]
							Idrep = strings.Replace(original, `"`,"",-1) //*Asigna el valor de id
						case ">ruta":
							original := valor_propiedad_comando[1]
							Rutarep = strings.Replace(original, `"`,"",-1) //*Asigna el valor de ruta
						default:
							return Structs.Resp{Res: "ERROR EN EL COMANDO DE ENTRADA: " + entradaI}
					}
				}
			}
			return GenerateRep()
		}
	}
	return Structs.Resp{Res: ""}
}

func strncmp(entrada string, comparacion string) bool {
	if len(comparacion) > len(entrada) {
		return false
	}

	for i := 0; i < len(comparacion); i++ {
		if i >= len(entrada) {
			return false
		}
		if entrada[i] != comparacion[i] {
			return false
		}
	}

	return true
}

func find(cadena string, substring string) int {
	i := strings.Index(cadena, substring)
	if i == -1 {
		i = len(cadena)
	}
	return i
}

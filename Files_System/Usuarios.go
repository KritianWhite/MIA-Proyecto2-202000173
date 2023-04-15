package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var IdUsuario = " "
var NameUsuario = " "
var PassUsuario = " "
var GrupoUsuario = " "
var CambioCont = false
var sbU Structs.SuperBloque
var fileU *os.File
var errfU error

func logout() Structs.Resp {
	if UsuarioL.IdU == 0 {
		return Structs.Resp{Res: "NO SE HA INICIADO UNA SESION CON ANTERIORIDAD"}
	}
	UsuarioL.NombreU = " "
	UsuarioL.IdMount = " "
	UsuarioL.IdU = 0
	UsuarioL.IdG = 0
	UsuarioL.Login = false
	return Structs.Resp{Res: "SE CERRO LA SESION"}
}

func login() Structs.Resp {
	defer func() {
		IdUsuario = " "
		NameUsuario = " "
		PassUsuario = " "
		GrupoUsuario = " "
		sbU = Structs.SuperBloque{}
		fileU = nil
	}()

	if UsuarioL.IdU != 0 {
		return Structs.Resp{Res: "YA HAY UNA SESION ACTIVA"}
	}

	nodo := Mlist.buscar(IdUsuario)
	if nodo != nil {
		fileU, errfU = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
		if errfU == nil {
			banderaU := false
			banderaG := false
			if nodo.Type == 'p' {
				mbr := Structs.MBR{}
				fileU.Seek(0, 0)
				errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
				if mbr.Mbr_partition[nodo.Pos].Part_status != '2' {
					return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
				}
				fileU.Seek(int64(nodo.Start), 0)
				errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
			} else if nodo.Type == 'l' {
				ebr := Structs.EBR{}
				fileU.Seek(int64(nodo.Start), 0)
				errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
				if ebr.Part_status != '2' {
					return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
				}
				fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
				errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
			}
			content := getConten(int(sbU.S_inode_start) + int(unsafe.Sizeof(Structs.TablaInodo{})))
			usuarios := splitUsr(content)
			grupos := splitGrp(content)
			var datosU []string
			var datosG []string
			for i := 0; i < len(usuarios); i++ {
				datosU = strings.Split(usuarios[i], ",")
				if datosU[3] == NameUsuario {
					banderaU = true
					for j := 0; j < len(usuarios); j++ {
						datosG = strings.Split(grupos[j], ",")
						if datosG[2] == datosU[2] {
							banderaG = true
							goto t0
						}
					}
				}
			}
		t0:
			if banderaU {
				if banderaG {
					UsuarioL.NombreU = datosU[3]
					UsuarioL.IdMount = IdUsuario
					IdU, _ := strconv.Atoi(datosU[0])
					UsuarioL.IdU = int32(IdU)
					IdG, _ := strconv.Atoi(datosG[0])
					UsuarioL.IdG = int32(IdG)
					UsuarioL.Login = true
					fileU.Close()
					return Structs.Resp{Res: "SE INICIO SESION COMO " + UsuarioL.NombreU}
				}
				fileU.Close()
				return Structs.Resp{Res: "GRUPO NO ENCONTRADO"}
			}
			fileU.Close()
			return Structs.Resp{Res: "USUARIO NO ENCONTRADO"}
		}
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + IdUsuario}
}

func mkgrp() Structs.Resp {
	defer func() {
		IdUsuario = " "
		NameUsuario = " "
		PassUsuario = " "
		GrupoUsuario = " "
		sbU = Structs.SuperBloque{}
		fileU = nil
	}()
	if UsuarioL.IdU == 1 && UsuarioL.IdG == 1 {
		if NameUsuario == " " {
			return Structs.Resp{Res: "TIENE QUE INGRESAR COMO EL USUARIO root DEL GRUPO 1 PARA EJECUTAR ESTE COMANDO"}
		}

		nodo := Mlist.buscar(UsuarioL.IdMount)
		if nodo != nil {
			fileU, errfU = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
			if errfU == nil {
				if nodo.Type == 'p' {
					fileU.Seek(int64(nodo.Start), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				} else if nodo.Type == 'l' {
					fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				}

				content := getConten(int(sbU.S_inode_start) + int(unsafe.Sizeof(Structs.TablaInodo{})))
				nuevoGrp := ""
				cantBlockAnt := len(splitConten(content))
				grupos := splitGrp(content)

				if !grupoExist(grupos, NameUsuario) {
					nuevoGrp = getGID(grupos) + ",G," + NameUsuario + "\n"
					content += nuevoGrp
					usersTxt := splitConten(content)
					cantBlockAct := len(usersTxt)

					if len(usersTxt) > 16 {
						return Structs.Resp{Res: "NO SE PUEDE GUARDAR EL GRUPO"}
					}

					if int(sbU.S_free_blocks_count) < (cantBlockAct - cantBlockAnt) {
						return Structs.Resp{Res: "NO HAY SUFICIENTES BLOQUES LIBRES PARA ACTUALIZAR EL ARCHIVO"}
					}

					//Se busca espacion en el bitmap de bloques
					var bit byte
					start := int(sbU.S_bm_block_start)
					end := start + int(sbU.S_block_start)
					cantContiguos := 0
					inicioBM := -1
					inicioB := -1
					contadorA := 0
					if (cantBlockAct - cantBlockAnt) > 0 {
						for i := start; i < end; i++ {
							fileU.Seek(int64(i), 0)
							errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
							if bit == '1' {
								cantContiguos = 0
								inicioBM = -1
								inicioB = -1
							} else {
								if cantContiguos == 0 {
									inicioBM = i
									inicioB = contadorA
								}
								cantContiguos++
							}
							if cantContiguos >= (cantBlockAct - cantBlockAnt) {
								break
							}
							contadorA++
						}
						if inicioBM == -1 || cantContiguos != (cantBlockAct-cantBlockAnt) {
							return Structs.Resp{Res: "NO HAY SUFICIENTES BLOQUES CONTIGUOS PARA ACTUALIZAR EL ARCHIVO"}
						}

						for i := inicioBM; i < (inicioBM + (cantBlockAct - cantBlockAnt)); i++ {
							var uno byte = '1'
							fileU.Seek(int64(i), 0)
							var bufferByte bytes.Buffer
							errf := binary.Write(&bufferByte, binary.BigEndian, uno)
							if errf != nil {
								fmt.Println(errf)
							}
							EscribirFile(fileU, bufferByte.Bytes())
						}
						sbU.S_free_blocks_count = int32(int(sbU.S_free_blocks_count) - (cantBlockAct - cantBlockAnt))
						bit2 := 0
						for k := start; k < end; k++ {
							fileU.Seek(int64(k), 0)
							errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
							if bit == '0' {
								break
							}
							bit2++
						}
						sbU.S_first_blo = int32(bit2)
					}
					inodo := Structs.TablaInodo{}
					seekInodo := int(sbU.S_inode_start) + int(unsafe.Sizeof(inodo))
					fileU.Seek(int64(seekInodo), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

					tamanio := 0
					for tm := 0; tm < len(usersTxt); tm++ {
						tamanio += len(usersTxt[tm])
					}
					inodo.I_s = int32(tamanio)
					inodo.I_mtime = time.Now().Unix()

					contador := 0
					j := 0
					for j < len(usersTxt) {
						CambioCont = false
						inodo = agregarArchivo(usersTxt[j], inodo, j, inicioB+contador)
						if CambioCont {
							contador++
						}
						j++
					}

					fileU.Seek(int64(seekInodo), 0)
					var bufferInodo bytes.Buffer
					errf := binary.Write(&bufferInodo, binary.BigEndian, inodo)
					if errf != nil {
						fmt.Println(errf)
					}
					EscribirFile(fileU, bufferInodo.Bytes())
					if nodo.Type == 'p' {
						fileU.Seek(int64(nodo.Start), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					} else if nodo.Type == 'l' {
						fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					}

					fileU.Close()
					return Structs.Resp{Res: "SE GUARDO EL GRUPO " + NameUsuario}

				}
				return Structs.Resp{Res: "GRUPO " + NameUsuario + " YA EXISTE"}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}
	return Structs.Resp{Res: "TIENE QUE INGRESAR COMO EL USUARIO root DEL GRUPO 1 PARA EJECUTAR ESTE COMANDO"}
}

func rmgrp() Structs.Resp {
	defer func() {
		IdUsuario = " "
		NameUsuario = " "
		PassUsuario = " "
		GrupoUsuario = " "
		sbU = Structs.SuperBloque{}
		fileU = nil
	}()
	if UsuarioL.IdU == 1 && UsuarioL.IdG == 1 {
		if NameUsuario == " " {
			return Structs.Resp{Res: "DEBE INGRESAR NOMBRE DEL GRUPO"}
		}

		nodo := Mlist.buscar(UsuarioL.IdMount)
		if nodo != nil {
			fileU, errfU = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
			if errfU == nil {
				if nodo.Type == 'p' {
					fileU.Seek(int64(nodo.Start), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				} else if nodo.Type == 'l' {
					fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				}

				content := getConten(int(sbU.S_inode_start) + int(unsafe.Sizeof(Structs.TablaInodo{})))
				grupos := splitGrp(content)
				if grupoExist(grupos, NameUsuario) {
					for i := 0; i < len(grupos); i++ {
						datos := strings.Split(grupos[i], ",")
						if datos[2] == NameUsuario {
							datos[0] = "0"
							newS := datos[0] + "," + datos[1] + "," + datos[2]
							content = strings.Replace(content, grupos[i], newS, 1)
						}
					}
					usersTxt := splitConten(content)
					inodo := Structs.TablaInodo{}
					seekInodo := int(sbU.S_inode_start) + int(unsafe.Sizeof(inodo))
					fileU.Seek(int64(seekInodo), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

					tamanio := 0
					for tm := 0; tm < len(usersTxt); tm++ {
						tamanio += len(usersTxt[tm])
					}
					inodo.I_s = int32(tamanio)
					inodo.I_mtime = time.Now().Unix()
					inodo.I_atime = time.Now().Unix()

					j := 0
					for j < len(usersTxt) {
						inodo = agregarArchivo(usersTxt[j], inodo, j, -1)
						j++
					}

					fileU.Seek(int64(seekInodo), 0)
					var bufferInodo bytes.Buffer
					errf := binary.Write(&bufferInodo, binary.BigEndian, inodo)
					if errf != nil {
						fmt.Println(errf)
					}
					EscribirFile(fileU, bufferInodo.Bytes())
					if nodo.Type == 'p' {
						fileU.Seek(int64(nodo.Start), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					} else if nodo.Type == 'l' {
						fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					}

					fileU.Close()

					return Structs.Resp{Res: "SE REMOVIO EL GRUPO " + NameUsuario}
				}
				return Structs.Resp{Res: "NO EXISTE GRUPO " + NameUsuario}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}
	return Structs.Resp{Res: "TIENE QUE INGRESAR COMO EL USUARIO root DEL GRUPO 1 PARA EJECUTAR ESTE COMANDO"}
}

func rmusr() Structs.Resp {
	defer func() {
		IdUsuario = " "
		NameUsuario = " "
		PassUsuario = " "
		GrupoUsuario = " "
		sbU = Structs.SuperBloque{}
		fileU = nil
	}()
	if UsuarioL.IdU == 1 && UsuarioL.IdG == 1 {
		if NameUsuario == " " {
			return Structs.Resp{Res: "DEBE INGRESAR NOMBRE DEL GRUPO"}
		}
		nodo := Mlist.buscar(UsuarioL.IdMount)
		if nodo != nil {
			fileU, errfU = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
			if errfU == nil {
				if nodo.Type == 'p' {
					fileU.Seek(int64(nodo.Start), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				} else if nodo.Type == 'l' {
					fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				}

				content := getConten(int(sbU.S_inode_start) + int(unsafe.Sizeof(Structs.TablaInodo{})))
				usuarios := splitUsr(content)
				if usrExist(usuarios, NameUsuario) {
					for i := 0; i < len(usuarios); i++ {
						datos := strings.Split(usuarios[i], ",")
						if datos[3] == NameUsuario {
							datos[0] = "0"
							newS := datos[0] + "," + datos[1] + "," + datos[2] + "," + datos[3] + "," + datos[4]
							content = strings.Replace(content, usuarios[i], newS, 1)
						}
					}

					usersTxt := splitConten(content)
					inodo := Structs.TablaInodo{}
					seekInodo := int(sbU.S_inode_start) + int(unsafe.Sizeof(inodo))
					fileU.Seek(int64(seekInodo), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

					tamanio := 0
					for tm := 0; tm < len(usersTxt); tm++ {
						tamanio += len(usersTxt[tm])
					}
					inodo.I_s = int32(tamanio)
					inodo.I_mtime = time.Now().Unix()
					inodo.I_atime = time.Now().Unix()

					j := 0
					for j < len(usersTxt) {
						inodo = agregarArchivo(usersTxt[j], inodo, j, -1)
						j++
					}

					fileU.Seek(int64(seekInodo), 0)
					var bufferInodo bytes.Buffer
					errf := binary.Write(&bufferInodo, binary.BigEndian, inodo)
					if errf != nil {
						fmt.Println(errf)
					}
					EscribirFile(fileU, bufferInodo.Bytes())
					if nodo.Type == 'p' {
						fileU.Seek(int64(nodo.Start), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					} else if nodo.Type == 'l' {
						fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					}

					fileU.Close()
					return Structs.Resp{Res: "SE REMOVIO EL USUARIO " + NameUsuario}
				}
				return Structs.Resp{Res: "NO EXISTE EL USUARIO " + NameUsuario}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}
	return Structs.Resp{Res: "TIENE QUE INGRESAR COMO EL USUARIO root DEL GRUPO 1 PARA EJECUTAR ESTE COMANDO"}
}

func mkusr() Structs.Resp {
	defer func() {
		IdUsuario = " "
		NameUsuario = " "
		PassUsuario = " "
		GrupoUsuario = " "
		sbU = Structs.SuperBloque{}
		fileU = nil
	}()

	if UsuarioL.IdU == 1 && UsuarioL.IdG == 1 {
		if NameUsuario == " " || PassUsuario == " " || GrupoUsuario == " " {
			return Structs.Resp{Res: "ASEGURESE DE INGRESAR TODOS LOS DATOS DEL USUARIO"}
		}

		nodo := Mlist.buscar(UsuarioL.IdMount)
		if nodo != nil {
			fileU, errfU = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
			if errfU == nil {
				if nodo.Type == 'p' {
					fileU.Seek(int64(nodo.Start), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				} else if nodo.Type == 'l' {
					fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(sbU))), binary.BigEndian, &sbU)
				}

				content := getConten(int(sbU.S_inode_start) + int(unsafe.Sizeof(Structs.TablaInodo{})))
				nuevoUsr := ""
				cantBlockAnt := len(splitConten(content))
				grupos := splitGrp(content)
				usuarios := splitUsr(content)

				if usrExist(usuarios, NameUsuario) {
					return Structs.Resp{Res: "NOMBRE DE USUARIO YA REGISTRADO EN EL SISTEMA"}
				}
				if grupoExist(grupos, GrupoUsuario) {
					nuevoUsr = getUID(usuarios) + ",U," + GrupoUsuario + "," + NameUsuario + "," + PassUsuario + "\n"
					content += nuevoUsr
					usersTxt := splitConten(content)
					cantBlockAct := len(usersTxt)

					if len(usersTxt) > 16 {
						return Structs.Resp{Res: "NO SE PUEDE GUARDAR EL USUARIO"}
					}

					if int(sbU.S_free_blocks_count) < (cantBlockAct - cantBlockAnt) {
						return Structs.Resp{Res: "NO HAY SUFICIENTES BLOQUES LIBRES PARA ACTUALIZAR EL ARCHIVO"}
					}

					//Se busca espacion en el bitmap de bloques
					var bit byte
					start := int(sbU.S_bm_block_start)
					end := start + int(sbU.S_block_start)
					cantContiguos := 0
					inicioBM := -1
					inicioB := -1
					contadorA := 0
					if (cantBlockAct - cantBlockAnt) > 0 {
						for i := start; i < end; i++ {
							fileU.Seek(int64(i), 0)
							errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
							if bit == '1' {
								cantContiguos = 0
								inicioBM = -1
								inicioB = -1
							} else {
								if cantContiguos == 0 {
									inicioBM = i
									inicioB = contadorA
								}
								cantContiguos++
							}
							if cantContiguos >= (cantBlockAct - cantBlockAnt) {
								break
							}
							contadorA++
						}
						if inicioBM == -1 || cantContiguos != (cantBlockAct-cantBlockAnt) {
							return Structs.Resp{Res: "NO HAY SUFICIENTES BLOQUES CONTIGUOS PARA ACTUALIZAR EL ARCHIVO"}
						}

						for i := inicioBM; i < (inicioBM + (cantBlockAct - cantBlockAnt)); i++ {
							var uno byte = '1'
							fileU.Seek(int64(i), 0)
							var bufferByte bytes.Buffer
							errf := binary.Write(&bufferByte, binary.BigEndian, uno)
							if errf != nil {
								fmt.Println(errf)
							}
							EscribirFile(fileU, bufferByte.Bytes())
						}
						sbU.S_free_blocks_count = int32(int(sbU.S_free_blocks_count) - (cantBlockAct - cantBlockAnt))
						bit2 := 0
						for k := start; k < end; k++ {
							fileU.Seek(int64(k), 0)
							errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
							if bit == '0' {
								break
							}
							bit2++
						}
						sbU.S_first_blo = int32(bit2)
					}

					inodo := Structs.TablaInodo{}
					seekInodo := int(sbU.S_inode_start) + int(unsafe.Sizeof(inodo))
					fileU.Seek(int64(seekInodo), 0)
					errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

					tamanio := 0
					for tm := 0; tm < len(usersTxt); tm++ {
						tamanio += len(usersTxt[tm])
					}
					inodo.I_s = int32(tamanio)
					inodo.I_mtime = time.Now().Unix()

					contador := 0
					j := 0
					for j < len(usersTxt) {
						CambioCont = false
						inodo = agregarArchivo(usersTxt[j], inodo, j, inicioB+contador)
						if CambioCont {
							contador++
						}
						j++
					}

					fileU.Seek(int64(seekInodo), 0)
					var bufferInodo bytes.Buffer
					errf := binary.Write(&bufferInodo, binary.BigEndian, inodo)
					if errf != nil {
						fmt.Println(errf)
					}
					EscribirFile(fileU, bufferInodo.Bytes())
					if nodo.Type == 'p' {
						fileU.Seek(int64(nodo.Start), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					} else if nodo.Type == 'l' {
						fileU.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sbU)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileU, bufferSB.Bytes())
					}

					fileU.Close()
					return Structs.Resp{Res: "SE REGISTRO AL USUARIO " + NameUsuario}
				}
				return Structs.Resp{Res: "GRUPO " + GrupoUsuario + " NO EXISTE"}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}

	return Structs.Resp{Res: "TIENE QUE INGRESAR COMO EL USUARIO root DEL GRUPO 1 PARA EJECUTAR ESTE COMANDO"}
}

func getConten(inodoStart int) string {
	var inodo Structs.TablaInodo
	var archivo Structs.BloqueArchivo
	fileU.Seek(int64(inodoStart), 0)
	errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)
	content := ""
	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			fileU.Seek(int64(inodo.I_block[i]), 0)
			errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(archivo))), binary.BigEndian, &archivo)
			content += archivoContent2(archivo.B_content)
		}
	}
	return content
}

func splitUsr(cadena string) []string {
	var split []string
	content := strings.Split(cadena, "\n")
	for i := 0; i < len(content); i++ {
		if content[i] != "" {
			datos := strings.Split(content[i], ",")
			if datos[1] == "U" && datos[0] != "0" {
				split = append(split, content[i])
			}
		}

	}
	return split
}

func splitGrp(cadena string) []string {
	var split []string
	content := strings.Split(cadena, "\n")
	for i := 0; i < len(content); i++ {
		if content[i] != "" {
			datos := strings.Split(content[i], ",")
			if datos[1] == "G" && datos[0] != "0" {
				split = append(split, content[i])
			}
		}
	}
	return split
}

func splitConten(cadena string) []string {
	controlador := 0
	var split []string
	aux := ""
	for i := 0; i < len(cadena); i++ {
		if controlador < 64 {
			aux += string([]byte{cadena[i]})
			controlador++
		}
		if len(aux) == 64 {
			split = append(split, aux)
			aux = ""
			controlador = 0
		}
	}
	if controlador != 0 {
		split = append(split, aux)
		aux = ""
		controlador = 0
	}
	return split
}

func grupoExist(grupos []string, name string) bool {
	var datosG []string
	for i := 0; i < len(grupos); i++ {
		datosG = strings.Split(grupos[i], ",")
		if datosG[2] == name && datosG[0] != "0" {
			return true
		}
	}
	return false
}

func usrExist(usuarios []string, name string) bool {
	var datosG []string
	for i := 0; i < len(usuarios); i++ {
		datosG = strings.Split(usuarios[i], ",")
		if datosG[3] == name && datosG[0] != "0" {
			return true
		}
	}
	return false
}

func getGID(grupos []string) string {
	var datosG []string
	gid := 0
	for i := 0; i < len(grupos); i++ {
		datosG = strings.Split(grupos[i], ",")
		id, _ := strconv.Atoi(datosG[0])
		if gid < id {
			gid++
		}
	}
	return strconv.Itoa(gid + 1)
}

func getUID(usuarios []string) string {
	var datosG []string
	gid := 0
	for i := 0; i < len(usuarios); i++ {
		datosG = strings.Split(usuarios[i], ",")
		id, _ := strconv.Atoi(datosG[0])
		if gid < id {
			gid++
		}
	}
	return strconv.Itoa(gid + 1)
}

func agregarArchivo(cadena string, inodo Structs.TablaInodo, j int, aux int) Structs.TablaInodo {
	in := Structs.TablaInodo{}
	in.I_type = 'F'
	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 && i == j {
			var archivo Structs.BloqueArchivo
			fileU.Seek(int64(inodo.I_block[i]), 0)
			errfU = binary.Read(LeerFile(fileU, int(unsafe.Sizeof(archivo))), binary.BigEndian, &archivo)
			archivo.B_content = archivoContent(cadena)
			fileU.Seek(int64(inodo.I_block[i]), 0)
			var bufferArchivo bytes.Buffer
			errf := binary.Write(&bufferArchivo, binary.BigEndian, archivo)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileU, bufferArchivo.Bytes())
			return inodo
		} else if inodo.I_block[i] == -1 && aux != -1 {
			var archivo Structs.BloqueArchivo
			seek := int(sbU.S_block_start) + (aux * int(unsafe.Sizeof(archivo)))
			archivo.B_content = archivoContent(cadena)
			fileU.Seek(int64(seek), 0)
			var bufferArchivo bytes.Buffer
			errf := binary.Write(&bufferArchivo, binary.BigEndian, archivo)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileU, bufferArchivo.Bytes())
			inodo.I_block[i] = int32(seek)
			CambioCont = true
			return inodo
		}
	}
	return in
}

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

var PathArchivos = " "
var ContArchivos = " "
var DestinoArchivos = " "
var UGO = 0
var NameArchivos = " "
var RArchivos = false
var SArchivos = 0
var cambioContArchivos = false
var sbA Structs.SuperBloque
var fileA *os.File
var errfA error


func mkfile() Structs.Resp {
	defer func() {
		PathArchivos = " "
		ContArchivos = " "
		DestinoArchivos = " "
		UGO = 0
		NameArchivos = " "
		RArchivos = false
		SArchivos = 0
		cambioContArchivos = false
		sbA = Structs.SuperBloque{}
		fileA = nil
	}()

	if PathArchivos == " " {
		return Structs.Resp{Res: "DEBE INGRESAR LA RUTA DEL ARCHIVO A CREAR"}
	}

	nodo := Mlist.buscar(UsuarioL.IdMount)

	if nodo == nil {
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}

	if ContArchivos == " " {
		if SArchivos < 0 {
			return Structs.Resp{Res: "EL VALOR DE size DEBER SER MAYOR O IGUAL A 0"}
		}
	}
	
	if(SArchivos == 0){
		SArchivos = 1 //*Para que no se cree el archivo con el tamaño de 0
	}

	fileA, errfA = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
	if errfA == nil {
		if nodo.Type == 'p' {
			mbr := Structs.MBR{}
			fileA.Seek(0, 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
			if mbr.Mbr_partition[nodo.Pos].Part_status != '2' {
				return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
			}
			fileA.Seek(int64(nodo.Start), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(sbA))), binary.BigEndian, &sbA)
		} else if nodo.Type == 'l' {
			ebr := Structs.EBR{}
			fileA.Seek(int64(nodo.Start), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
			if ebr.Part_status != '2' {
				return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
			}
			fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(sbA))), binary.BigEndian, &sbA)
		}

		rutaS := splitRuta(PathArchivos)
		if rutaS == nil {
			fileA.Close()
			return Structs.Resp{Res: "RUTA INVALIDA"}
		}

		exist := getInodoF(rutaS, 0, len(rutaS)-1, int(sbA.S_inode_start), fileA)
		if exist != -1 {
			fileA.Close()
			return Structs.Resp{Res: "YA EXISTE EL DIRECTORIO " + PathArchivos}
		}

		posInodoI := int(sbA.S_inode_start)
		existP := true
		var inodo Structs.TablaInodo

		if len(rutaS) > 1 {
			for i := 0; i < (len(rutaS) - 1); i++ {
				if existP {
					aux := posInodoI
					posInodoI = existPath(rutaS, i, posInodoI)
					if posInodoI == aux {
						existP = false
					}
				}
				if !existP {
					fmt.Println(RArchivos)
					if RArchivos {
						posInodoI = crearCarpeta(rutaS, i, posInodoI)
						if nodo.Type == 'p' {
							fileA.Seek(int64(nodo.Start), 0)
						} else if nodo.Type == 'l' {
							fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						}
						var bufferSB bytes.Buffer
						errfA = binary.Write(&bufferSB, binary.BigEndian, sbA)
						EscribirFile(fileA, bufferSB.Bytes())
						if posInodoI == -1 {
							fmt.Println("no se puede crear el archivo 1")
							return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
						}
					} else {
						fmt.Println("no se puede crear el archivo 2")
						return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
					}

				}
			}
		}

		if posInodoI == -1 {
			return Structs.Resp{Res: "Algo salio mal"}
		}

		if !validarPermisoW(posInodoI) {
			return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO " + rutaS[len(rutaS)-1] + " POR FALTA DE PERMISOS"}
		}

		texto := ""
		if ContArchivos != " " {
			fileC, err := os.ReadFile(ContArchivos)
			if err == nil {
				texto = string(fileC)
			} else {
				return Structs.Resp{Res: "NO SE ENCONTRO EL ARCHIVO"}
			}
		} else {
			texto = getSize()
		}
		contenido := splitConten(texto)

		if len(contenido) > 16 {
			return Structs.Resp{Res: "ARCHIVO EXCECE LA CANTIDAD DE BLOQUES DISPONIBLES POR INODO"}
		}

		if int(sbA.S_free_blocks_count) < len(contenido) {
			return Structs.Resp{Res: "NO HAY BLOQUES SUFICIENTES PARA CREAR ARCHIVO"}
			fileA.Close()
		}

		fileA.Seek(int64(posInodoI), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

		var bit byte
		start := int(sbA.S_bm_block_start)
		end := start + int(sbA.S_block_start)
		cantContiguos := 0
		inicioBM := -1
		inicioB := -1
		contadorA := 0
		for i := start; i < end; i++ {
			fileA.Seek(int64(i), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
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

			if cantContiguos >= len(contenido) {
				break
			}
			contadorA++
		}

		if (inicioBM == -1 || cantContiguos != len(contenido)) && SArchivos != 0 {
			fileA.Close()
			return Structs.Resp{Res: "NO HAY SUFICIENTES BLOQUES CONTIGUOS PARA ACTUALIZAR EL ARCHIVO"}
		}

		for i := inicioBM; i < (inicioBM + len(contenido)); i++ {
			var uno byte = '1'
			fileA.Seek(int64(i), 0)
			var bufferByte bytes.Buffer
			errf := binary.Write(&bufferByte, binary.BigEndian, uno)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileA, bufferByte.Bytes())
		}

		sbA.S_free_blocks_count = int32(int(sbA.S_free_blocks_count) - len(contenido))

		bit2 := 0
		for k := start; k < end; k++ {
			fileA.Seek(int64(k), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
			if bit == '0' {
				break
			}
			bit2++
		}
		sbA.S_first_blo = int32(bit2)
		if nodo.Type == 'p' {
			fileA.Seek(int64(nodo.Start), 0)
		} else if nodo.Type == 'l' {
			fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
		}
		var bufferSB bytes.Buffer
		errfA = binary.Write(&bufferSB, binary.BigEndian, sbA)
		EscribirFile(fileA, bufferSB.Bytes())

		var newInodoA Structs.TablaInodo
		posNewI := buscarPosicionNewInodo()
		crearInodoArchivo(posNewI)
		fileA.Seek(int64(posNewI), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(newInodoA))), binary.BigEndian, &newInodoA)
		atc := agregarCarpeta(posNewI, posInodoI, rutaS[len(rutaS)-1])

		if atc == -1 {
			return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
		}

		tamanio := 0
		for tm := 0; tm < len(contenido); tm++ {
			tamanio += len(contenido[tm])
		}
		newInodoA.I_s = int32(tamanio)
		newInodoA.I_mtime = time.Now().Unix()
		newInodoA.I_atime = time.Now().Unix()

		contador := 0
		j := 0
		for j < len(contenido) {
			CambioCont = false
			newInodoA = agregarArchivo2(contenido[j], newInodoA, j, inicioB+contador)
			if CambioCont {
				contador++
			}
			j++
		}

		fileA.Seek(int64(posNewI), 0)
		var bufferInodo bytes.Buffer
		errf := binary.Write(&bufferInodo, binary.BigEndian, newInodoA)
		if errf != nil {
			fmt.Println(errf)
		}
		EscribirFile(fileA, bufferInodo.Bytes())
		if nodo.Type == 'p' {
			fileA.Seek(int64(nodo.Start), 0)
		} else if nodo.Type == 'l' {
			fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
		}
		var bufferSB2 bytes.Buffer
		errf = binary.Write(&bufferSB2, binary.BigEndian, sbA)
		if errf != nil {
			fmt.Println(errf)
		}
		EscribirFile(fileA, bufferSB2.Bytes())
		fileA.Close()
		return Structs.Resp{Res: "SE CREO EL ARCHIVO " + PathArchivos}
	}
	return Structs.Resp{Res: "DISCO INEXISTENTE"}
}

func mkdir() Structs.Resp {
	defer func() {
		PathArchivos = " "
		ContArchivos = " "
		DestinoArchivos = " "
		UGO = 0
		NameArchivos = " "
		RArchivos = false
		SArchivos = 0
		cambioContArchivos = false
		sbA = Structs.SuperBloque{}
		fileA = nil
	}()

	if PathArchivos == " " {
		return Structs.Resp{Res: "DEBE INGRESAR LA RUTA DEL ARCHIVO A CREAR"}
	}

	nodo := Mlist.buscar(UsuarioL.IdMount)

	if nodo == nil {
		return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + UsuarioL.IdMount}
	}

	fileA, errfA = os.OpenFile(nodo.Path, os.O_RDWR, 0777)
	if errfA == nil {
		if nodo.Type == 'p' {
			mbr := Structs.MBR{}
			fileA.Seek(0, 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
			if mbr.Mbr_partition[nodo.Pos].Part_status != '2' {
				return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
			}
			fileA.Seek(int64(nodo.Start), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(sbA))), binary.BigEndian, &sbA)
		} else if nodo.Type == 'l' {
			ebr := Structs.EBR{}
			fileA.Seek(int64(nodo.Start), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
			if ebr.Part_status != '2' {
				return Structs.Resp{Res: "NO SE HA FORMATEADO LA PARTICION"}
			}
			fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(sbA))), binary.BigEndian, &sbA)
		}

		rutaS := splitRuta(PathArchivos)
		if rutaS == nil {
			fileA.Close()
			return Structs.Resp{Res: "RUTA INVALIDA"}
		}

		exist := getInodoF(rutaS, 0, len(rutaS)-1, int(sbA.S_inode_start), fileA)

		if exist != -1 {
			fileA.Close()
			return Structs.Resp{Res: "YA EXISTE EL DIRECTORIO " + PathArchivos}
		}

		posInodoI := int(sbA.S_inode_start)
		existP := true
		if len(rutaS) > 1 {
			for i := 0; i < (len(rutaS) - 1); i++ {
				if existP {
					aux := posInodoI
					posInodoI = existPath(rutaS, i, posInodoI)
					if posInodoI == aux {
						existP = false
					}
				}
				if !existP {
					if RArchivos {
						posInodoI = crearCarpeta(rutaS, i, posInodoI)
						if nodo.Type == 'p' {
							fileA.Seek(int64(nodo.Start), 0)
						} else if nodo.Type == 'l' {
							fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
						}
						var bufferSB bytes.Buffer
						errfA = binary.Write(&bufferSB, binary.BigEndian, sbA)
						EscribirFile(fileA, bufferSB.Bytes())
						if posInodoI == -1 {
							return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
						}
					} else {
						return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
					}

				}
			}
		}

		if posInodoI == -1 {
			return Structs.Resp{Res: "Algo salio mal"}
		}

		posInodoI = crearCarpeta(rutaS, len(rutaS)-1, posInodoI)
		if posInodoI == -1 {
			return Structs.Resp{Res: "NO SE PUEDE CREAR EL ARCHIVO"}
		}
		if nodo.Type == 'p' {
			fileA.Seek(int64(nodo.Start), 0)
		} else if nodo.Type == 'l' {
			fileA.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
		}
		var bufferSB bytes.Buffer
		errfA = binary.Write(&bufferSB, binary.BigEndian, sbA)
		EscribirFile(fileA, bufferSB.Bytes())
		fileA.Close()
		return Structs.Resp{Res: "SE CREO EL DIRECTORIO " + PathArchivos}
	}
	return Structs.Resp{Res: "DISCO INEXISTENTE"}
}

func splitRuta(ruta string) []string {
	var splitO []string
	var splitF []string

	splitO = strings.Split(ruta, "/")

	for i := 0; i < len(splitO); i++ {
		if splitO[i] != "" {
			splitF = append(splitF, splitO[i])
		}
	}

	return splitF
}

func existPath(rutaS []string, posAct int, start int) int {
	var inodo Structs.TablaInodo
	var carpeta Structs.BloqueCarpeta

	fileA.Seek(int64(start), 0)
	errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

	if inodo.I_type == '1' {
		return start
	}

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			fileA.Seek(int64(inodo.I_block[i]), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(carpeta))), binary.BigEndian, &carpeta)
			for c := 0; c < 4; c++ {
				name := getContentName(carpeta.B_content[c].B_name)
				if name == rutaS[posAct] {
					return int(carpeta.B_content[c].B_inodo)
				}
			}
		}
	}

	return start
}

func crearCarpeta(rutaS []string, posAct int, posI int) int {
	if !validarPermisoW(posI) {
		return -1
	}

	var inodo Structs.TablaInodo
	var carpeta Structs.BloqueCarpeta

	fileA.Seek(int64(posI), 0)
	errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

	if inodo.I_type == '1' {
		return -1
	}

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			fileA.Seek(int64(inodo.I_block[i]), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(carpeta))), binary.BigEndian, &carpeta)
			for c := 0; c < 4; c++ {
				if carpeta.B_content[c].B_inodo == -1 {
					if sbA.S_free_inodes_count > 0 && sbA.S_free_blocks_count > 0 {
						posInodo := buscarPosicionNewInodo()
						posCarpetaI := crearBloqueCarpetaInicial(posInodo, posI)
						crearInodoCarpeta(posInodo, posCarpetaI)
						carpeta.B_content[c].B_inodo = int32(posInodo)
						carpeta.B_content[c].B_name = nameConten(rutaS[posAct])
						fileA.Seek(int64(inodo.I_block[i]), 0)
						var bufferByte bytes.Buffer
						errf := binary.Write(&bufferByte, binary.BigEndian, carpeta)
						if errf != nil {
							fmt.Println(errf)
						}
						EscribirFile(fileA, bufferByte.Bytes())
						return posInodo
					} else {
						return -1
					}
				}
			}
		} else if inodo.I_block[i] == -1 {
			if sbA.S_free_inodes_count > 1 && sbA.S_free_blocks_count > 0 {
				posInodo := buscarPosicionNewInodo()
				posCarpetaI := crearBloqueCarpetaInicial(posInodo, posI)
				crearInodoCarpeta(posInodo, posCarpetaI)
				posCarpetaO := crearBloqueCarpetaOtra(posInodo, rutaS[posAct])
				inodo.I_block[i] = int32(posCarpetaO)
				fileA.Seek(int64(posI), 0)
				var bufferByte bytes.Buffer
				errf := binary.Write(&bufferByte, binary.BigEndian, inodo)
				if errf != nil {
					fmt.Println(errf)
				}
				EscribirFile(fileA, bufferByte.Bytes())
				return posInodo
			} else {
				return -1
			}
		}
	}
	return -1
}

func buscarPosicionNewInodo() int {
	bitI := 0
	var bit byte
	var one byte = '1'
	startI := int(sbA.S_bm_inode_start)
	endI := startI + int(sbA.S_inodes_count)

	for j := startI; j < endI; j++ {
		fileA.Seek(int64(j), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
		if bit == '0' {
			fileA.Seek(int64(j), 0)
			var bufferByte bytes.Buffer
			errf := binary.Write(&bufferByte, binary.BigEndian, one)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileA, bufferByte.Bytes())
			break
		}
		bitI++
	}
	sbA.S_free_inodes_count -= 1
	posInodo := int(sbA.S_inode_start) + (int(unsafe.Sizeof(Structs.TablaInodo{})) * bitI)
	buscarPrimerInodoVacio()
	return posInodo
}

func buscarPosicionNewBLoque() int {
	bitI := 0
	var bit byte
	var one byte = '1'
	startI := int(sbA.S_bm_block_start)
	endI := startI + int(sbA.S_blocks_count)
	posBloque := -1

	for j := startI; j < endI; j++ {
		fileA.Seek(int64(j), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
		if bit == '0' {
			fileA.Seek(int64(j), 0)
			var bufferByte bytes.Buffer
			errf := binary.Write(&bufferByte, binary.BigEndian, one)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileA, bufferByte.Bytes())
			break
		}
		bitI++
	}
	sbA.S_free_blocks_count -= 1
	posBloque = int(sbA.S_block_start) + (int(unsafe.Sizeof(Structs.BloqueArchivo{})) * bitI)
	buscarPrimerBLoqueVacio()
	return posBloque
}

func buscarPrimerBLoqueVacio() {
	bitI := 0
	var bit byte
	startI := int(sbA.S_bm_block_start)
	endI := startI + int(sbA.S_blocks_count)

	for j := startI; j < endI; j++ {
		fileA.Seek(int64(j), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
		if bit == '0' {
			bitI++
			break
		}
		bitI++
	}
	sbA.S_first_blo = int32(bitI)
}

func buscarPrimerInodoVacio() {
	bitI := 0
	var bit byte
	startI := int(sbA.S_bm_inode_start)
	endI := startI + int(sbA.S_inodes_count)

	for j := startI; j < endI; j++ {
		fileA.Seek(int64(j), 0)
		errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
		if bit == '0' {
			bitI++
			break
		}
		bitI++
	}
	sbA.S_firts_ino = int32(bitI)
}

func crearInodoCarpeta(pos int, bloque int) {
	var newInodo Structs.TablaInodo
	newInodo.I_uid = UsuarioL.IdU
	newInodo.I_gid = UsuarioL.IdG
	newInodo.I_s = 0
	newInodo.I_atime = time.Now().Unix()
	newInodo.I_ctime = time.Now().Unix()
	newInodo.I_mtime = time.Now().Unix()
	newInodo.I_type = '0'
	newInodo.I_perm = 664
	newInodo.I_block[0] = int32(bloque)
	for j := 1; j < 16; j++ {
		newInodo.I_block[j] = -1
	}

	fileA.Seek(int64(pos), 0)
	var bufferByte bytes.Buffer
	errf := binary.Write(&bufferByte, binary.BigEndian, newInodo)
	if errf != nil {
		fmt.Println(errf)
	}
	EscribirFile(fileA, bufferByte.Bytes())
}

func crearInodoArchivo(pos int) {
	var newInodo Structs.TablaInodo
	newInodo.I_uid = UsuarioL.IdU
	newInodo.I_gid = UsuarioL.IdG
	newInodo.I_s = 0
	newInodo.I_atime = time.Now().Unix()
	newInodo.I_ctime = time.Now().Unix()
	newInodo.I_mtime = time.Now().Unix()
	newInodo.I_type = '1'
	newInodo.I_perm = 664
	for i := 0; i < 16; i++ {
		newInodo.I_block[i] = -1
	}
	fileA.Seek(int64(pos), 0)
	var bufferByte bytes.Buffer
	errf := binary.Write(&bufferByte, binary.BigEndian, newInodo)
	if errf != nil {
		fmt.Println(errf)
	}
	EscribirFile(fileA, bufferByte.Bytes())
}

func agregarCarpeta(posA int, posC int, name string) int {
	if !validarPermisoW(posC) {
		return -1
	}

	var inodo Structs.TablaInodo
	var carpeta Structs.BloqueCarpeta

	fileA.Seek(int64(posC), 0)
	errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

	if inodo.I_type == '1' {
		return -1
	}

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			fileA.Seek(int64(inodo.I_block[i]), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(carpeta))), binary.BigEndian, &carpeta)
			for c := 0; c < 4; c++ {
				if carpeta.B_content[c].B_inodo == -1 {
					carpeta.B_content[c].B_name = nameConten(name)
					carpeta.B_content[c].B_inodo = int32(posA)
					fileA.Seek(int64(inodo.I_block[i]), 0)
					var bufferByte bytes.Buffer
					errf := binary.Write(&bufferByte, binary.BigEndian, carpeta)
					if errf != nil {
						fmt.Println(errf)
					}
					EscribirFile(fileA, bufferByte.Bytes())
					return 0
				}
			}
		} else if inodo.I_block[i] == -1 {
			if sbA.S_free_blocks_count > 0 {
				posCarpetaO := crearBloqueCarpetaOtra(posA, name)
				inodo.I_block[i] = int32(posCarpetaO)
				fileA.Seek(int64(posC), 0)
				var bufferByte bytes.Buffer
				errf := binary.Write(&bufferByte, binary.BigEndian, inodo)
				if errf != nil {
					fmt.Println(errf)
				}
				EscribirFile(fileA, bufferByte.Bytes())
				return 0
			} else {
				return -1
			}
		}
	}
	return -1
}

func crearBloqueCarpetaInicial(posActual int, posPadre int) int {
	posBloque := buscarPosicionNewBLoque()

	var newCarpeta Structs.BloqueCarpeta
	newCarpeta.B_content[0].B_name = nameConten(".")
	newCarpeta.B_content[0].B_inodo = int32(posActual)
	newCarpeta.B_content[1].B_name = nameConten("..")
	newCarpeta.B_content[1].B_inodo = int32(posPadre)
	newCarpeta.B_content[2].B_name = nameConten("")
	newCarpeta.B_content[2].B_inodo = -1
	newCarpeta.B_content[3].B_name = nameConten("")
	newCarpeta.B_content[3].B_inodo = -1

	fileA.Seek(int64(posBloque), 0)
	var bufferByte bytes.Buffer
	errf := binary.Write(&bufferByte, binary.BigEndian, newCarpeta)
	if errf != nil {
		fmt.Println(errf)
	}
	EscribirFile(fileA, bufferByte.Bytes())

	return posBloque
}

func crearBloqueCarpetaOtra(hijo int, nombreH string) int {
	posBloque := buscarPosicionNewBLoque()

	var newCarpeta Structs.BloqueCarpeta
	newCarpeta.B_content[0].B_name = nameConten(nombreH)
	newCarpeta.B_content[0].B_inodo = int32(hijo)
	for i := 1; i < 4; i++ {
		newCarpeta.B_content[i].B_name = nameConten("")
		newCarpeta.B_content[i].B_inodo = -1
	}

	fileA.Seek(int64(posBloque), 0)
	var bufferByte bytes.Buffer
	errf := binary.Write(&bufferByte, binary.BigEndian, newCarpeta)
	if errf != nil {
		fmt.Println(errf)
	}
	EscribirFile(fileA, bufferByte.Bytes())

	return posBloque
}

func getSize() string {
	content := ""
	cont1 := 0
	cont2 := 0
	i := 0

	for i < SArchivos {
		content += strconv.Itoa(cont1)
		cont1++
		cont2++
		if cont1 == 10 {
			cont1 = 0
		}
		i++
		if cont2 == 19 {
			content += "\n"
			cont2 = 0
			i++
		}
	}

	return content
}

func validarPermisoW(posI int) bool {
	var inodo Structs.TablaInodo
	fileA.Seek(int64(posI), 0)
	errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)
	permiso := strconv.Itoa(int(inodo.I_perm))

	if len(permiso) == 1 {
		permiso = "00" + permiso
	} else if len(permiso) == 2 {
		permiso = "0" + permiso
	}

	if UsuarioL.IdU == 1 && UsuarioL.IdG == 1 {
		return true
	} else if inodo.I_uid == UsuarioL.IdU && inodo.I_gid == UsuarioL.IdG {
		if permiso[0] == '2' || permiso[0] == '3' || permiso[0] == '6' || permiso[0] == '7' {
			return true
		}
	} else if inodo.I_gid == UsuarioL.IdG {
		if permiso[1] == '2' || permiso[1] == '3' || permiso[1] == '6' || permiso[1] == '7' {
			return true
		}
	} else if permiso[2] == '2' || permiso[2] == '3' || permiso[2] == '6' || permiso[2] == '7' {
		return true
	}
	return false
}

func agregarArchivo2(cadena string, inodo Structs.TablaInodo, j int, aux int) Structs.TablaInodo {
	in := Structs.TablaInodo{}
	in.I_type = 'F'
	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 && i == j {
			var archivo Structs.BloqueArchivo
			fileA.Seek(int64(inodo.I_block[i]), 0)
			errfA = binary.Read(LeerFile(fileA, int(unsafe.Sizeof(archivo))), binary.BigEndian, &archivo)
			archivo.B_content = archivoContent(cadena)
			fileA.Seek(int64(inodo.I_block[i]), 0)
			var bufferArchivo bytes.Buffer
			errf := binary.Write(&bufferArchivo, binary.BigEndian, archivo)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileA, bufferArchivo.Bytes())
			return inodo
		} else if inodo.I_block[i] == -1 && aux != -1 {
			var archivo Structs.BloqueArchivo
			seek := int(sbA.S_block_start) + (aux * int(unsafe.Sizeof(archivo)))
			archivo.B_content = archivoContent(cadena)
			fileA.Seek(int64(seek), 0)
			var bufferArchivo bytes.Buffer
			errf := binary.Write(&bufferArchivo, binary.BigEndian, archivo)
			if errf != nil {
				fmt.Println(errf)
			}
			EscribirFile(fileA, bufferArchivo.Bytes())
			inodo.I_block[i] = int32(seek)
			CambioCont = true
			return inodo
		}
	}
	return in
}

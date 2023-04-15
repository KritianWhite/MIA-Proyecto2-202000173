package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

var Sdisk = 0
var Fdisk = "ff"
var Udisk = "m"
var Pdisk = " "
var Directorio_disk = ""

func MkDisk() Structs.Resp {
	val := validar()
	if !val.Val {
		return Structs.Resp{Res: val.Men}
	}

	var file *os.File
	var errf error
	//fmt.Println("\033[31m" + Pdisk + "\033[0m")
	file, errf = os.OpenFile(Pdisk, os.O_RDWR, 0666)
	defer func() {
		reco := recover()
		if reco != nil {
			fmt.Println(reco)
		}
		Sdisk = 0
		Fdisk = "ff"
		Udisk = "m"
		Pdisk = " "
		Directorio_disk = ""
		if file != nil {
			file.Close()
		}

	}()
	if errf == nil {
		return Structs.Resp{Res: "EL DISCO YA EXISTE"}
	}

	Directorio_disk = GetDirectorio(Pdisk)
	err := os.MkdirAll(Directorio_disk, 0777)
	if err != nil {
		fmt.Printf("%s", err)
	}

	size := Sdisk
	if Udisk == "m" {
		size = size * 1024
	}

	file, errf = os.OpenFile(Pdisk, os.O_RDWR|os.O_CREATE, 0664)

	var contenedor bytes.Buffer
	var buffer [1024]int8
	for i := 0; i < 1024; i++ {
		buffer[i] = 0
	}

	binary.Write(&contenedor, binary.BigEndian, &buffer)

	for i := 0; i < size; i++ {
		EscribirFile(file, contenedor.Bytes())
	}

	mbr := Structs.MBR{}
	mbr.Mbr_fecha_creacion = time.Now().Unix()
	mbr.Mbr_disk_signature = int32(int(binary.BigEndian.Uint64([]byte(time.Now().String()))))
	mbr.Mbr_tamanio = int32(size * 1024)
	mbr.Disk_fit = Fdisk[0]
	for j := 0; j < 4; j++ {
		mbr.Mbr_partition[j].Part_start = -1
	}

	file.Seek(0, 0)
	var bufferControl bytes.Buffer
	err = binary.Write(&bufferControl, binary.BigEndian, mbr)
	EscribirFile(file, bufferControl.Bytes())

	return Structs.Resp{Res: "SE CREO EL DISCO EXITOSAMENTE"}
}

func RmDisk() Structs.Resp {
	extension := ""
	if Pdisk != " " {
		i := find(Pdisk, ".")
		extension = Pdisk[i+1:]
		//fmt.Println(extension)
		if !strncmp(extension, "dsk") {
			return Structs.Resp{Res: "EXTENSION INCORRECTA"}
		}
	} else {
		return Structs.Resp{Res: "ASEGURESE DE ESCRIBIR UN RUTA"}
	}

	var file *os.File
	var errf error
	file, errf = os.OpenFile(Pdisk, os.O_RDWR, 0666)
	file, errf = os.OpenFile(Pdisk, os.O_RDWR, 0666)
	defer func() {
		reco := recover()
		if reco != nil {
			fmt.Println(reco)
		}
		Sdisk = 0
		Fdisk = "ff"
		Udisk = "m"
		Pdisk = " "
		Directorio_disk = ""
		if file != nil {
			file.Close()
		}

	}()

	if errf != nil {
		return Structs.Resp{Res: "NO EXISTE EL DISCO"}
	}

	err := os.Remove(Pdisk)
	if err != nil {
		return Structs.Resp{Res: "Error al eliminar el disco"}
	}

	return Structs.Resp{Res: "DISCO ELIMINADO"}
}

func validar() Structs.Bandera {
	if Sdisk > 0 {
		if strncmp(Fdisk, "bf") || strncmp(Fdisk, "ff") || strncmp(Fdisk, "wf") {
			if strncmp(Udisk, "k") || strncmp(Udisk, "m") {
				if Pdisk != " " {
					i := find(Pdisk, ".")
					extension := Pdisk[i+1:]
					if strncmp(extension, "dsk") {
						return Structs.Bandera{Val: true, Men: ""}
					} else {
						return Structs.Bandera{Val: false, Men: "EXTENSION INCORRECTA"}
					}
				} else {
					return Structs.Bandera{Val: false, Men: "ASEGURESE DE ESCRIBIR UN RUTA"}
				}
			} else {
				return Structs.Bandera{Val: false, Men: "CONFIGURACION DE UNIDADES DEL TAMAÑO DE MEMORIA INVALIDO"}
			}
		} else {
			return Structs.Bandera{Val: false, Men: "CONFIGURACION DE AJUSTE INVALIDO"}
		}
	} else {
		return Structs.Bandera{Val: false, Men: "EL TAMAÑO DEL DISCO TIENE QUE SER MAYOR A 0"}
	}
}

func GetDirectorio(path string) string {
	directorio := ""
	aux := path
	p := strings.Index(aux, "/")
	for p != -1 {
		directorio += aux[:p] + "/"
		aux = aux[p+1:]
		p = strings.Index(aux, "/")
	}

	return directorio
}

func GetExtension(path string) string {
	i := find(path, ".")
	aux := path[i+1:]
	return aux
}

func EscribirFile(file *os.File, data []byte) {
	_, err := file.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerFile(file *os.File, number int) *bytes.Buffer {
	data := make([]byte, number)
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	bufferD := bytes.NewBuffer(data)
	return bufferD
}

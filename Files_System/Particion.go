package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

var Spart = 0
var Upart byte = 'k'
var Ppart = " "
var Tpart byte = 'p'
var Fpart byte = 'w'
var Namepart = " "

func fdisk() Structs.Resp {
	defer func() {
		Spart = 0
		Upart = 'k'
		Ppart = " "
		Tpart = 'p'
		Fpart = 'w'
		Namepart = " "
	}()
	if Ppart != " " {
		if Namepart != " " {
			if Spart > 0 {
				if Upart == 'b' || Upart == 'k' || Upart == 'm' {
					if Tpart == 'p' {
						return particionPrimaria()
					} else if Tpart == 'e' {
						return particionExtendida()
					} else if Tpart == 'l' {
						return particionLogica()
					} else {
						return Structs.Resp{Res: "TIPO DE PARTICION INVALIDA"}
					}
				} else {
					return Structs.Resp{Res: "UNIDAD ERRONEA"}
				}
			} else {
				return Structs.Resp{Res: "ASEGURESE DE INGRESAR UN VALOR MAYOR A 0 CON -size"}
			}
		} else {
			return Structs.Resp{Res: "ASEGURESE EL NOMBRE DE LA PARTICION"}
		}
	} else {
		return Structs.Resp{Res: "ASEGURESE DE ESCRIBIR UN RUTA"}
	}
	return Structs.Resp{Res: "Algo salio mal"}
}

func particionPrimaria() Structs.Resp {
	p := Structs.Partition{}
	pos := -1
	file, errf := os.OpenFile(Ppart, os.O_RDWR, 0777)
	if errf == nil {
		file.Seek(0, 0)
		mbr := Structs.MBR{}
		errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)

		if errf != nil {
			fmt.Println(errf)
		}

		file.Close()
		for i := 0; i < 4; i++ {
			if mbr.Mbr_partition[i].Part_start == -1 {
				pos = i
				break
			}
		}

		if espacioDisponible(Spart, Ppart, Upart, pos) {
			if !existeNombreP(Ppart, Namepart) {
				p.Part_fit = Fpart
				p.Part_type = Tpart
				for i := 0; i < 16; i++ {
					if i == len(Namepart) {
						break
					}
					p.Part_name[i] = Namepart[i]
				}
				p.Part_status = '0'
				if Upart == 'b' {
					p.Part_s = int32(Spart)
				} else if Upart == 'k' {
					p.Part_s = int32(Spart) * 1024
				} else if Upart == 'm' {
					p.Part_s = int32(Spart) * 1024 * 1024
				}

				if pos == 0 {
					p.Part_start = int32(unsafe.Sizeof(Structs.MBR{}))
				} else {
					p.Part_start = mbr.Mbr_partition[pos-1].Part_start + mbr.Mbr_partition[pos-1].Part_s
				}

				mbr.Mbr_partition[pos] = p
				file, errf = os.OpenFile(Ppart, os.O_RDWR, 0777)
				file.Seek(0, 0)
				var bufferControl bytes.Buffer
				errf = binary.Write(&bufferControl, binary.BigEndian, mbr)
				EscribirFile(file, bufferControl.Bytes())
				file.Close()

				men := ""
				comprobacion := Structs.MBR{}
				file, errf = os.OpenFile(Ppart, os.O_RDWR, 0777)
				file.Seek(0, 0)
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(comprobacion))), binary.BigEndian, &comprobacion)
				file.Close()
				men += "SE CREO LA PARTICION # " + strconv.Itoa(pos+1) + "\n"
				men += "Particion: " + string(comprobacion.Mbr_partition[pos].Part_name[:]) + "\n"
				men += "Tipo: Primaria\n"
				men += "Inicio: " + strconv.Itoa(int(comprobacion.Mbr_partition[pos].Part_start)) + "\n"
				men += "Tamanio: " + strconv.Itoa(int(comprobacion.Mbr_partition[pos].Part_s))
				return Structs.Resp{Res: men}
			}
			return Structs.Resp{Res: "YA EXISTE LA PARTICION "}
		}
		return Structs.Resp{Res: "NO HAY SUFICIENTE ESPACIO PARA CREAR LA PARTICION"}
	} else {
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	return Structs.Resp{Res: "Algo salio mal"}
}

func particionExtendida() Structs.Resp {
	p := Structs.Partition{}
	pos, posEBR := -1, -1
	file, errf := os.OpenFile(Ppart, os.O_RDWR, 0777)
	if errf == nil {
		file.Seek(0, 0)
		mbr := Structs.MBR{}
		errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
		if errf != nil {
			fmt.Println(errf)
		}

		for i := 0; i < 4; i++ {
			if mbr.Mbr_partition[i].Part_start == -1 {
				pos = i
				break
			}
		}

		if espacioDisponible(Spart, Ppart, Upart, pos) {
			if !existeNombreP(Ppart, Namepart) {
				if !existeParticionExt(Ppart) {
					p.Part_fit = Fpart
					p.Part_type = Tpart
					for i := 0; i < 16; i++ {
						if i == len(Namepart) {
							break
						}
						p.Part_name[i] = Namepart[i]
					}
					p.Part_status = '0'
					if Upart == 'b' {
						p.Part_s = int32(Spart)
					} else if Upart == 'k' {
						p.Part_s = int32(Spart) * 1024
					} else if Upart == 'm' {
						p.Part_s = int32(Spart) * 1024 * 1024
					}
					if pos == 0 {
						p.Part_start = int32(unsafe.Sizeof(Structs.MBR{}))
					} else {
						p.Part_start = mbr.Mbr_partition[pos-1].Part_start + mbr.Mbr_partition[pos-1].Part_s
					}
					posEBR = int(p.Part_start)
					mbr.Mbr_partition[pos] = p
					file.Seek(0, 0)
					var bufferMBR bytes.Buffer
					errf = binary.Write(&bufferMBR, binary.BigEndian, mbr)
					EscribirFile(file, bufferMBR.Bytes())

					ebr := Structs.EBR{}
					ebr.Part_next = -1
					ebr.Part_start = int32(posEBR)
					ebr.Part_s = -1
					ebr.Part_status = '0'
					file.Seek(int64(posEBR), 0)
					var bufferEBR bytes.Buffer
					errf = binary.Write(&bufferEBR, binary.BigEndian, ebr)
					EscribirFile(file, bufferEBR.Bytes())
					file.Close()

					men := ""
					comprobacion := Structs.MBR{}
					file, errf = os.OpenFile(Ppart, os.O_RDWR, 0777)
					file.Seek(0, 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(comprobacion))), binary.BigEndian, &comprobacion)
					file.Close()
					men += "SE CREO LA PARTICION # " + strconv.Itoa(pos+1) + "\n"
					men += "Particion: " + string(comprobacion.Mbr_partition[pos].Part_name[:]) + "\n"
					men += "Tipo: Extendida\n"
					men += "Inicio: " + strconv.Itoa(int(comprobacion.Mbr_partition[pos].Part_start)) + "\n"
					men += "Tamanio: " + strconv.Itoa(int(comprobacion.Mbr_partition[pos].Part_s))
					return Structs.Resp{Res: men}

				}
				return Structs.Resp{Res: "YA EXISTE UNA PARTICION EXTENDIDA"}
			}
			return Structs.Resp{Res: "YA EXISTE LA PARTICION "}
		}
		return Structs.Resp{Res: "NO HAY SUFICIENTE ESPACIO PARA CREAR LA PARTICION"}
	} else {
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	fmt.Println(p, pos, posEBR)
	return Structs.Resp{Res: "Algo salio mal"}
}

func particionLogica() Structs.Resp {
	pos := -1
	file, errf := os.OpenFile(Ppart, os.O_RDWR, 0777)
	if errf == nil {
		if existeParticionExt(Ppart) {
			if !existeNombreP(Ppart, Namepart) {
				file.Seek(0, 0)
				mbr := Structs.MBR{}
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)

				for i := 0; i < 4; i++ {
					if mbr.Mbr_partition[i].Part_type == 'e' {
						pos = i
						break
					}
				}

				if pos != -1 {
					ebrAux := Structs.EBR{}
					espacioOcupado := 0
					file.Seek(int64(mbr.Mbr_partition[pos].Part_start), 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebrAux))), binary.BigEndian, &ebrAux)
					if ebrAux.Part_next != -1 || ebrAux.Part_s != -1 {
						espacioOcupado += int(unsafe.Sizeof(Structs.EBR{})) + int(ebrAux.Part_s)
						tell, _ := file.Seek(0, os.SEEK_CUR)
						for ebrAux.Part_next != -1 &&
							(int(tell) < int(mbr.Mbr_partition[pos].Part_start)+int(mbr.Mbr_partition[pos].Part_s)) {
							file.Seek(int64(ebrAux.Part_next), 0)
							errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebrAux))), binary.BigEndian, &ebrAux)
							espacioOcupado += int(unsafe.Sizeof(Structs.EBR{})) + int(ebrAux.Part_s)
						}
						var newE Structs.EBR
						newE.Part_fit = Fpart
						newE.Part_start = ebrAux.Part_start + int32(unsafe.Sizeof(Structs.EBR{})) + ebrAux.Part_s
						newE.Part_status = '0'
						newE.Part_next = -1
						if Upart == 'b' {
							newE.Part_s = int32(Spart)
						} else if Upart == 'k' {
							newE.Part_s = int32(Spart * 1024)
						} else if Upart == 'm' {
							newE.Part_s = int32(Spart * 1024 * 1024)
						}
						for i := 0; i < 16; i++ {
							if i == len(Namepart) {
								break
							}
							newE.Part_name[i] = Namepart[i]
						}
						espacioD := int(mbr.Mbr_partition[pos].Part_s) - espacioOcupado
						espacioNewE := int(unsafe.Sizeof(Structs.EBR{})) + int(newE.Part_s)
						if espacioD >= espacioNewE {
							ebrAux.Part_next = newE.Part_start
							file.Seek(int64(ebrAux.Part_start), 0)
							var bufferEBRA bytes.Buffer
							errf = binary.Write(&bufferEBRA, binary.BigEndian, ebrAux)
							EscribirFile(file, bufferEBRA.Bytes())

							file.Seek(int64(newE.Part_start), 0)
							var bufferEBRN bytes.Buffer
							errf = binary.Write(&bufferEBRN, binary.BigEndian, newE)
							EscribirFile(file, bufferEBRN.Bytes())
							file.Close()
							men := ""
							men += "SE CREO LA PARTICION \n"
							men += "Particion: " + string(newE.Part_name[:]) + "\n"
							men += "Tipo: Logica\n"
							men += "Inicio: " + strconv.Itoa(int(newE.Part_start)) + "\n"
							men += "Tamanio: " + strconv.Itoa(int(newE.Part_s))
							return Structs.Resp{Res: men}
						}
						return Structs.Resp{Res: "NO HAY ESPACIO PARA CREAR LA PARTICION LOGICA"}
					} else {
						ebrAux.Part_fit = Fpart
						ebrAux.Part_start = mbr.Mbr_partition[pos].Part_start
						ebrAux.Part_status = '0'
						if Upart == 'b' {
							ebrAux.Part_s = int32(Spart)
						} else if Upart == 'k' {
							ebrAux.Part_s = int32(Spart * 1024)
						} else if Upart == 'm' {
							ebrAux.Part_s = int32(Spart * 1024 * 1024)
						}
						ebrAux.Part_next = -1
						for i := 0; i < 16; i++ {
							if i == len(Namepart) {
								break
							}
							ebrAux.Part_name[i] = Namepart[i]
						}

						if mbr.Mbr_partition[pos].Part_s >= (ebrAux.Part_s + int32(unsafe.Sizeof(Structs.EBR{}))) {
							file.Seek(0, 0)
							var bufferMBR bytes.Buffer
							errf = binary.Write(&bufferMBR, binary.BigEndian, mbr)
							EscribirFile(file, bufferMBR.Bytes())

							file.Seek(int64(mbr.Mbr_partition[pos].Part_start), 0)
							var bufferEBR bytes.Buffer
							errf = binary.Write(&bufferEBR, binary.BigEndian, ebrAux)
							EscribirFile(file, bufferEBR.Bytes())
							file.Close()
							men := ""
							men += "SE CREO LA PARTICION \n"
							men += "Particion: " + string(ebrAux.Part_name[:]) + "\n"
							men += "Tipo: Logica\n"
							men += "Inicio: " + strconv.Itoa(int(ebrAux.Part_start)) + "\n"
							men += "Tamanio: " + strconv.Itoa(int(ebrAux.Part_s))
							return Structs.Resp{Res: men}
						}
						file.Close()
						return Structs.Resp{Res: "NO HAY ESPACIO PARA CREAR LA PARTICION LOGICA"}
					}
				}

			}
			return Structs.Resp{Res: "YA EXISTE LA PARTICION"}
		}
		return Structs.Resp{Res: "NO EXISTE UNA PARTICION EXTENDIDA PARA ALMACENAR LA PARTICION LOGICA"}
	}
	return Structs.Resp{Res: "DISCO INEXISTENTE"}
}

func espacioDisponible(s int, p string, u byte, pos int) bool {
	file, _ := os.OpenFile(p, os.O_RDWR, 0777)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	file.Seek(0, 0)
	mbr := Structs.MBR{}
	errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
	if errf != nil {
		fmt.Println(errf)
	}
	if pos > -1 {
		size := 0
		if u == 'b' {
			size = s
		} else if u == 'k' {
			size = s * 1024
		} else if u == 'm' {
			size = s * 1024 * 1024
		}
		if size > 0 {
			espacioRestante := 0
			if pos == 0 {
				espacioRestante = int(mbr.Mbr_tamanio) - int(unsafe.Sizeof(Structs.MBR{}))
			} else {
				espacioRestante = int(mbr.Mbr_tamanio) - int(mbr.Mbr_partition[pos-1].Part_start) - int(mbr.Mbr_partition[pos-1].Part_s)
			}
			return espacioRestante >= size
		}
		return false
	}
	return false
}

func existeNombreP(p string, name string) bool {
	file, _ := os.OpenFile(p, os.O_RDWR, 0777)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	file.Seek(0, 0)
	mbr := Structs.MBR{}
	errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
	if errf != nil {
		fmt.Println(errf)
	}
	for i := 0; i < 4; i++ {
		name1 := string(mbr.Mbr_partition[i].Part_name[:])
		if strncmp(name1, name) {
			file.Close()
			return true
		}
		if mbr.Mbr_partition[i].Part_type == 'e' {
			ebr := Structs.EBR{}
			file.Seek(int64(mbr.Mbr_partition[i].Part_start), 0)
			errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
			if errf != nil {
				fmt.Println(errf)
			}
			if ebr.Part_next != -1 || ebr.Part_s != -1 {
				name1 = string(ebr.Part_name[:])
				if strncmp(name, name1) {
					file.Close()
					return true
				}
				for ebr.Part_next != -1 {
					name1 = string(ebr.Part_name[:])
					if strncmp(name, name1) {
						file.Close()
						return true
					}
					file.Seek(int64(ebr.Part_next), 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
				}
			}

		}
	}
	return false
}

func existeParticionExt(p string) bool {
	file, _ := os.OpenFile(p, os.O_RDWR, 0777)
	file.Seek(0, 0)
	mbr := Structs.MBR{}
	errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
	if errf != nil {
		fmt.Println(errf)
	}
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'e' {
			return true
		}
	}
	return false
}

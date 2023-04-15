package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"
)

var Pmontar = " "
var Namemontar = " "
var IdMontar = " "
var Tmontar = "full"

func mount() Structs.Resp {
	defer func() {
		Pmontar = " "
		Namemontar = " "
		IdMontar = " "
		Tmontar = "full"
	}()

	if Pmontar != " " {
		if Namemontar != " " {
			pos := -1
			file, errf := os.OpenFile(Pmontar, os.O_RDWR, 0777)
			if errf == nil {
				file.Seek(0, 0)
				mbr := Structs.MBR{}
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
				for i := 0; i < 4; i++ {
					name1 := string(mbr.Mbr_partition[i].Part_name[:])
					if strncmp(name1, Namemontar) {
						pos = i
						break
					} else if mbr.Mbr_partition[i].Part_type == 'e' {
						ebr := Structs.EBR{}
						sb := Structs.SuperBloque{}
						file.Seek(int64(mbr.Mbr_partition[i].Part_start), 0)
						errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
						if errf != nil {
							fmt.Println(errf)
						}
						if ebr.Part_next != -1 || ebr.Part_s != -1 {
							name1 = string(ebr.Part_name[:])
							if strncmp(name1, Namemontar) {
								if ebr.Part_status == '0' || ebr.Part_status == '1' {
									ebr.Part_status = '1'
								}

								file.Seek(int64(ebr.Part_start), 0)
								var bufferEBRN bytes.Buffer
								errf = binary.Write(&bufferEBRN, binary.BigEndian, ebr)
								EscribirFile(file, bufferEBRN.Bytes())

								if ebr.Part_status == '2' {
									file.Seek(int64(ebr.Part_start)+int64(unsafe.Sizeof(Structs.EBR{})), 0)
									errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)
									sb.S_mtime = time.Now().Unix()
									sb.S_mnt_count += 1
									file.Seek(int64(ebr.Part_start)+int64(unsafe.Sizeof(Structs.EBR{})), 0)
									var bufferSB bytes.Buffer
									errf = binary.Write(&bufferSB, binary.BigEndian, sb)
									EscribirFile(file, bufferSB.Bytes())
								}
								file.Close()
								return Mlist.add(Pmontar, Namemontar, 'l', int(ebr.Part_start), -1)

							} else if ebr.Part_next != -1 {
								file.Seek(int64(ebr.Part_next), 0)
								errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
								for true {
									name1 = string(ebr.Part_name[:])
									if strncmp(name1, Namemontar) {
										if ebr.Part_status == '0' || ebr.Part_status == '1' {
											ebr.Part_status = '1'
										}

										file.Seek(int64(ebr.Part_start), 0)
										var bufferEBRN bytes.Buffer
										errf = binary.Write(&bufferEBRN, binary.BigEndian, ebr)
										EscribirFile(file, bufferEBRN.Bytes())

										if ebr.Part_status == '2' {
											file.Seek(int64(ebr.Part_start)+int64(unsafe.Sizeof(Structs.EBR{})), 0)
											errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)
											sb.S_mtime = time.Now().Unix()
											sb.S_mnt_count += 1
											file.Seek(int64(ebr.Part_start)+int64(unsafe.Sizeof(Structs.EBR{})), 0)
											var bufferSB bytes.Buffer
											errf = binary.Write(&bufferSB, binary.BigEndian, sb)
											EscribirFile(file, bufferSB.Bytes())
										}
										file.Close()
										return Mlist.add(Pmontar, Namemontar, 'l', int(ebr.Part_start), -1)

									}

									if ebr.Part_next == -1 {
										break
									}
									file.Seek(int64(ebr.Part_next), 0)
									errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
								}
							}
						}
					}
				}

				if pos != -1 {
					if mbr.Mbr_partition[pos].Part_type == 'e' {
						file.Close()
						return Structs.Resp{Res: "NO SE PUEDE MONTAR UNA PARTICION EXTENDIDA "}
					}
					sb := Structs.SuperBloque{}
					if mbr.Mbr_partition[pos].Part_status == '0' || mbr.Mbr_partition[pos].Part_status == '1' {
						mbr.Mbr_partition[pos].Part_status = '1'
					}

					file.Seek(0, 0)
					var bufferMBR bytes.Buffer
					errf = binary.Write(&bufferMBR, binary.BigEndian, mbr)
					EscribirFile(file, bufferMBR.Bytes())

					if mbr.Mbr_partition[pos].Part_status == '2' {
						file.Seek(int64(mbr.Mbr_partition[pos].Part_start), 0)
						errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)
						sb.S_mtime = time.Now().Unix()
						sb.S_mnt_count += 1
						file.Seek(int64(mbr.Mbr_partition[pos].Part_start), 0)
						var bufferSB bytes.Buffer
						errf = binary.Write(&bufferSB, binary.BigEndian, sb)
						EscribirFile(file, bufferSB.Bytes())
					}
					file.Close()
					return Mlist.add(Pmontar, Namemontar, 'p', int(mbr.Mbr_partition[pos].Part_start), pos)

				}
				file.Close()
				return Structs.Resp{Res: "NO EXISTE ESA PARTICION"}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "ASEGURESE DE ESCRIBIR EL NOMBRE DE LA PARTICION"}
	}
	return Structs.Resp{Res: "ASEGURESE DE ESCRIBIR UN RUTA"}
}

func mkfs() Structs.Resp {
	defer func() {
		Pmontar = " "
		Namemontar = " "
		IdMontar = " "
		Tmontar = "full"
	}()
	if IdMontar != " " {
		if Tmontar == "full" {
			nodoM := Mlist.buscar(IdMontar)
			if nodoM == nil {
				return Structs.Resp{Res: "NO EXISTE PARTCION MONTADA CON EL ID " + IdMontar}
			}
			file, errf := os.OpenFile(nodoM.Path, os.O_RDWR, 0777)
			if errf == nil {
				sb := Structs.SuperBloque{}
				inodo := Structs.TablaInodo{}
				carpeta := Structs.BloqueCarpeta{}

				if nodoM.Type == 'p' {
					mbr := Structs.MBR{}
					file.Seek(0, 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)

					tamanio := int(mbr.Mbr_partition[nodoM.Pos].Part_s)
					var n float64
					n = float64(tamanio-int(unsafe.Sizeof(Structs.SuperBloque{}))) / float64(4+int(unsafe.Sizeof(Structs.TablaInodo{}))+(3*int(unsafe.Sizeof(Structs.BloqueArchivo{}))))
					numEstructuras := int(math.Floor(n))
					nBloques := 3 * numEstructuras
					inoding := numEstructuras * int(unsafe.Sizeof(Structs.TablaInodo{}))

					sb.S_filesystem_type = 2
					sb.S_inodes_count = int32(numEstructuras)
					sb.S_blocks_count = int32(nBloques)
					sb.S_free_blocks_count = int32(nBloques - 2)
					sb.S_free_inodes_count = int32(numEstructuras - 2)
					sb.S_mtime = time.Now().Unix()
					sb.S_umtime = 0
					sb.S_mnt_count = 1
					sb.S_magic = 0xEF53
					sb.S_inode_s = int32(unsafe.Sizeof(Structs.TablaInodo{}))
					sb.S_block_s = int32(unsafe.Sizeof(Structs.BloqueArchivo{}))
					sb.S_firts_ino = 2
					sb.S_first_blo = 2
					sb.S_bm_inode_start = int32(nodoM.Start) + int32(unsafe.Sizeof(Structs.SuperBloque{}))
					sb.S_bm_block_start = sb.S_bm_inode_start + int32(numEstructuras)
					sb.S_inode_start = sb.S_bm_block_start + int32(nBloques)
					sb.S_block_start = sb.S_inode_start + int32(inoding)

					inodo.I_uid = 1
					inodo.I_gid = 1
					inodo.I_atime = time.Now().Unix()
					inodo.I_ctime = time.Now().Unix()
					inodo.I_mtime = time.Now().Unix()
					inodo.I_perm = 664
					inodo.I_block[0] = sb.S_block_start
					for i := 1; i < 16; i++ {
						inodo.I_block[i] = -1
					}
					inodo.I_type = '0'
					inodo.I_s = 0
					file.Seek(int64(sb.S_inode_start), 0)
					var bufferInode bytes.Buffer
					errf = binary.Write(&bufferInode, binary.BigEndian, inodo)
					EscribirFile(file, bufferInode.Bytes())

					carpeta.B_content[0].B_name = nameConten(".")
					carpeta.B_content[0].B_inodo = sb.S_inode_start
					carpeta.B_content[1].B_name = nameConten("..")
					carpeta.B_content[1].B_inodo = sb.S_inode_start
					carpeta.B_content[2].B_name = nameConten("users.txt")
					carpeta.B_content[2].B_inodo = sb.S_inode_start + int32(unsafe.Sizeof(Structs.TablaInodo{}))
					carpeta.B_content[3].B_name = nameConten("")
					carpeta.B_content[3].B_inodo = -1
					file.Seek(int64(sb.S_block_start), 0)
					var bufferCarpeta bytes.Buffer
					errf = binary.Write(&bufferCarpeta, binary.BigEndian, carpeta)
					EscribirFile(file, bufferCarpeta.Bytes())

					inodoU := Structs.TablaInodo{}
					archivoU := Structs.BloqueArchivo{}

					inodoU.I_uid = 1
					inodoU.I_gid = 1
					inodoU.I_atime = time.Now().Unix()
					inodoU.I_ctime = time.Now().Unix()
					inodoU.I_mtime = time.Now().Unix()
					inodoU.I_perm = 700
					inodoU.I_block[0] = sb.S_block_start + int32(unsafe.Sizeof(Structs.BloqueCarpeta{}))
					for i := 1; i < 16; i++ {
						inodoU.I_block[i] = -1
					}
					s := "1,G,root\n1,U,root,root,123\n"
					inodoU.I_s = int32(unsafe.Sizeof(s))
					inodoU.I_type = '1'
					file.Seek(int64(sb.S_inode_start+int32(unsafe.Sizeof(Structs.TablaInodo{}))), 0)

					var bufferInodo2 bytes.Buffer
					errf = binary.Write(&bufferInodo2, binary.BigEndian, inodoU)
					EscribirFile(file, bufferInodo2.Bytes())

					for i := 0; i < 64; i++ {
						archivoU.B_content[i] = '\000'
					}
					archivoU.B_content = archivoContent(s)
					file.Seek(int64(sb.S_block_start+int32(unsafe.Sizeof(Structs.BloqueCarpeta{}))), 0)
					var bufferArchivo bytes.Buffer
					errf = binary.Write(&bufferArchivo, binary.BigEndian, archivoU)
					EscribirFile(file, bufferArchivo.Bytes())

					mbr.Mbr_partition[nodoM.Pos].Part_status = '2'
					file.Seek(0, 0)
					var bufferMBR bytes.Buffer
					errf = binary.Write(&bufferMBR, binary.BigEndian, mbr)
					EscribirFile(file, bufferMBR.Bytes())
					file.Seek(int64(mbr.Mbr_partition[nodoM.Pos].Part_start), 0)
					var bufferSB bytes.Buffer
					errf = binary.Write(&bufferSB, binary.BigEndian, sb)
					EscribirFile(file, bufferSB.Bytes())

					var ch0 byte = '0'
					var ch1 byte = '1'

					for i := 0; i < numEstructuras; i++ {

						file.Seek(int64(int(sb.S_bm_inode_start)+i), 0)
						var bufferC0 bytes.Buffer
						errf = binary.Write(&bufferC0, binary.BigEndian, ch0)
						EscribirFile(file, bufferC0.Bytes())
					}
					var bufferC1 bytes.Buffer
					errf = binary.Write(&bufferC1, binary.BigEndian, ch1)
					file.Seek(int64(sb.S_bm_inode_start), 0)
					EscribirFile(file, bufferC1.Bytes())
					file.Seek(int64(int(sb.S_bm_inode_start)+1), 0)
					EscribirFile(file, bufferC1.Bytes())

					for i := 0; i < nBloques; i++ {
						file.Seek(int64(int(sb.S_bm_block_start)+i), 0)
						var bufferC0 bytes.Buffer
						errf = binary.Write(&bufferC0, binary.BigEndian, ch0)
						EscribirFile(file, bufferC0.Bytes())
					}
					file.Seek(int64(sb.S_bm_block_start), 0)
					EscribirFile(file, bufferC1.Bytes())
					file.Seek(int64(int(sb.S_bm_block_start)+1), 0)
					EscribirFile(file, bufferC1.Bytes())

					file.Close()
					return Structs.Resp{Res: "SE FORMATEO LA PARTICION CON EXITO"}

				} else if nodoM.Type == 'l' {
					ebr := Structs.EBR{}
					file.Seek(int64(nodoM.Start), 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)

					tamanio := int(ebr.Part_s)
					var n float64
					n = float64(tamanio-int(unsafe.Sizeof(Structs.SuperBloque{}))) / float64(4+int(unsafe.Sizeof(Structs.TablaInodo{}))+(3*int(unsafe.Sizeof(Structs.BloqueArchivo{}))))
					numEstructuras := int(math.Floor(n))
					nBloques := 3 * numEstructuras
					inoding := numEstructuras * int(unsafe.Sizeof(Structs.TablaInodo{}))

					sb.S_filesystem_type = 2
					sb.S_inodes_count = int32(numEstructuras)
					sb.S_blocks_count = int32(nBloques)
					sb.S_free_blocks_count = int32(nBloques - 2)
					sb.S_free_inodes_count = int32(numEstructuras - 2)
					sb.S_mtime = time.Now().Unix()
					sb.S_umtime = 0
					sb.S_mnt_count = 1
					sb.S_magic = 0xEF53
					sb.S_inode_s = int32(unsafe.Sizeof(Structs.TablaInodo{}))
					sb.S_block_s = int32(unsafe.Sizeof(Structs.BloqueArchivo{}))
					sb.S_firts_ino = 2
					sb.S_first_blo = 2
					sb.S_bm_inode_start = int32(nodoM.Start) + int32(unsafe.Sizeof(Structs.EBR{})) + int32(unsafe.Sizeof(Structs.SuperBloque{}))
					sb.S_bm_block_start = sb.S_bm_inode_start + int32(numEstructuras)
					sb.S_inode_start = sb.S_bm_block_start + int32(nBloques)
					sb.S_block_start = sb.S_inode_start + int32(inoding)

					inodo.I_uid = 1
					inodo.I_gid = 1
					inodo.I_atime = time.Now().Unix()
					inodo.I_ctime = time.Now().Unix()
					inodo.I_mtime = time.Now().Unix()
					inodo.I_perm = 664
					inodo.I_block[0] = sb.S_block_start
					for i := 1; i < 16; i++ {
						inodo.I_block[i] = -1
					}
					inodo.I_type = '0'
					inodo.I_s = 0
					file.Seek(int64(sb.S_inode_start), 0)
					var bufferInode bytes.Buffer
					errf = binary.Write(&bufferInode, binary.BigEndian, inodo)
					EscribirFile(file, bufferInode.Bytes())

					carpeta.B_content[0].B_name = nameConten(".")
					carpeta.B_content[0].B_inodo = sb.S_inode_start
					carpeta.B_content[1].B_name = nameConten("..")
					carpeta.B_content[1].B_inodo = sb.S_inode_start
					carpeta.B_content[2].B_name = nameConten("users.txt")
					carpeta.B_content[2].B_inodo = sb.S_inode_start + int32(unsafe.Sizeof(Structs.TablaInodo{}))
					carpeta.B_content[3].B_name = nameConten("")
					carpeta.B_content[3].B_inodo = -1
					file.Seek(int64(sb.S_block_start), 0)
					var bufferCarpeta bytes.Buffer
					errf = binary.Write(&bufferCarpeta, binary.BigEndian, carpeta)
					EscribirFile(file, bufferCarpeta.Bytes())

					inodoU := Structs.TablaInodo{}
					archivoU := Structs.BloqueArchivo{}

					inodoU.I_uid = 1
					inodoU.I_gid = 1
					inodoU.I_atime = time.Now().Unix()
					inodoU.I_ctime = time.Now().Unix()
					inodoU.I_mtime = time.Now().Unix()
					inodoU.I_perm = 700
					inodoU.I_block[0] = sb.S_block_start + int32(unsafe.Sizeof(Structs.BloqueCarpeta{}))
					for i := 1; i < 16; i++ {
						inodoU.I_block[i] = -1
					}
					s := "1,G,root\n1,U,root,root,123\n"
					inodoU.I_s = int32(unsafe.Sizeof(s))
					inodoU.I_type = '1'
					file.Seek(int64(sb.S_inode_start+int32(unsafe.Sizeof(Structs.TablaInodo{}))), 0)
					var bufferInodo2 bytes.Buffer
					errf = binary.Write(&bufferInodo2, binary.BigEndian, inodoU)
					EscribirFile(file, bufferInodo2.Bytes())

					for i := 0; i < 64; i++ {
						archivoU.B_content[i] = '\000'
					}
					archivoU.B_content = archivoContent(s)
					file.Seek(int64(sb.S_block_start+int32(unsafe.Sizeof(Structs.BloqueCarpeta{}))), 0)
					var bufferArchivo bytes.Buffer
					errf = binary.Write(&bufferArchivo, binary.BigEndian, archivoU)
					EscribirFile(file, bufferArchivo.Bytes())

					ebr.Part_status = '2'
					file.Seek(int64(nodoM.Start), 0)
					var bufferEBR bytes.Buffer
					errf = binary.Write(&bufferEBR, binary.BigEndian, ebr)
					EscribirFile(file, bufferEBR.Bytes())
					file.Seek(int64(int(ebr.Part_start)+int(unsafe.Sizeof(Structs.EBR{}))), 0)
					var bufferSB bytes.Buffer
					errf = binary.Write(&bufferSB, binary.BigEndian, sb)
					EscribirFile(file, bufferSB.Bytes())

					ch0 := '0'
					ch1 := '1'

					for i := 0; i < numEstructuras; i++ {
						file.Seek(int64(int(sb.S_bm_inode_start)+i), 0)
						var bufferC0 bytes.Buffer
						errf = binary.Write(&bufferC0, binary.BigEndian, ch0)
						EscribirFile(file, bufferC0.Bytes())
					}
					var bufferC1 bytes.Buffer
					errf = binary.Write(&bufferC1, binary.BigEndian, ch1)
					file.Seek(int64(sb.S_bm_inode_start), 0)
					EscribirFile(file, bufferC1.Bytes())
					file.Seek(int64(int(sb.S_bm_inode_start)+1), 0)
					EscribirFile(file, bufferC1.Bytes())

					for i := 0; i < nBloques; i++ {
						file.Seek(int64(int(sb.S_bm_block_start)+i), 0)
						var bufferC0 bytes.Buffer
						errf = binary.Write(&bufferC0, binary.BigEndian, ch0)
						EscribirFile(file, bufferC0.Bytes())
					}
					file.Seek(int64(sb.S_bm_block_start), 0)
					EscribirFile(file, bufferC1.Bytes())
					file.Seek(int64(int(sb.S_bm_block_start)+1), 0)
					EscribirFile(file, bufferC1.Bytes())

					file.Close()
					return Structs.Resp{Res: "SE FORMATEO LA PARTICION CON EXITO"}
				}
			}
			return Structs.Resp{Res: "ERROR, NO SE PUDO ENCONTRAR EL DISCO DURO"}
		}
		return Structs.Resp{Res: "TIPO DE FORMATEO EQUIVOCADO"}
	}
	return Structs.Resp{Res: "SE NECESITA UN ID DE LA MONTURA"}
}

func nameConten(cadena string) [12]byte {
	var name [12]byte

	for i := 0; i < 12; i++ {
		if i >= len(cadena) {
			break
		}
		name[i] = cadena[i]
	}

	return name
}

func nameConten2(array [12]byte) string {
	bname := ""

	for j := 0; j < len(array); j++ {
		if array[j] == '\000' {
			break
		}
		bname += string([]byte{array[j]})
	}

	return bname
}

func archivoContent(cadena string) [64]byte {
	var content [64]byte

	for i := 0; i < 64; i++ {
		if i >= len(cadena) {
			break
		}
		content[i] = cadena[i]
	}

	return content
}

func archivoContent2(array [64]byte) string {
	var cadena string

	for i := 0; i < 64; i++ {
		if array[i] == '\000' {
			break
		}
		cadena += string([]byte{array[i]})
	}

	return cadena
}

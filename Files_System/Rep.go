package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var Prep = " "
var Namerep = " "
var Idrep = " "
var Rutarep = " "
var Dirrep = " "
var Extrep = " "

func GenerateRep() Structs.Resp {
	defer func() {
		Prep = " "
		Namerep = " "
		Idrep = " "
		Rutarep = " "
		Dirrep = " "
		Extrep = " "
	}()
	if Prep != " " {
		if Idrep != " " {
			if Namerep == "disk" {
				return disk()
			} else if Namerep == "tree" {
				return tree()
			} else if Namerep == "file" {
				return fileR()
			} else if Namerep == "sb" {
				return sbR()
			}
			return Structs.Resp{Res: "NOMBRE DE REPORTE INVALIDO"}
		}
		return Structs.Resp{Res: "FALTA EL ID DE LA PARTICION"}
	}
	return Structs.Resp{Res: "FALTA LA UBICACION DONDE SE GUARDARA EL REPORTE"}
}


func disk() Structs.Resp {
	nodo := Mlist.buscar(Idrep)
	if nodo != nil {
		Dirrep = GetDirectorio(Prep)
		Extrep = GetExtension(Prep)
		nombreD := nombre(Prep)
		err := os.MkdirAll(Dirrep, 0777)
		if err != nil {
			fmt.Printf("%s", err)
		}
		file, errf := os.OpenFile(nodo.Path, os.O_RDWR, 0777)
		if errf == nil {
			var mbr Structs.MBR
			file.Seek(0, 0)
			errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
			tamanioT := int(mbr.Mbr_tamanio)
			dotS := ""
			dot, errD := os.OpenFile("Reportes/"+nombreD+".dot", os.O_CREATE, 0777)
			dot.Close()
			if errD != nil {
				fmt.Println(errD)
			}

			dotS += "digraph G {\n"
			dotS += "node[shape=none]\n"
			dotS += "start[label=<<table color='orange'><tr>\n"
			dotS += "<td height='30' width='75' rowspan=\"2\">MBR</td>\n"

			i := 0
			inicio := int(unsafe.Sizeof(Structs.MBR{}))
			for i < 4 {
				var parcial = mbr.Mbr_partition[i].Part_s
				if mbr.Mbr_partition[i].Part_start != -1 {
					var porcentaje_real = (float64(parcial) / float64(tamanioT)) * 100
					//var porcentaje_aux = float64(int(porcentaje_real*100)) / 100
					if mbr.Mbr_partition[i].Part_type == 'p' {
						porcentaje := (float64(mbr.Mbr_partition[i].Part_s) / float64(tamanioT)) * 100
						trunc := float64(int(porcentaje*100)) / 100
						name1 := getPartName(mbr.Mbr_partition[i].Part_name)
						dotS += "<td rowspan=\"2\" width='" + fmt.Sprintf("%v", porcentaje_real) +"'>PRIMARIA <br/> " + name1 + " <br/>Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
						//* Verificar que no haya espacio fragmentado
						if i != 3 {
							if (mbr.Mbr_partition[i].Part_start + mbr.Mbr_partition[i].Part_s) < mbr.Mbr_partition[i+1].Part_start {
								porcentaje = (float64(mbr.Mbr_partition[i+1].Part_start-(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s)) / float64(tamanioT)) * 100
								trunc = float64(int(porcentaje*100)) / 100
								dotS += "<td rowspan=\"2\">LIBRE <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
							}
						} else if int(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s) < tamanioT {
							porcentaje = (float64(tamanioT-int(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s)) / float64(tamanioT)) * 100
							trunc = float64(int(porcentaje*100)) / 100
							dotS += "<td rowspan=\"2\">LIBRE <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
						}
						//*Extendida
					} else if mbr.Mbr_partition[i].Part_type == 'e' {
						porcentaje := (float64(mbr.Mbr_partition[i].Part_s) / float64(tamanioT)) * 100
						dotS += "<td rowspan=\"2\" height='60' colspan='100%'>EXTENDIDA</td>\n"
						ebr := Structs.EBR{}
						file.Seek(int64(mbr.Mbr_partition[i].Part_start), 0)
						errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
						if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
							if ebr.Part_s > -1 {
								name1 := getPartName(ebr.Part_name)
								dotS += "<td rowspan=\"2\" height='30'>EBR <br/>" + name1 + "</td>\n"
								porcentaje = (float64(ebr.Part_s) / float64(tamanioT)) * 100.0
								trunc := float64(int(porcentaje*100)) / 100
								dotS += "<td rowspan=\"2\" height='30'>Logica <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
							} else {
								dotS += "<td rowspan=\"2\" height='30'>EBR</td>\n"
								porcentaje = ((float64(ebr.Part_next - ebr.Part_start)) / float64(tamanioT)) * 100.0
								trunc := float64(int(porcentaje*100)) / 100
								dotS += "<td rowspan=\"2\" height='30'>Libre <br/> Ocupa: " + fmt.Sprintf("%v", trunc) + "%</td>\n"
							}
							if ebr.Part_next != -1 {
								file.Seek(int64(ebr.Part_next), 0)
								errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
								for true {
									name1 := getPartName(ebr.Part_name)
									dotS += "<td rowspan=\"2\">EBR <br/>" + name1 + "</td>\n"
									porcentaje = (float64(ebr.Part_s) / float64(tamanioT)) * 100.0
									trunc := float64(int(porcentaje*100)) / 100
									dotS += "<td rowspan=\"2\">Logica <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"

									if ebr.Part_next == -1 {
										if (ebr.Part_start + ebr.Part_s) < mbr.Mbr_partition[i].Part_s {
											porcentaje = (float64(mbr.Mbr_partition[i].Part_s-(ebr.Part_start+ebr.Part_s)) / float64(tamanioT)) * 100
											trunc = float64(int(porcentaje*100)) / 100
											dotS += "<td rowspan=\"2\">Libre <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
										}
										break
									}
									if (ebr.Part_start + ebr.Part_s) < ebr.Part_next {
										porcentaje = (float64(ebr.Part_next-(ebr.Part_start+ebr.Part_s)) / float64(tamanioT)) * 100
										trunc = float64(int(porcentaje*100)) / 100
										dotS += "<td rowspan=\"2\">Libre <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
									}
									file.Seek(int64(ebr.Part_next), 0)
									errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
								}
							}
						}
						dotS += "<td rowspan=\"2\" height='60' colspan='100%'>EXTENDIDA</td>\n"
					}
					inicio = int(mbr.Mbr_partition[i].Part_start + mbr.Mbr_partition[i].Part_s)
				} else {
					i++
					for i < 4 {
						if mbr.Mbr_partition[i].Part_start != -1 {
							porcentaje := (float64(int(mbr.Mbr_partition[i].Part_start)-inicio) / float64(tamanioT)) * 100
							trunc := float64(int(porcentaje*100)) / 100
							dotS += "<td rowspan=\"2\" height='30' width='100%'>Libre <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
							break
						}
						i++
					}
					if i == 4 {
						porcentaje := float64(tamanioT-inicio) / float64(tamanioT) * 100
						trunc := float64(int(porcentaje*100)) / 100
						dotS += "<td rowspan=\"2\" height='30' width='100%'>Libre <br/> Ocupa:" + fmt.Sprintf("%v", trunc) + "%</td>\n"
						goto t0
					}
					i--
				}
				i++
			}
		t0:
			
			dotS += "</tr></table>>];\n"
			dotS += "}"
			errD = os.WriteFile("Reportes/"+nombreD+".dot", []byte(dotS), 0777)
			if errD != nil {
				fmt.Println(errD)
			}

			file.Close()
			ext := Extrep
			_, errD = exec.Command("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD+"."+ext).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}
			_, errD = exec.Command("dot", "-T"+ext, "Reportes/"+nombreD+".dot", "-o", Dirrep+nombreD+"."+ext).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}
			//fmt.Println("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD+"."+ext)
			//fmt.Println("dot", "-T"+ext, "Reportes/"+nombreD+".dot", "-o", Dirrep+nombreD+"."+ext)
			
			return Structs.Resp{Res: "SE GENERO EL REPORTE DISK", Dot1:dotS}
		}
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + Idrep}

}

func tree() Structs.Resp {
	nodo := Mlist.buscar(Idrep)
	dotS := ""
	if nodo != nil {
		Dirrep = GetDirectorio(Prep)
		Extrep = GetExtension(Prep)
		nombreD := nombre(Prep)
		err := os.MkdirAll(Dirrep, 0777)
		if err != nil {
			fmt.Printf("%s", err)
		}

		file, errf := os.OpenFile(nodo.Path, os.O_RDWR, 0777)
		if errf == nil {
			sb := Structs.SuperBloque{}
			if nodo.Type == 'p' {
				file.Seek(int64(nodo.Start), 0)
			} else if nodo.Type == 'l' {
				file.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
			}
			errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)

			start := int(sb.S_bm_inode_start)
			end := start + int(sb.S_inodes_count)

			inodo := Structs.TablaInodo{}

			var bit byte
			cont := 0

			dot, errD := os.OpenFile("Reportes/"+nombreD+".dot", os.O_CREATE, 0777)
			dot.Close()
			if errD != nil {
				fmt.Println(errD)
			}

			
			dotS += "digraph G {\n"
			dotS += "rankdir=LR;\n"
			dotS += "node[shape=none]\n"

			for i := start; i < end; i++ {
				file.Seek(int64(i), 0)
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(bit))), binary.BigEndian, &bit)
				if bit == '1' {
					posInodo := int(sb.S_inode_start) + (cont * int(unsafe.Sizeof(Structs.TablaInodo{})))
					file.Seek(int64(posInodo), 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(Structs.TablaInodo{}))), binary.BigEndian, &inodo)

					dotS += treeInodo(posInodo, file)

					for j := 0; j < 16; j++ {
						if inodo.I_block[j] != -1 {
							if inodo.I_type == '0' {
								dotS += treeBlock(int(inodo.I_block[j]), 0, file)
							} else if inodo.I_type == '1' {
								dotS += treeBlock(int(inodo.I_block[j]), 1, file)
							}
							dotS += conexiones(posInodo, int(inodo.I_block[j]))
						}
					}
				}
				cont++
			}
			dotS += "}"
			errD = os.WriteFile("Reportes/"+nombreD+".dot", []byte(dotS), 0777)
			if errD != nil {
				fmt.Println(errD)
			}

			file.Close()

			ext := Extrep
			_, errD = exec.Command("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD + "." + ext).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}
			_, errD = exec.Command("dot", "-T"+ext, "Reportes/"+nombreD+".dot", "-o", Dirrep+nombreD + "." + ext ).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}
			//mt.Println("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD + "." + ext)
			//fmt.Println("dot", "-T"+ext, "Reportes/"+nombreD+".dot", "-o", Dirrep+nombreD + "." + ext)
			return Structs.Resp{Res: "SE GENERO EL REPORTE TREE"}
		}
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + Idrep, Dot3:dotS}
}

func treeBlock(pos int, typ int, file *os.File) string {
	dot := ""

	if typ == 0 {
		carpeta := Structs.BloqueCarpeta{}
		file.Seek(int64(pos), 0)
		errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(Structs.BloqueCarpeta{}))), binary.BigEndian, &carpeta)
		if errf != nil {
			fmt.Println(errf)
		}

		dot += "n" + strconv.Itoa(pos) + "[label=<<table>\n"
		dot += "<tr>\n"
		dot += "<td colspan=\"2\" bgcolor=\"#f34037\">Bloque Carpeta</td>"
		dot += "</tr>\n"

		for i := 0; i < 4; i++ {
			bname := nameConten2(carpeta.B_content[i].B_name)
			dot += "<tr>\n"
			dot += "<td>" + bname + "</td>\n"
			dot += "<td port=\"" + strconv.Itoa(int(carpeta.B_content[i].B_inodo)) + "\">" + strconv.Itoa(int(carpeta.B_content[i].B_inodo)) + "</td>\n"
			dot += "</tr>\n"
		}
		dot += "</table>>]\n"

		for i := 0; i < 4; i++ {
			name1 := nameConten2(carpeta.B_content[i].B_name)
			if carpeta.B_content[i].B_inodo != -1 && (name1 != "." && name1 != "..") {
				dot += conexiones(pos, int(carpeta.B_content[i].B_inodo))
			}
		}

	} else if typ == 1 {
		content := ""
		archivo := Structs.BloqueArchivo{}
		file.Seek(int64(pos), 0)
		errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(Structs.BloqueArchivo{}))), binary.BigEndian, &archivo)
		if errf != nil {
			fmt.Println(errf)
		}
		content = archivoContent2(archivo.B_content)

		dot += "n" + strconv.Itoa(pos) + "[label=<<table>\n"
		dot += "<tr>\n"
		dot += "<td colspan=\"2\" bgcolor=\"#c3f8b6\">Bloque Archivo</td>"
		dot += "</tr>\n"
		dot += "<tr>\n"
		dot += "<td>" + content + "</td>\n"
		dot += "</tr>\n"
		dot += "</table>>]\n"
	}
	return dot
}

func treeInodo(pos int, file *os.File) string {
	dot := ""
	inodo := Structs.TablaInodo{}
	file.Seek(int64(pos), 0)
	errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(Structs.TablaInodo{}))), binary.BigEndian, &inodo)
	if errf != nil {
		fmt.Println(errf)
	}
	dot += "n" + strconv.Itoa(pos) + "[label=<<table><tr><td colspan=\"2\" bgcolor=\"#376ef3\">INODO " + strconv.Itoa(pos) + "</td></tr>\n"

	dot += "<tr>\n"
	dot += "<td>i_uid</td>\n"
	dot += "<td>" + strconv.Itoa(int(inodo.I_uid)) + "</td>\n"
	dot += "</tr>\n"

	dot += "<tr>\n"
	dot += "<td>i_gid</td>\n"
	dot += "<td>" + strconv.Itoa(int(inodo.I_gid)) + "</td>\n"
	dot += "</tr>\n"

	dot += "<tr>\n"
	dot += "<td>i_s</td>\n"
	dot += "<td>" + strconv.Itoa(int(inodo.I_s)) + "</td>\n"
	dot += "</tr>\n"

	tm := time.Unix(inodo.I_atime, 0)
	dot += "<tr>\n"
	dot += "<td>i_atime</td>\n"
	dot += "<td>" + tm.Format("2006-01-02 15:04:05") + "</td>\n"
	dot += "</tr>\n"

	tm = time.Unix(inodo.I_ctime, 0)
	dot += "<tr>\n"
	dot += "<td>i_ctime</td>\n"
	dot += "<td>" + tm.Format("2006-01-02 15:04:05") + "</td>\n"
	dot += "</tr>\n"

	tm = time.Unix(inodo.I_mtime, 0)
	dot += "<tr>\n"
	dot += "<td>i_mtime</td>\n"
	dot += "<td>" + tm.Format("2006-01-02 15:04:05") + "</td>\n"
	dot += "</tr>\n"

	for j := 0; j < 16; j++ {
		if inodo.I_block[j] != -1 {
			dot += "<tr>\n"
			dot += "<td>ap" + strconv.Itoa(j) + "</td>\n"
			dot += "<td port=\"" + strconv.Itoa(int(inodo.I_block[j])) + "\">" + strconv.Itoa(int(inodo.I_block[j])) + "</td>\n"
			dot += "</tr>\n"
		} else {
			dot += "<tr>\n"
			dot += "<td>i_block</td>\n"
			dot += "<td>-1</td>\n"
			dot += "</tr>\n"
		}
	}

	dot += "<tr>\n"
	dot += "<td>i_type</td>\n"
	dot += "<td>" + string(inodo.I_type) + "</td>\n"
	dot += "</tr>\n"

	dot += "<tr>\n"
	dot += "<td>i_perm</td>\n"
	dot += "<td>" + strconv.Itoa(int(inodo.I_perm)) + "</td>\n"
	dot += "</tr>\n"

	dot += "</table>>]\n"

	return dot
}

func conexiones(inicio int, final int) string {
	dot := "n" + strconv.Itoa(inicio) + ":" + strconv.Itoa(final) + "->n" + strconv.Itoa(final) + ";\n"
	return dot
}

func sbR() Structs.Resp {
	nodo := Mlist.buscar(Idrep)
	if nodo != nil {
		Dirrep = GetDirectorio(Prep)
		Extrep = GetExtension(Prep)
		nombreD := nombre(Prep)
		err := os.MkdirAll(Dirrep, 0777)
		if err != nil {
			fmt.Printf("%s", err)
		}
		file, errf := os.OpenFile(nodo.Path, os.O_RDWR, 0777)
		if errf == nil {
			sb := Structs.SuperBloque{}
			if nodo.Type == 'p' {
				mbr := Structs.MBR{}
				file.Seek(0, 0)
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
				if mbr.Mbr_partition[nodo.Pos].Part_status != '2' {
					file.Close()
					return Structs.Resp{Res: "NO SE HA FORMATEADO LA MONTURA DE LA PARTICION " + nodo.Name}
				}
				file.Seek(int64(nodo.Start), 0)
			} else if nodo.Type == 'l' {
				ebr := Structs.EBR{}
				file.Seek(int64(nodo.Start), 0)
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
				if ebr.Part_status != '2' {
					file.Close()
					return Structs.Resp{Res: "NO SE HA FORMATEADO LA MONTURA DE LA PARTICION " + nodo.Name}
				}
				file.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
			}
			errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)
			file.Close()

			dot, errD := os.OpenFile("Reportes/"+nombreD+".dot", os.O_CREATE, 0777)
			dot.Close()
			if errD != nil {
				fmt.Println(errD)
			}

			dotS := ""
			dotS += "digraph G {\n"
			dotS += "node[shape=none]\n"
			dotS += "start[label=<<table>\n"
			dotS += "<tr><td colspan=\"2\" bgcolor=\"#147e0d\"><font color=\"white\">REPORTE DE SUPERBLOQUE</font></td></tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">sb_nombre_hd</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + nombreD + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_filesystem_type</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_filesystem_type)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_inodes_count</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_inodes_count)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_blocks_count</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_blocks_count)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_free_blocks_count</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_free_blocks_count)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_free_inodes_count</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_free_inodes_count)) + "</td>\n"
			dotS += "</tr>\n"

			tm := time.Unix(sb.S_mtime, 0)
			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_mtime</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + tm.Format("2006-01-02 15:04:05") + "</td>\n"
			dotS += "</tr>\n"

			tm = time.Unix(sb.S_umtime, 0)
			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_umtime</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + tm.Format("2006-01-02 15:04:05") + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_mnt_count</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_mnt_count)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_magic</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_magic)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_inode_s</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_inode_s)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_block_s</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_block_s)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_firts_ino</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_firts_ino)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_first_blo</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_first_blo)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_bm_inode_start</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_bm_inode_start)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_bm_block_start</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_bm_block_start)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">s_inode_start</td>\n"
			dotS += "<td bgcolor=\"#b4f0b1\">" + strconv.Itoa(int(sb.S_inode_start)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "<tr>\n"
			dotS += "<td bgcolor=\"#27ba40\">s_block_start</td>\n"
			dotS += "<td bgcolor=\"#27ba40\">" + strconv.Itoa(int(sb.S_block_start)) + "</td>\n"
			dotS += "</tr>\n"

			dotS += "</table>>];\n"
			dotS += "}"

			errD = os.WriteFile("Reportes/"+nombreD+".dot", []byte(dotS), 0777)
			if errD != nil {
				fmt.Println(errD)
			}

			ext := Extrep
			_, errD = exec.Command("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD + "." + ext).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}
			_, errD = exec.Command("dot", "-T"+ext, "Reportes/"+nombreD+".dot", "-o", Dirrep+nombreD + "." + ext).Output()
			if errD != nil {
				fmt.Printf("%s", errD)
			}

			return Structs.Resp{Res: "SE GENERO EL REPORTE SB", Dot2:dotS}
		}
		return Structs.Resp{Res: "DISCO INEXISTENTE"}
	}
	return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + Idrep}
}

func fileR() Structs.Resp {
	nodo := Mlist.buscar(Idrep)
	if nodo != nil {
		if Rutarep != " " {
			Dirrep = GetDirectorio(Prep)
			Extrep = GetExtension(Prep)
			nombreD := nombre(Prep)

			err := os.MkdirAll(Dirrep, 0777)
			if err != nil {
				fmt.Printf("%s", err)
			}
			file, errf := os.OpenFile(nodo.Path, os.O_RDWR, 0777)
			if errf == nil {
				sb := Structs.SuperBloque{}
				if nodo.Type == 'p' {
					mbr := Structs.MBR{}
					file.Seek(0, 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(mbr))), binary.BigEndian, &mbr)
					if mbr.Mbr_partition[nodo.Pos].Part_status != '2' {
						file.Close()
						return Structs.Resp{Res: "NO SE HA FORMATEADO LA MONTURA DE LA PARTICION " + nodo.Name}
					}
					file.Seek(int64(nodo.Start), 0)
				} else if nodo.Type == 'l' {
					ebr := Structs.EBR{}
					file.Seek(int64(nodo.Start), 0)
					errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(ebr))), binary.BigEndian, &ebr)
					if ebr.Part_status != '2' {
						file.Close()
						return Structs.Resp{Res: "NO SE HA FORMATEADO LA MONTURA DE LA PARTICION " + nodo.Name}
					}
					file.Seek(int64(nodo.Start+int(unsafe.Sizeof(Structs.EBR{}))), 0)
				}
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(sb))), binary.BigEndian, &sb)
				rutaS := strings.Split(Rutarep, "/")
				if len(rutaS) < 2 {
					file.Close()
					return Structs.Resp{Res: "RUTA INVALIDA"}
				}

				posInodoF := getInodoF(rutaS, 1, len(rutaS)-1, int(sb.S_inode_start), file)

				if posInodoF == -1 {
					return Structs.Resp{Res: "ARCHIVO NO ENCONTRADO"}
				}

				var inodo Structs.TablaInodo
				file.Seek(int64(posInodoF), 0)
				errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)

				if inodo.I_type != '1' {
					return Structs.Resp{Res: "LA RUTA NO HACE REFERENCIA A UN ARCHIVO"}
				}

				inodo.I_atime = time.Now().Unix()
				file.Seek(int64(posInodoF), 0)
				var bufferInode bytes.Buffer
				errf = binary.Write(&bufferInode, binary.BigEndian, inodo)
				EscribirFile(file, bufferInode.Bytes())

				dot, errD := os.OpenFile("Reportes/"+nombreD+".dot", os.O_CREATE, 0777)
				dot.Close()
				if errD != nil {
					fmt.Println(errD)
				}

				content := getContenArchivo(posInodoF, file)
				dotS := ""
				dotS += "digraph G {\n"
				dotS += "node[shape=none, lblstyle=\"align=left\"]\n"
				dotS += "start[label=\""

				dotS += rutaS[len(rutaS)-1] + "\n"
				dotS += content
				dotS += "\"]"
				dotS += "}"
				errD = os.WriteFile("Reportes/"+nombreD+".dot", []byte(dotS), 0777)
				if errD != nil {
					fmt.Println(errD)
				}

				ext := Extrep
				_, errD = exec.Command("dot", "-T"+Extrep, "Reportes/"+nombreD+".dot", "-o", "Reportes/"+nombreD).Output()
				if errD != nil {
					fmt.Printf("%s", errD)
				}
				_, errD = exec.Command("dot", "-T"+ext, "Reports/"+nombreD+".dot", "-o", Dirrep+nombreD).Output()
				if errD != nil {
					fmt.Printf("%s", errD)
				}

				file.Close()
				return Structs.Resp{Res: "SE GENERO EL REPORTE FILE DE " + rutaS[len(rutaS)-1], Dot4:dotS}
			}
			return Structs.Resp{Res: "DISCO INEXISTENTE"}
		}
		return Structs.Resp{Res: "DEBE ESCRIBIR RUTA DEL ARCHIVO"}
	}
	return Structs.Resp{Res: "NO SE HA ENCONTRADO ALGUNA MONTURA CON EL ID: " + Idrep}
}

func nombre(path string) string {
	directorio := ""
	aux := path
	p := strings.Index(aux, "/")
	for p != -1 {
		directorio += aux[:p] + "/"
		aux = aux[p+1:]
		p = strings.Index(aux, "/")
	}
	i := find(aux, ".")

	return aux[:i]
}

func getPartName(partName [16]byte) string {
	name := ""
	for i := 0; i < 16; i++ {
		if partName[i] == '\000' {
			break
		}
		name += string(partName[i])
	}
	return name
}

func getInodoF(rutaS []string, posAct int, rutaSize int, start int, file *os.File) int {
	var inodo Structs.TablaInodo
	var carpeta Structs.BloqueCarpeta

	file.Seek(int64(start), 0)
	errf := binary.Read(LeerFile(file, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)
	if errf != nil {
		fmt.Println(errf)
	}

	if inodo.I_type == '1' {
		return -1
	}

	for i := 0; i < len(inodo.I_block); i++ {
		if inodo.I_block[i] != -1 {
			file.Seek(int64(int(inodo.I_block[i])), 0)
			errf = binary.Read(LeerFile(file, int(unsafe.Sizeof(carpeta))), binary.BigEndian, &carpeta)
			for c := 0; c < 4; c++ {
				name := getContentName(carpeta.B_content[c].B_name)
				if name == rutaS[posAct] {
					if posAct < rutaSize {
						return getInodoF(rutaS, posAct+1, rutaSize, int(carpeta.B_content[c].B_inodo), file)
					}

					if posAct == rutaSize {
						return int(carpeta.B_content[c].B_inodo)
					}
				}
			}
		}
	}
	return -1
}

func getContentName(name [12]byte) string {
	var cadena string

	for i := 0; i < 12; i++ {
		if name[i] == '\000' {
			break
		}
		cadena += string([]byte{name[i]})
	}

	return cadena
}

func getContenArchivo(inodoStart int, file *os.File) string {
	var inodo Structs.TablaInodo
	var archivo Structs.BloqueArchivo
	file.Seek(int64(inodoStart), 0)
	errfU = binary.Read(LeerFile(file, int(unsafe.Sizeof(inodo))), binary.BigEndian, &inodo)
	content := ""
	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			file.Seek(int64(inodo.I_block[i]), 0)
			errfU = binary.Read(LeerFile(file, int(unsafe.Sizeof(archivo))), binary.BigEndian, &archivo)
			content += archivoContent2(archivo.B_content)
		}
	}
	return content
}
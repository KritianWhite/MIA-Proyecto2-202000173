package Files_System

import (
	"MIA-Proyecto2-202000173/Structs"
	"strconv"
)

type Nodo_M struct {
	Path  string
	Name  string
	Id    string
	Num   int
	Pos   int
	Type  byte
	Letra byte
	Start int
	Sig   *Nodo_M
}

type MountList struct {
	Primero *Nodo_M
	Ultimo  *Nodo_M
}

func (L *MountList) add(path string, name string, ty byte, start int, pos int) Structs.Resp {
	if !L.existMount(path, name) {
		num := L.getNum(path)
		letra := L.getLetra(path)

		nuevo := &Nodo_M{}
		nuevo.Path = path
		nuevo.Name = name
		nuevo.Type = ty
		nuevo.Num = num
		nuevo.Letra = letra
		nuevo.Start = start
		nuevo.Pos = pos
		nuevo.Id = "73" + strconv.Itoa(nuevo.Num) + string(nuevo.Letra)
		if L.Primero == nil {
			L.Primero = nuevo
			L.Ultimo = nuevo
		} else {
			L.Ultimo.Sig = nuevo
			L.Ultimo = nuevo
		}
		return Structs.Resp{Res: "SE MONTO LA PARTICION CON ID " + L.Ultimo.Id}
	}
	return Structs.Resp{Res: "LA PARTICION " + name + " YA ESTA MONTADA"}
}

func (L *MountList) existMount(path string, name string) bool {
	aux := L.Primero
	for aux != nil {
		if aux.Path == path && aux.Name == name {
			return true
		}
		aux = aux.Sig
	}
	return false
}

func (L *MountList) getNum(path string) int {
	mayor := 0
	aux := L.Primero
	for aux != nil {
		if aux.Path == path && aux.Num > mayor {
			mayor = aux.Num
		}
		aux = aux.Sig
	}
	return mayor + 1
}

func (L *MountList) getLetra(path string) byte {
	aux := L.Primero
	var letraMayor byte = 64
	for aux != nil {
		letraAct := aux.Letra
		if aux.Path == path {
			return aux.Letra
		}
		if letraAct > letraMayor {
			letraMayor = letraAct
		}
		aux = aux.Sig
	}
	return letraMayor + 1
}

func (L *MountList) buscar(id string) *Nodo_M {
	aux := L.Primero
	for aux != nil {
		if aux.Id == id {
			return aux
		}
		aux = aux.Sig
	}
	return aux
}

func (L *MountList) eliminar(id string) Structs.Bandera {
	if L.Primero != nil {
		if L.Primero == L.Ultimo && L.Primero.Id == id {
			L.Primero = nil
			L.Ultimo = nil
			return Structs.Bandera{Val: true}
		} else if L.Primero.Id == id {
			L.Primero = L.Primero.Sig
			return Structs.Bandera{Val: true}
		} else {
			aux := L.Primero.Sig
			ant := L.Primero
			for aux != nil {
				if aux.Id == id {
					ant.Sig = aux.Sig
					aux.Sig = nil
					return Structs.Bandera{Val: true}
				}
				ant = aux
				aux = aux.Sig
			}
			return Structs.Bandera{Val: false, Men: "NO EXISTE LA PARTICION " + id}
		}
	} else {
		return Structs.Bandera{Val: false, Men: "NO HAY PARTICIONES MONTADAS"}
	}
}

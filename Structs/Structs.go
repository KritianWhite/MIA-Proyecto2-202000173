package Structs

// Structs de administrador
type Resp struct {
	Res string  `json:"res"`
	U   Usuario `json:"usuario"`
	Dot1 string  `json:"dot1"`
	Dot2 string  `json:"dot2"`
	Dot3 string  `json:"dot3"`
	Dot4 string  `json:"dot4"`
}

type Bandera struct {
	Val bool
	Men string
}

type Inicio struct {
	Res string  `json:"res"`
	U   Usuario `json:"usuario"`
}

type Entrada struct {
	Command string `json:"comando"`
	IdU     int32  `json:"idU"`
	IdG     int32  `json:"idG"`
	IdMount string `json:"idMount"`
	NombreU string `json:"nombreU"`
	Login   bool   `json:"login"`
}

type Exec struct {
	Commands []string `json:"comandos"`
	IdU      int32    `json:"idU"`
	IdG      int32    `json:"idG"`
	IdMount  string   `json:"idMoun"`
	NombreU  string   `json:"nombreU"`
	Login    bool     `json:"login"`
	I        int      `json:"i"`
}

type Usuario struct {
	IdU     int32  `json:"id_u"`
	IdG     int32  `json:"id_g"`
	IdMount string `json:"id_mount"`
	NombreU string `json:"nombre_u"`
	Login   bool   `json:"login"`
}

// Structs del Sistema de Archivos
type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int32
	Part_s      int32
	Part_name   [16]byte
}

type MBR struct {
	Mbr_tamanio        int32
	Mbr_fecha_creacion int64
	Mbr_disk_signature int32
	Disk_fit           byte
	Mbr_partition      [4]Partition
}

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int32
	Part_s      int32
	Part_next   int32
	Part_name   [16]byte
}

type SuperBloque struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             int64
	S_umtime            int64
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_firts_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type TablaInodo struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime int64
	I_ctime int64
	I_mtime int64
	I_block [16]int32
	I_type  byte
	I_perm  int32
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type BloqueCarpeta struct {
	B_content [4]Content
}
type BloqueArchivo struct {
	B_content [64]byte
}

type Propiedad struct {
	Name string
	Val  string
}

type Comando struct {
	Name        string
	Propiedades []Propiedad
}
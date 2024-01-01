package main

import (
	"crypto/sha256"
	"fmt"
)

type Noeud struct {
	Type          int
	Data          []byte
	HashCalculate []byte
	HashReceive   []byte
	FilsGauche    *Noeud
	FilsDroit     *Noeud
}

/*
	func CreateNoeudChunk(data string, name string) *Noeud {
		return &Noeud{Type: 0, Name: name, Hash: StringToHash(data), FilsGauche: nil, FilsDroit: nil}
	}

	func CreateNoeudBigFile(data string, name string) *Noeud {
		if len(data) < 1024 {
			return &Noeud{Type: 0, Name: name, Hash: StringToHash(data), FilsGauche: nil, FilsDroit: nil}
		} else {
			milieu := len(data) / 2
			gauche := CreateNoeudBigFile(data[:milieu], name)
			droite := CreateNoeudBigFile(data[milieu:], name)
			return &Noeud{Type: 1, Name: name, Hash: StringToHash(string(gauche.Hash) + string(droite.Hash)), FilsGauche: gauche, FilsDroit: droite}
		}
	}

	func CreateDirectory(data string, name string) *Noeud {
		if len(data) < 1024 {
			return &Noeud{Type: 0, Name: name, Hash: StringToHash(data), FilsGauche: nil, FilsDroit: nil}
		} else {
			milieu := len(data) / 2
			gauche := CreateNoeudBigFile(data[:milieu], name)
			droite := CreateNoeudBigFile(data[milieu:], name)
			return &Noeud{Type: 2, Name: name, Hash: StringToHash(string(gauche.Hash) + string(droite.Hash)), FilsGauche: gauche, FilsDroit: droite}
		}
	}
*/
func StringToHash(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}
func Affiche(tab Noeud, profondeur int) {
	indent := ""
	for i := 0; i < profondeur; i++ {
		indent += "    "
	}
	fmt.Printf("%sHash: %s \n", indent, string(tab.Data))
	if tab.FilsGauche != nil {
		Affiche(*tab.FilsGauche, profondeur+1)
	}
	if tab.FilsDroit != nil {
		Affiche(*tab.FilsDroit, profondeur+1)
	}
}

/*
func RecalculateHashes(node *Noeud) {
	if node == nil {
		return
	}
	if node.FilsGauche != nil && node.FilsDroit != nil {
		RecalculateHashes(node.FilsGauche)
		RecalculateHashes(node.FilsDroit)
		node.Hash = StringToHash(node.FilsGauche.Hash + node.FilsDroit.Hash)
	}
}
*/
/*
func main() {
	root := CreateDirectory("data", "root")
	root.FilsGauche = CreateNoeudBigFile("bigfile_data", "bigfile")
	root.FilsDroit = CreateNoeudChunk("chunk_data", "chunk")

	root.FilsDroit.FilsGauche = CreateNoeudBigFile("directory_bigfile_data", "directory_bigfile")
	root.FilsDroit.FilsDroit = CreateNoeudChunk("directory_chunk_data", "directory_chunk")

	Affiche(*root, 0)
	RecalculateHashes(root)
	Affiche(*root, 0)
}
*/

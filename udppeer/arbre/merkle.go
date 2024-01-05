package arbre

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	. "projet-protocoles-internet/Tools"
	"sort"
)

type Noeud struct {
	//	HashCalculate []byte
	ID          int
	Type        int8
	HashReceive []byte
	NAME        string
	Data        []byte
	Fils        []*Noeud
}

func BuildImage(noeud *Noeud, chemin string) {
	if noeud == nil {
		return
	}

	cheminComplet := filepath.Join(chemin, noeud.NAME)

	fmt.Println(cheminComplet)

	switch noeud.Type {
	case ChunkType:
		if len(noeud.NAME) > 0 {
			err := os.WriteFile(cheminComplet, noeud.Data, 0644)
			if err != nil {
				fmt.Println("Error writing", err)
			}
		}
	case BigFileType:
		bytetab := make([]byte, 0)
		for _, child := range noeud.Fils {
			bytetab = append(bytetab, BuildBigFile(child)...)
		}
		//fmt.Println("tmp/peers/" + noeud.NAME)

		if len(bytetab) != 0 {
			fmt.Println(cheminComplet)
			err := os.WriteFile(cheminComplet, bytetab, 0644)
			if err != nil {
				fmt.Println("Error write file bigfile", err)
			}
		}

	case DirectoryType:
		err := os.MkdirAll(cheminComplet, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}

	for _, child := range noeud.Fils {
		BuildImage(child, cheminComplet)
	}
}

type ByID []*Noeud

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

func BuildBigFile(noeud *Noeud) []byte {
	bytetab := make([]byte, 0)

	sort.Sort(ByID(noeud.Fils))

	for _, child := range noeud.Fils {
		switch child.Type {
		case ChunkType:
			bytetab = append(bytetab, child.Data...)
		case BigFileType:
			bytetab = append(bytetab, BuildBigFile(child)...)
		}
	}
	return bytetab
}

func AfficherArbre(noeud *Noeud, niveau int) {

	if noeud == nil {
		return
	}
	indent := ""
	for i := 0; i < niveau; i++ {
		indent += "       -"
	}
	hashStr := hex.EncodeToString(noeud.HashReceive)
	fmt.Printf("%sNoeud : Type %d Fils: %d ID: %d Hash: %.5s, Name: %s, Data: %d\n", indent, noeud.Type, len(noeud.Fils), noeud.ID, hashStr, noeud.NAME, len(noeud.Data))

	for _, enfant := range noeud.Fils {
		AfficherArbre(enfant, niveau+1)
	}
}

func ChangeDataFromHash(root *Noeud, hashATrouver []byte, newData []byte) bool {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if CompareHashes(currentNode.HashReceive, hashATrouver) {
			copybyte := make([]byte, len(newData))
			copy(copybyte, newData)
			currentNode.Data = copybyte
			currentNode.Type = 0
			return true // Retourne vrai si les données ont été modifiées
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func SetType(root *Noeud, hashATrouver []byte, typeFile int8) bool {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if CompareHashes(currentNode.HashReceive, hashATrouver) {
			currentNode.Type = typeFile
			return true // Retourne vrai si les données ont été modifiées
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func AddNodeFromHash(root *Noeud, hash []byte, noeudToAdd *Noeud) {
	var queue []*Noeud
	queue = append(queue, root)
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		// Vérifie si le nœud actuel a le hash recherché
		//fmt.Println(hex.EncodeToString(currentNode.HashReceive), " avec ", hex.EncodeToString(hash))

		if CompareHashes(currentNode.HashReceive, hash) {
			currentNode.Fils = append(currentNode.Fils, noeudToAdd)
			currentNode.Type = 1
			//fmt.Println("Hash trouvé j'ajoute")
			return
		}

		// Ajoute les fils du nœud actuel à la file
		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	//fmt.Println("Aucun Hash j'ajoute au noeud racine")
	root.HashReceive = hash
	root.Fils = append(root.Fils, noeudToAdd)
}

// compareHashes compare deux slices de bytes (hashes).
func CompareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		//fmt.Println("lenght hash  ", len(hash1), " != ", len(hash2))
		//fmt.Println(hex.EncodeToString(hash1), " ", hex.EncodeToString(hash2))
		return false
	}
	for i, b := range hash1 {
		if b != hash2[i] {
			//fmt.Println("different")
			//fmt.Println(hex.EncodeToString(hash1), " ", hex.EncodeToString(hash2))
			return false
		}
	}
	return true
}

func ParcourirRepertoire(chemin string) (*Noeud, error) {
	fichierInfo, err := os.Stat(chemin)
	if err != nil {
		return nil, err
	}

	noeud := &Noeud{NAME: fichierInfo.Name()}

	if fichierInfo.IsDir() {
		noeud.Type = DirectoryType

		body := make([]byte, 0)
		body = append(body, DirectoryType)

		fichiers, err := os.ReadDir(chemin)
		if err != nil {
			return nil, err
		}

		for _, fi := range fichiers {
			fils, err := ParcourirRepertoire(filepath.Join(chemin, fi.Name()))
			if err != nil {
				return nil, err
			}
			noeud.Fils = append(noeud.Fils, fils)
		}

		for i := range noeud.Fils {
			body = append(body, []byte(noeud.Fils[i].NAME)...)
			body = append(body, noeud.Fils[i].HashReceive...)
		}
		bodyConvert := sha256.Sum256(body)
		noeud.HashReceive = bodyConvert[:]

	} else {

		if fichierInfo.Size() <= 1024 {

			body := make([]byte, 0)
			body = append(body, ChunkType)
			body = append(body, noeud.Data...)

			noeud.Type = ChunkType
			data, err := os.ReadFile(chemin)
			noeud.Data = data

			bodyConvert := sha256.Sum256(body)
			noeud.HashReceive = bodyConvert[:]

			if err != nil {
				fmt.Println("Chunck generation failed", err)
			}

		} else {

			noeud.Type = BigFileType

			body := make([]byte, 0)
			body = append(body, BigFileType)

			data, err := os.ReadFile(chemin)
			if err != nil {
				return nil, err
			}

			nbfilsChuck := len(data)/ChunkSize + 1

			if nbfilsChuck > 32 {

				var tabTempoNoeud []*Noeud

				nbfils := nbfilsChuck/32 + 1

				for i := 0; i < nbfils; i++ {
					noeudCreate := &Noeud{Type: BigFileType, Data: make([]byte, 0), Fils: make([]*Noeud, 0), ID: i}
					tabTempoNoeud = append(tabTempoNoeud, noeudCreate)
					noeud.Fils = append(noeud.Fils, noeudCreate)
				}

				poseCounter := 0
				position := 0

				for i := 0; i < len(data); i += ChunkSize {
					fin := i + ChunkSize
					if fin > len(data) {
						fin = len(data)
					}

					tabTempoNoeud[position].Fils = append(tabTempoNoeud[position].Fils, &Noeud{
						Type: ChunkType,
						Data: data[i:fin],
						ID:   i / ChunkSize,
					})
					if poseCounter == 32 {
						poseCounter = 0
						position = position + 1
					}
					poseCounter = poseCounter + 1
				}

			} else {

				for i := 0; i < len(data); i += ChunkSize {
					fin := i + ChunkSize
					if fin > len(data) {
						fin = len(data)
					}

					chunckBody := make([]byte, 0)
					chunckBody = append(chunckBody, ChunkType)
					chunckBody = append(chunckBody, data[i:fin]...)

					bodyConvert := sha256.Sum256(chunckBody)

					noeud.Fils = append(noeud.Fils, &Noeud{
						Type:        ChunkType,
						Data:        data[i:fin],
						HashReceive: bodyConvert[:],
						ID:          i / ChunkSize,
					})

					body = append(body, bodyConvert[:]...)

				}
			}
		}

	}
	return noeud, nil
}

/*
func HashDFS(noeud *Noeud) {

	if noeud.Type == ChunkType {
		hashCalculate := sha256.Sum256([]byte(noeud.Data))
		noeud.HashReceive = hashCalculate[:]
		return
	}

	hash := make([]byte, 0)
	for _, fils := range noeud.Fils {
		HashDFS(fils)
		hash = append(hash, fils.HashReceive...)
	}
	hashCalculate := sha256.Sum256(hash)
	copy(noeud.HashReceive, hashCalculate[:])
}

func chunckEmpty(noeud *Noeud) []byte {
	return make([]byte, 0)
}*/

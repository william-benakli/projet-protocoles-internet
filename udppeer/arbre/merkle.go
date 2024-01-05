package arbre

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// TODO verifier quand on recoit un helloreply si on est deja en communication avec juliuz avec l'historique
type Noeud struct {
	//	HashCalculate []byte
	Pos         int32
	Type        int8
	HashReceive []byte
	NAME        string
	Data        []byte
	Fils        []*Noeud
}

const (
	DirectoryType = 2
	BigFileType   = 1
	ChunkType     = 0
	ChunkSize     = 1024
)

func ParcoursRec(noeud *Noeud) error {

	switch noeud.Type {
	case 1: // bigfile
		var donnéesComplètes []byte
		for _, fils := range noeud.Fils {
			donnéesComplètes = append(donnéesComplètes, fils.Data...)
		}
		return os.WriteFile("tmp/peers/"+noeud.NAME, donnéesComplètes, 0644)
	case 2: // directory
		for _, fils := range noeud.Fils {
			if err := ParcoursRec(fils); err != nil {
				return err
			}
		}
	default: // chunk
		return os.WriteFile("tmp/peers/"+noeud.NAME, noeud.Data, 0644)
	}

	return nil
}

func BuildImage(root *Noeud) {

	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		//fmt.Println(removeEmpty(string(currentNode.Data)))

		if currentNode.Type == 0 {
			if len(currentNode.NAME) > 0 {
				err := os.WriteFile("tmp/peers/"+currentNode.NAME, currentNode.Data, 0644)
				if err != nil {
					fmt.Println("Erreur writing")
				}
			}
		} else if currentNode.Type == 1 {
			bytetab := make([]byte, 0)

			for i := 0; i < len(currentNode.Fils); i++ {
				for j := 0; j < len(currentNode.Fils[i].Fils); j++ {
					for k := 0; k < len(currentNode.Fils[i].Fils[j].Data); k++ {
						bytetab = append(bytetab, currentNode.Fils[i].Fils[j].Data[k])
					}
				}
			}
			os.WriteFile("tmp/peers/"+currentNode.NAME, bytetab, 0644)

		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

}

func AfficherArbre(noeud *Noeud, niveau int) {

	if noeud == nil {
		return
	}
	/*	if noeud.Type == 0 {
			return
		}
	*/
	indent := ""
	for i := 0; i < niveau; i++ {
		indent += "       -"
	}

	hashStr := hex.EncodeToString(noeud.HashReceive)
	//dataStr := string(noeud.Data)
	/*if len(noeud.Data) == 0 && noeud.Type == 0 {
		return
	}*/
	fmt.Printf("%sNoeud : Type %d Fils: %d Hash: %.5s, Name: %s, Data: %d\n", indent, noeud.Type, len(noeud.Fils), hashStr, noeud.NAME, len(noeud.Data))

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
	} else {
		noeud.Type = BigFileType
		data, err := os.ReadFile(chemin)
		if err != nil {
			return nil, err
		}

		nbfilsChuck := len(data)/ChunkSize + 1

		if nbfilsChuck > 32 {

			var tabTempoNoeud []*Noeud

			nbfils := nbfilsChuck/32 + 1

			for i := 0; i < nbfils; i++ {
				noeudCreate := &Noeud{Type: BigFileType, Data: make([]byte, 0), Fils: make([]*Noeud, 0)}
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
				noeud.Fils = append(noeud.Fils, &Noeud{
					Type: ChunkType,
					Data: data[i:fin],
				})
			}
		}
	}
	return noeud, nil
}

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
}

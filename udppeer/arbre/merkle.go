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

var Root Noeud

func GetRoot() *Noeud {
	return &Root
}

func ResetRoot() {
	Root = Noeud{}
}

type Noeud struct {
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

	switch noeud.Type {
	case DirectoryType:

		err := os.MkdirAll(cheminComplet, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}

	case ChunkType:
		if len(noeud.NAME) > 0 {
			err := os.WriteFile(cheminComplet, noeud.Data, 0766)
			if err != nil {
				fmt.Println("Error writing", err)
				err = os.MkdirAll(cheminComplet, os.ModePerm)
				err = os.WriteFile(cheminComplet, noeud.Data, 0766)
			}
		}
	case BigFileType:
		bytetab := make([]byte, 0)
		for _, child := range noeud.Fils {
			bytetab = append(bytetab, BuildBigFile(child)...)
		}

		if len(bytetab) != 0 && len(noeud.NAME) > 0 {
			err := os.WriteFile(cheminComplet, bytetab, 0766)
			if err != nil {
				fmt.Println("Error write file bigfile", err)
			}
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

func ChangeDataFromHashRec(noeud *Noeud, hashATrouver []byte, newData []byte) {
	if CompareHashes(noeud.HashReceive, hashATrouver) {
		copybyte := make([]byte, len(newData))
		copy(copybyte, newData)
		noeud.Data = copybyte
		noeud.Type = 0
		fmt.Println("trouvé")
		return
	}

	for _, child := range noeud.Fils {
		ChangeDataFromHashRec(child, hashATrouver, newData)
	}
	//return 0 // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func SetTypeRec(noeud *Noeud, hashATrouver []byte, typeFile int8) bool {
	if CompareHashes(noeud.HashReceive, hashATrouver) {
		noeud.Type = typeFile
		return true
	}

	for _, child := range noeud.Fils {
		SetTypeRec(child, hashATrouver, typeFile)
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func AddNodeFromHashRec(noeud *Noeud, hash []byte, noeudToAdd *Noeud) {

	if len(Root.HashReceive) == 0 {
		Root.HashReceive = hash
		Root.Type = 2
		noeud.Fils = append(noeud.Fils, noeudToAdd)
		return
	}

	if CompareHashes(noeud.HashReceive, hash) {
		noeud.Fils = append(noeud.Fils, noeudToAdd)
		return
	}

	for _, child := range noeud.Fils {
		AddNodeFromHashRec(child, hash, noeudToAdd)
	}

	//fmt.Println("Pas trouvé mon papa")

}

// compareHashes compare deux slices de bytes (hashes).
func CompareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i, b := range hash1 {
		if b != hash2[i] {
			return false
		}
	}
	return true
}
func ParcourirRepertoire2(chemin string) (*Noeud, error) {
	fichierInfo, err := os.Stat(chemin)
	if err != nil {
		return nil, err
	}

	noeud := &Noeud{NAME: fichierInfo.Name()}

	if fichierInfo.IsDir() {

		body := make([]byte, 0)
		body = append(body, DirectoryType)

		noeud.Type = DirectoryType
		fichiers, err := os.ReadDir(chemin)
		if err != nil {
			return nil, err
		}

		for _, fi := range fichiers {
			fils, err := ParcourirRepertoire2(filepath.Join(chemin, fi.Name()))
			if err != nil {
				return nil, err
			}
			noeud.Fils = append(noeud.Fils, fils)
		}

		for i := range noeud.Fils {
			var byteArray [32]byte
			copy(byteArray[:], noeud.Fils[i].NAME)
			body = append(body, byteArray[:]...)
			body = append(body, noeud.Fils[i].HashReceive...)
		}
		bodyConvert := sha256.Sum256(body)
		noeud.HashReceive = bodyConvert[:]

	} else {
		if fichierInfo.Size() <= 1024 {
			//CHUNCK
			body := make([]byte, 0)
			body = append(body, ChunkType)

			noeud.Type = ChunkType
			data, err := os.ReadFile(chemin)

			noeud.Data = data
			body = append(body, noeud.Data...)

			bodyConvert := sha256.Sum256(body)
			noeud.HashReceive = bodyConvert[:]

			if err != nil {
				fmt.Println("Chunck generation failed", err)
			}

		} else {
			noeud.Type = BigFileType
			data, err := os.ReadFile(chemin)
			if err != nil {
				return nil, err
			}

			nbfilsChuck := len(data)/ChunkSize + 1

			if nbfilsChuck > 32 {

				body := make([]byte, 0)
				body = append(body, BigFileType)

				var tabTempoNoeud []*Noeud

				nbfils := nbfilsChuck/32 + 1

				for i := 0; i < nbfils; i++ {
					noeudCreate := &Noeud{Type: BigFileType, Data: make([]byte, 0), Fils: make([]*Noeud, 0), ID: i}
					tabTempoNoeud = append(tabTempoNoeud, noeudCreate)
					noeud.Fils = append(noeud.Fils, noeudCreate)
				}

				poseCounter := 0
				position := 0

				chunckSaveForDad := make([]byte, 0)
				chunckSaveForDad = append(chunckSaveForDad, BigFileType)
				for i := 0; i < len(data); i += ChunkSize {
					fin := i + ChunkSize
					if fin > len(data) {
						fin = len(data)
					}

					chunck := sha256.Sum256(data[i:fin])
					fmt.Println(position, "et ", len(tabTempoNoeud))
					tabTempoNoeud[position].Fils = append(tabTempoNoeud[position].Fils, &Noeud{
						Type:        ChunkType,
						Data:        data[i:fin],
						HashReceive: chunck[:],
						ID:          i / ChunkSize,
					})

					chunckSaveForDad = append(chunckSaveForDad, chunck[:]...)

					if poseCounter == 32 {
						bigfile := sha256.Sum256(chunckSaveForDad)
						tabTempoNoeud[position].HashReceive = bigfile[:]

						chunckSaveForDad = make([]byte, 0)
						chunckSaveForDad = append(chunckSaveForDad, BigFileType)

						poseCounter = 0
						position = position + 1
					}

					if position == nbfils-1 {
						bigfile := sha256.Sum256(chunckSaveForDad)
						tabTempoNoeud[position].HashReceive = bigfile[:]
					}

					poseCounter = poseCounter + 1
				}

				bigFileLast := make([]byte, 0)
				bigFileLast = append(bigFileLast, BigFileType)

				for i := 0; i < len(tabTempoNoeud); i += 1 {
					hash := tabTempoNoeud[i].HashReceive
					bigFileLast = append(bigFileLast, hash...)
				}

				bodyConvertBigFile := sha256.Sum256(bigFileLast)
				noeud.HashReceive = bodyConvertBigFile[:]

			} else {
				body := make([]byte, 0)
				body = append(body, BigFileType)

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
				bodyConvert := sha256.Sum256(body)
				noeud.HashReceive = bodyConvert[:]
			}
		}
	}
	return noeud, nil
}

func GetHashDFS(n *Noeud, hash []byte) *Noeud {
	if CompareHashes(n.HashReceive, hash) {
		fmt.Printf("Le nœud avec ID %d a un HashReceive correspondant.\n", n.ID)
		return n
	}
	for _, fils := range n.Fils {
		result := GetHashDFS(fils, hash)
		if result != nil {
			return result
		}
	}

	return nil
}

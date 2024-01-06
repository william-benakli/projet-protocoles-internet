package udppeer

import (
	"fmt"
	"projet-protocoles-internet/udppeer/arbre"
)

var racine arbre.Noeud

func InitRoot() {
	racine.Type = 2
	racine, err := arbre.ParcourirRepertoire("tmp/user")
	if err != nil {
		fmt.Println("Creation de l'arbre impossible")
	}
	arbre.AfficherArbre(racine, 0)
}

func GetRacine() *arbre.Noeud {
	return &racine
}

func getNoeudFromHash(hash []byte) *arbre.Noeud {
	var queue []*arbre.Noeud

	queue = append(queue, &racine)
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		// Vérifie si le nœud actuel a le hash recherché
		//fmt.Println(hex.EncodeToString(currentNode.HashReceive), " avec ", hex.EncodeToString(hash))

		if arbre.CompareHashes(currentNode.HashReceive, hash) {
			return currentNode
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}
	return nil
}

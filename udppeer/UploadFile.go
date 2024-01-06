package udppeer

import (
	"projet-protocoles-internet/udppeer/arbre"
)

var racine *arbre.Noeud

func InitRoot() {
	racine = &arbre.Noeud{}
	racine.Type = 2
	racine, _ = arbre.ParcourirRepertoire("tmp/user")
	arbre.AfficherArbre(racine, 0)
}

func GetRacine() *arbre.Noeud {
	return racine
}

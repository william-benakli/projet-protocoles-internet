package udppeer

import (
	"projet-protocoles-internet/udppeer/arbre"
)

var Racine *arbre.Noeud

func InitRoot() {
	Racine = &arbre.Noeud{}
	Racine.Type = 2
	Racine, _ = arbre.ParcourirRepertoire2("tmp/user")
	arbre.AfficherArbre(Racine, 0)
}

func GetRacine() *arbre.Noeud {
	return Racine
}

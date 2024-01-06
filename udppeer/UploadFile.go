package udppeer

import (
	"projet-protocoles-internet/udppeer/arbre"
)

var Racine *arbre.Noeud

func InitRoot() {
	Racine = &arbre.Noeud{}
	Racine.Type = 2
	Racine, _ = arbre.ParcourirRepertoire("tmp/user")
	arbre.AfficherArbre(Racine, 0)
	//arbre.BuildImage(Racine, "tmp/peers/")

}

func GetRacine() *arbre.Noeud {
	return Racine
}

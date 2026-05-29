package main

import "fmt"

const (
	Nom          = "Léo"
	IMCMaigreur  = 18.5
	IMCNormal    = 25.0
	IMCSurpoids  = 30.0
)

func main() {
	var poids float64 = 70.5
	var taille float64 = 1.75

	imc := poids / (taille * taille)

	fmt.Printf("Bonjour %s !\n", Nom)
	fmt.Printf("Poids : %.1f kg | Taille : %.2f m\n", poids, taille)
	fmt.Printf("IMC : %.2f\n", imc)

	var categorie string
	if imc < IMCMaigreur {
		categorie = "Maigreur"
	} else if imc < IMCNormal {
		categorie = "Normal"
	} else if imc < IMCSurpoids {
		categorie = "Surpoids"
	} else {
		categorie = "Obésité"
	}

	fmt.Printf("Catégorie : %s\n", categorie)
}

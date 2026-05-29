package main

import (
	"errors"
	"fmt"
	"os"
)

func operer(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("erreur : division par zéro")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("erreur : opération inconnue '%s'", op)
	}
}

func creerOperation(op string) func(float64, float64) float64 {
	return func(a, b float64) float64 {
		resultat, err := operer(a, b, op)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return resultat
	}
}

func main() {
	fmt.Println("Calculatrice — entrez : <nombre> <nombre> <opération> (ou 'quit' pour quitter)")

	var a, b float64
	var op string

	for {
		fmt.Print("> ")
		_, err := fmt.Scan(&a, &b, &op)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erreur de lecture, fin du programme")
			break
		}

		if op == "quit" {
			fmt.Println("Au revoir !")
			break
		}

		resultat, err := operer(a, b, op)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("= %g\n", resultat)
		}
	}
}

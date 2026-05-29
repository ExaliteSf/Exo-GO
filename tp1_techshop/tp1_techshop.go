package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Produit struct {
	ID        int
	Nom       string
	Marque    string
	Prix      float64
	Stock     int
	Categorie string
	Actif     bool
}

type Catalogue struct {
	produits []Produit
}

func (c *Catalogue) AjouterProduit(p Produit) error {
	for _, existing := range c.produits {
		if existing.ID == p.ID {
			return fmt.Errorf("ID %d déjà utilisé par \"%s\"", p.ID, existing.Nom)
		}
	}
	c.produits = append(c.produits, p)
	return nil
}

func (c *Catalogue) TrouverParID(id int) (Produit, error) {
	for _, p := range c.produits {
		if p.ID == id {
			return p, nil
		}
	}
	return Produit{}, fmt.Errorf("aucun produit avec l'ID %d", id)
}

func (c *Catalogue) TrouverParCategorie(cat string) []Produit {
	var result []Produit
	for _, p := range c.produits {
		if strings.EqualFold(p.Categorie, cat) {
			result = append(result, p)
		}
	}
	return result
}

func (c *Catalogue) AppliquerReduction(categorie string, pct float64) int {
	count := 0
	for i, p := range c.produits {
		if strings.EqualFold(p.Categorie, categorie) {
			c.produits[i].Prix = p.Prix * (1 - pct/100)
			count++
		}
	}
	return count
}

func (c *Catalogue) Vendre(id int, qte int) error {
	if qte <= 0 {
		return errors.New("la quantité doit être positive")
	}
	for i, p := range c.produits {
		if p.ID == id {
			if p.Stock < qte {
				return fmt.Errorf("stock insuffisant : %d disponible(s), %d demandé(s)", p.Stock, qte)
			}
			c.produits[i].Stock -= qte
			return nil
		}
	}
	return fmt.Errorf("aucun produit avec l'ID %d", id)
}

func (c *Catalogue) Rapport() string {
	valeurTotale := 0.0
	for _, p := range c.produits {
		valeurTotale += p.Prix * float64(p.Stock)
	}
	return fmt.Sprintf("Nb produits : %d | Valeur totale du stock : %.2f €", len(c.produits), valeurTotale)
}

var scanner = bufio.NewScanner(os.Stdin)

func lireLigne(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func lireInt(prompt string) int {
	var n int
	fmt.Print(prompt)
	fmt.Scan(&n)
	scanner.Scan() // vide le buffer
	return n
}

func lireFloat(prompt string) float64 {
	var f float64
	fmt.Print(prompt)
	fmt.Scan(&f)
	scanner.Scan()
	return f
}

func afficherProduit(p Produit) {
	fmt.Printf("  [%d] %-25s %-12s %8.2f € | Stock : %d | Catégorie : %s\n",
		p.ID, p.Nom, p.Marque, p.Prix, p.Stock, p.Categorie)
}

func main() {
	cat := &Catalogue{}

	initProduits := []Produit{
		{1, "iPhone 15 Pro", "Apple", 1199.99, 10, "Smartphone", true},
		{2, "MacBook Air M3", "Apple", 1499.99, 5, "Ordinateur", true},
		{3, "Galaxy S24 Ultra", "Samsung", 1099.99, 8, "Smartphone", true},
		{4, "iPad Pro 12.9", "Apple", 1299.99, 6, "Tablette", true},
		{5, "Surface Pro 9", "Microsoft", 999.99, 4, "Ordinateur", true},
	}
	for _, p := range initProduits {
		cat.AjouterProduit(p)
	}

	for {
		fmt.Println("\n========== TechShop CLI ==========")
		fmt.Println("[1] Ajouter  [2] Chercher  [3] Soldes  [4] Vendre  [5] Rapport  [0] Quitter")
		fmt.Print("> ")

		var choix int
		fmt.Scan(&choix)
		scanner.Scan()

		switch choix {
		case 0:
			fmt.Println("Au revoir !")
			return

		case 1:
			fmt.Println("--- Ajouter un produit ---")
			id := lireInt("ID         : ")
			nom := lireLigne("Nom        : ")
			marque := lireLigne("Marque     : ")
			prix := lireFloat("Prix (€)   : ")
			stock := lireInt("Stock      : ")
			categorie := lireLigne("Catégorie  : ")

			p := Produit{ID: id, Nom: nom, Marque: marque, Prix: prix, Stock: stock, Categorie: categorie, Actif: true}
			if err := cat.AjouterProduit(p); err != nil {
				fmt.Println("Erreur :", err)
			} else {
				fmt.Printf("Produit \"%s\" ajouté avec succès.\n", nom)
			}

		case 2:
			fmt.Println("--- Chercher par catégorie ---")
			categorie := lireLigne("Catégorie : ")
			produits := cat.TrouverParCategorie(categorie)
			if len(produits) == 0 {
				fmt.Printf("Aucun produit trouvé dans la catégorie \"%s\".\n", categorie)
			} else {
				fmt.Printf("%d produit(s) trouvé(s) :\n", len(produits))
				for _, p := range produits {
					afficherProduit(p)
				}
			}

		case 3:
			fmt.Println("--- Appliquer des soldes ---")
			categorie := lireLigne("Catégorie    : ")
			pct := lireFloat("Réduction (%) : ")
			n := cat.AppliquerReduction(categorie, pct)
			if n == 0 {
				fmt.Printf("Aucun produit dans la catégorie \"%s\".\n", categorie)
			} else {
				fmt.Printf("%.0f%% appliqués sur %d produit(s) de la catégorie \"%s\".\n", pct, n, categorie)
			}

		case 4:
			fmt.Println("--- Vendre un produit ---")
			id := lireInt("ID du produit : ")
			p, err := cat.TrouverParID(id)
			if err != nil {
				fmt.Println("Erreur :", err)
				break
			}
			fmt.Printf("Produit : %s %s | Stock : %d\n", p.Marque, p.Nom, p.Stock)
			qte := lireInt("Quantité      : ")
			if err := cat.Vendre(id, qte); err != nil {
				fmt.Println("Erreur :", err)
			} else {
				fmt.Printf("Vente de %d x \"%s\" effectuée.\n", qte, p.Nom)
			}

		case 5:
			fmt.Println("--- Rapport du catalogue ---")
			fmt.Println(cat.Rapport())
			fmt.Println("Détail :")
			for _, p := range cat.produits {
				afficherProduit(p)
			}

		default:
			fmt.Println("Choix invalide, entrez un nombre entre 0 et 5.")
		}
	}
}

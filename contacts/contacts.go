package main

import "fmt"

type Personne struct {
	Prenom string
	Nom    string
	Age    int
	Email  string
}

func (p Personne) NomComplet() string {
	return fmt.Sprintf("%s %s", p.Prenom, p.Nom)
}

func (p Personne) Presentation() string {
	return fmt.Sprintf("Nom : %s | Age : %d ans | Email : %s", p.NomComplet(), p.Age, p.Email)
}

type Adresse struct {
	Rue        string
	Ville      string
	CodePostal string
}

func (a Adresse) Format() string {
	return fmt.Sprintf("%s, %s %s", a.Rue, a.CodePostal, a.Ville)
}

type Employe struct {
	Personne
	Adresse
	Poste   string
	Salaire float64
}

func (e Employe) FicheEmploye() string {
	return fmt.Sprintf(
		"--- Fiche Employé ---\n%s\nPoste   : %s\nSalaire : %.2f €\nAdresse : %s",
		e.Presentation(), e.Poste, e.Salaire, e.Format(),
	)
}

func (e *Employe) AugmenterSalaire(pct float64) {
	e.Salaire += e.Salaire * pct / 100
}

type Etudiant struct {
	Personne
	Promo   string
	Moyenne float64
}

func (e Etudiant) MentionObtenue() string {
	switch {
	case e.Moyenne >= 16:
		return "Très Bien"
	case e.Moyenne >= 14:
		return "Bien"
	case e.Moyenne >= 12:
		return "Assez Bien"
	default:
		return "Passable"
	}
}

func (e Etudiant) FicheEtudiant() string {
	return fmt.Sprintf(
		"--- Fiche Étudiant ---\n%s\nPromo   : %s\nMoyenne : %.2f | Mention : %s",
		e.Presentation(), e.Promo, e.Moyenne, e.MentionObtenue(),
	)
}

func main() {
	e1 := Employe{
		Personne: Personne{"Alice", "Martin", 34, "alice.martin@corp.fr"},
		Adresse:  Adresse{"12 rue de la Paix", "Paris", "75001"},
		Poste:    "Développeuse Go",
		Salaire:  3800.0,
	}
	e2 := Employe{
		Personne: Personne{"Thomas", "Bernard", 41, "t.bernard@corp.fr"},
		Adresse:  Adresse{"8 avenue Foch", "Lyon", "69006"},
		Poste:    "Chef de projet",
		Salaire:  4500.0,
	}

	e1.AugmenterSalaire(10)
	e2.AugmenterSalaire(5)

	s1 := Etudiant{
		Personne: Personne{"Emma", "Dupont", 20, "emma.dupont@etu.fr"},
		Promo:    "BUT Info 2025",
		Moyenne:  15.8,
	}
	s2 := Etudiant{
		Personne: Personne{"Hugo", "Leroy", 22, "hugo.leroy@etu.fr"},
		Promo:    "Master IA 2025",
		Moyenne:  11.4,
	}

	fmt.Println(e1.FicheEmploye())
	fmt.Println()
	fmt.Println(e2.FicheEmploye())
	fmt.Println()
	fmt.Println(s1.FicheEtudiant())
	fmt.Println()
	fmt.Println(s2.FicheEtudiant())
}

package main

import (
	"errors"
	"fmt"
	"math"
)

// ══ Interface Payeur ══
type Payeur interface {
	Payer(montant float64) (string, error)
}

// ══ CarteCredit ══
type CarteCredit struct {
	Numero    string
	Titulaire string
	Solde     float64
}

func (cc *CarteCredit) Payer(montant float64) (string, error) {
	if montant > cc.Solde {
		return "", &ErrPaiement{
			Code: 402,
			Msg:  fmt.Sprintf("solde insuffisant : %.2f€ disponible, %.2f€ demandé", cc.Solde, montant),
		}
	}
	cc.Solde -= montant
	// Derniers 4 chiffres du numéro de carte
	last4 := cc.Numero
	if len(cc.Numero) >= 4 {
		last4 = cc.Numero[len(cc.Numero)-4:]
	}
	return fmt.Sprintf("Transaction CB #%s confirmée", last4), nil
}

// ══ PayPal ══
type PayPal struct {
	Email string
	Solde float64
}

func (pp *PayPal) Payer(montant float64) (string, error) {
	if montant > pp.Solde {
		return "", fmt.Errorf("solde PayPal insuffisant : %.2f€ disponible, %.2f€ demandé", pp.Solde, montant)
	}
	pp.Solde -= montant
	return fmt.Sprintf("Paiement PayPal de %.2f€ vers %s", montant, pp.Email), nil
}

// ══ Crypto ══
type Crypto struct {
	Adresse string
	Solde   float64 // en unités crypto (ex: BTC)
	Monnaie string
}

func (c *Crypto) Payer(montant float64) (string, error) {
	// 1 BTC = 50 000€
	const tauxBTC = 50000.0
	montantCrypto := math.Round(montant/tauxBTC*1000) / 1000
	if montantCrypto > c.Solde {
		return "", fmt.Errorf("solde %s insuffisant : %.4f disponible, %.4f requis",
			c.Monnaie, c.Solde, montantCrypto)
	}
	c.Solde -= montantCrypto
	return fmt.Sprintf("Paiement Crypto de %.2f€ = %.4f %s vers %s",
		montant, montantCrypto, c.Monnaie, c.Adresse), nil
}

// ══ Vérifications statiques à la compilation ══
var _ Payeur = &CarteCredit{}
var _ Payeur = &PayPal{}
var _ Payeur = &Crypto{}

// ══ ProcesserPanier ══
func ProcesserPanier(payeur Payeur, articles []float64) {
	// Calcul du total
	var total float64
	for _, prix := range articles {
		total += prix
	}
	fmt.Printf("🛒 Panier : %d article(s) — Total : %.2f€\n", len(articles), total)

	// Identification du mode de paiement via type switch
	switch p := payeur.(type) {
	case *CarteCredit:
		fmt.Printf("💳 Mode : Carte de crédit (%s — titulaire : %s)\n", p.Numero, p.Titulaire)
	case *PayPal:
		fmt.Printf("🅿️  Mode : PayPal (%s)\n", p.Email)
	case *Crypto:
		fmt.Printf("₿  Mode : Crypto (%s — solde : %.4f %s)\n", p.Adresse, p.Solde, p.Monnaie)
	default:
		fmt.Println("❓ Mode de paiement inconnu")
	}

	// Paiement
	msg, err := payeur.Payer(total)
	if err != nil {
		fmt.Printf("❌ Échec du paiement : %v\n", err)
		return
	}
	fmt.Printf("✅ %s\n", msg)
}

func main() {
	articles := []float64{29.99, 15.00, 4.49}

	fmt.Println("══════════════════════════════════")
	fmt.Println("       SYSTÈME DE PAIEMENT        ")
	fmt.Println("══════════════════════════════════")

	// ── Paiement par carte crédit (succès)
	fmt.Println("\n── Carte de crédit ──")
	carte := &CarteCredit{
		Numero:    "4532015112830366",
		Titulaire: "Alice Martin",
		Solde:     200.00,
	}
	ProcesserPanier(carte, articles)
	fmt.Printf("   Solde restant : %.2f€\n", carte.Solde)

	// ── Paiement par carte crédit (échec — solde insuffisant)
	fmt.Println("\n── Carte de crédit (solde insuffisant) ──")
	cartePauvre := &CarteCredit{
		Numero:    "4111111111111111",
		Titulaire: "Bob Faucheux",
		Solde:     10.00,
	}
	ProcesserPanier(cartePauvre, articles)

	// ── Paiement PayPal (succès)
	fmt.Println("\n── PayPal ──")
	paypal := &PayPal{
		Email: "alice@gmail.com",
		Solde: 100.00,
	}
	ProcesserPanier(paypal, articles)
	fmt.Printf("   Solde restant : %.2f€\n", paypal.Solde)

	// ── Paiement Crypto (succès)
	fmt.Println("\n── Bitcoin ──")
	bitcoin := &Crypto{
		Adresse: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		Solde:   0.01,
		Monnaie: "BTC",
	}
	ProcesserPanier(bitcoin, articles)
	fmt.Printf("   Solde restant : %.6f BTC\n", bitcoin.Solde)

	// ── Paiement Crypto (échec — solde insuffisant)
	fmt.Println("\n── Bitcoin (solde insuffisant) ──")
	bitcoinPauvre := &Crypto{
		Adresse: "bc1qfauxadresse",
		Solde:   0.000001,
		Monnaie: "BTC",
	}
	ProcesserPanier(bitcoinPauvre, []float64{999.99, 299.00})

	// ── Démonstration du piège interface nil ──
	fmt.Println("\n── Piège nil interface ──")
	fmt.Printf("fonctionDangereuse() == nil → %v\n", fonctionDangereuse() == nil)
	fmt.Printf("fonctionCorrecte()   == nil → %v\n", fonctionCorrecte() == nil)

	// ── Démonstration any + type assertion ──
	fmt.Println("\n── Type assertion ──")
	var val any = carte
	if cc, ok := val.(*CarteCredit); ok {
		fmt.Printf("Assertion réussie : titulaire = %s\n", cc.Titulaire)
	}

	// ── Erreur personnalisée ──
	fmt.Println("\n── Erreur personnalisée ──")
	_, err := cartePauvre.Payer(9999)
	var errPaiement *ErrPaiement
	if errors.As(err, &errPaiement) {
		fmt.Printf("ErrPaiement capturée : code=%d msg=%s\n", errPaiement.Code, errPaiement.Msg)
	}
}

// ══ Bonus : type d'erreur custom satisfaisant error ══
type ErrPaiement struct {
	Code int
	Msg  string
}

func (e *ErrPaiement) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

// ══ Démonstration du piège interface nil ══
type monErreurInterne struct{ msg string }

func (e *monErreurInterne) Error() string { return e.msg }

func fonctionDangereuse() error {
	var err *monErreurInterne = nil
	return err // ← interface non-nil ! (type concret = *monErreurInterne)
}

func fonctionCorrecte() error {
	return nil // ← interface réellement nil
}

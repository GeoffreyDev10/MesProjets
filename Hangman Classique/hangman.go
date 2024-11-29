package HangmanWeb

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func LireMotsDepuisFichier(fichier string) ([]string, error) {
	file, err := os.Open(fichier)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mots []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mot := strings.TrimSpace(scanner.Text())
		if len(mot) > 0 {
			mots = append(mots, mot)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mots, nil
}

func MotAleatoire(mots []string) string {
	rand.Seed(time.Now().UnixNano())
	indexAleatoire := rand.Intn(len(mots))
	return mots[indexAleatoire]
}

func CamouflerMot(mot string, lettresAReveler int) string {
	motCamoufle := strings.Repeat("_", len(mot))

	if lettresAReveler > len(mot) {
		lettresAReveler = len(mot)
	}

	revelees := make(map[int]struct{})
	lettresRevelees := make(map[rune]struct{})

	for i := 0; i < lettresAReveler; i++ {
		var index int
		for {
			index = rand.Intn(len(mot))
			if _, exists := revelees[index]; !exists {
				revelees[index] = struct{}{}
				lettresRevelees[rune(mot[index])] = struct{}{}
				motCamoufle = motCamoufle[:index] + string(mot[index]) + motCamoufle[index+1:]
				break
			}
		}
	}
	for lettre := range lettresRevelees {
		for i, l := range mot {
			if l == lettre {
				motCamoufle = motCamoufle[:i] + string(l) + motCamoufle[i+1:]
			}
		}
	}
	return motCamoufle
}

func RevelerLettres(mot string, lettresTrouvees []rune, motCamoufle string) string {
	for i, l := range motCamoufle {
		if l != '_' {
			lettresTrouvees[i] = rune(mot[i])
		}
	}
	return string(lettresTrouvees)
}

func RevelerLettresSimilaires(mot string, lettre string, lettresTrouvees []rune) {
	for i, l := range mot {
		if string(l) == lettre {
			lettresTrouvees[i] = l
		}
	}
}

func HangmanPositions(fichier string) ([]string, error) {
	file, err := os.Open(fichier)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hangs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hang := scanner.Text()
		hangs = append(hangs, hang)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hangs, nil
}

func AfficherHangman(hangs []string, essais int, lignesParEssai int) {
	debut := (essais - 1) * lignesParEssai
	fin := debut + lignesParEssai
	if fin > len(hangs) {
		fin = len(hangs)
	}
	for i := debut; i < fin; i++ {
		fmt.Println(hangs[i])
	}
}

func main() {
	mots, err := LireMotsDepuisFichier("words.txt")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	hangs, err := HangmanPositions("hangman.txt")
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	mot := MotAleatoire(mots)
	nLettresAReveler := len(mot)/2 - 1

	motCamoufle := CamouflerMot(mot, nLettresAReveler)

	lettresTrouvees := []rune(strings.Repeat("_", len(mot)))
	lettresTrouvees = []rune(RevelerLettres(mot, lettresTrouvees, motCamoufle))

	maxEssais := 10
	essais := 0
	lignesParEssai := 8

	fmt.Println("Bonne chance, tu as ", maxEssais-essais, "chance")
	fmt.Println(motCamoufle)

	for essais < maxEssais {
		fmt.Print("Entrez une lettre : ")
		var lettre string
		fmt.Scanln(&lettre)

		if strings.Contains(mot, lettre) {
			RevelerLettresSimilaires(mot, lettre, lettresTrouvees)

			fmt.Println("Mot à deviner :", string(lettresTrouvees))

			if !strings.Contains(string(lettresTrouvees), "_") {
				fmt.Println("Félicitations ! Vous avez trouvé le mot :", mot)
				break
			}
		} else {
			essais++
			fmt.Printf("La lettre '%s' n'est pas dans le mot.\n", lettre)
			fmt.Println("Il vous reste", maxEssais-essais, "chances")
			AfficherHangman(hangs, essais, lignesParEssai)
		}
	}
	if essais == 10 {
		fmt.Println("Désolé, vous avez utilisé tous vos essais. Le mot était :", mot)
	}
}

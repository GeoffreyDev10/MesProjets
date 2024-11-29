package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func lireMotsDepuisFichier(fichier string) ([]string, error) {
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

func motAleatoire(mots []string) string {
	rand.Seed(time.Now().UnixNano())
	indexAleatoire := rand.Intn(len(mots))
	return mots[indexAleatoire]
}

func camouflerMot(mot string, lettresAReveler int) string {
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

func revelerLettres(mot string, lettresTrouvees []rune, motCamoufle string) string {
	for i, l := range motCamoufle {
		if l != '_' {
			lettresTrouvees[i] = rune(mot[i])
		}
	}
	return string(lettresTrouvees)
}

func revelerLettresSimilaires(mot string, lettre string, lettresTrouvees []rune) {
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

func afficherHangman(hangs []string, essais int, lignesParEssai int) {
	debut := (essais - 1) * lignesParEssai
	fin := debut + lignesParEssai
	if fin > len(hangs) {
		fin = len(hangs)
	}
	for i := debut; i < fin; i++ {
		fmt.Println(hangs[i])
	}
}

func lireASCIIFile(fichier string) (map[rune][]string, error) {
	asciiMap := make(map[rune][]string)
	file, err := os.Open(fichier)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var currentChar rune
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			if currentChar != 0 {
				asciiMap[currentChar] = lines
			}
			currentChar = rune(line[1])
			lines = nil
		} else {
			lines = append(lines, line)
		}
	}
	if currentChar != 0 {
		asciiMap[currentChar] = lines
	}

	return asciiMap, scanner.Err()
}

func afficherASCIIArt(lettres []rune, asciiMap map[rune][]string) {
	var outputLines []string
	var maxLines int

	for _, l := range lettres {
		if art, exists := asciiMap[l]; exists {
			if len(outputLines) == 0 {
				outputLines = make([]string, len(art))
				maxLines = len(art)
			}
			for i := 0; i < maxLines; i++ {
				if i < len(art) {
					outputLines[i] += art[i] + "  "
				} else {
					outputLines[i] += "     "
				}
			}
		} else {
			for i := 0; i < maxLines; i++ {
				outputLines[i] += "     "
			}
		}
	}

	for _, line := range outputLines {
		fmt.Println(line)
	}
}

func main() {
	mots, err := lireMotsDepuisFichier("words.txt")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	hangs, err := HangmanPositions("hangman.txt")
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	asciiMap, err := lireASCIIFile("standard.txt")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier ASCII :", err)
		return
	}

	mot := motAleatoire(mots)
	nLettresAReveler := len(mot)/2 - 1
	motCamoufle := camouflerMot(mot, nLettresAReveler)
	lettresTrouvees := []rune(strings.Repeat("_", len(mot)))
	lettresTrouvees = []rune(revelerLettres(mot, lettresTrouvees, motCamoufle))

	maxEssais := 10
	essais := 0
	lignesParEssai := 8
	fmt.Println("Bonne chance, tu as ", maxEssais-essais, " chances")
	fmt.Println(motCamoufle)
	afficherASCIIArt(lettresTrouvees, asciiMap)

	for essais < maxEssais {
		fmt.Print("Entrez une lettre : ")
		var lettre string
		fmt.Scanln(&lettre)

		if strings.Contains(mot, lettre) {
			revelerLettresSimilaires(mot, lettre, lettresTrouvees)
			afficherASCIIArt(lettresTrouvees, asciiMap)

			if !strings.Contains(string(lettresTrouvees), "_") {
				fmt.Println("Félicitations ! Vous avez trouvé le mot :", mot)
				break
			}
		} else {
			essais++
			fmt.Printf("La lettre '%s' n'est pas dans le mot.\n", lettre)
			fmt.Println("Il vous reste", maxEssais-essais, "chances")
			afficherHangman(hangs, essais, lignesParEssai)
		}
	}

	if essais == maxEssais {
		fmt.Println("Désolé, vous avez utilisé tous vos essais. Le mot était :", mot)
	}
}

package main

import (
	"HangmanWeb/Hangman"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"text/template"
)

// go run hangmanweb.go
// se connecter au local host

type GameState struct {
	Word           string `json:"word"`
	MaskedWord     string `json:"masked_word"`
	AttemptsLeft   int    `json:"attempts_left"`
	GuessedLetters string `json:"guessed_letters"`
	CurrentStage   int    `json:"current_stage"`
	Win            bool   `json:"win"`
	Lose           bool   `json:"lose"`
}

var (
	gameState GameState
	gameMutex sync.Mutex
)

func loadGame() error {
	words, err := Hangman.LireMotsDepuisFichier("Hangman/words.txt")
	if err != nil {
		return err
	}

	hangs, err := Hangman.HangmanPositions("Hangman/hangman.txt")
	if err != nil {
		return err
	}
	fmt.Println("Positions hangman chargées : ", len(hangs))

	word := Hangman.MotAleatoire(words)
	maskedWord := Hangman.CamouflerMot(word, len(word)/2-1)

	gameState = GameState{
		Word:           word,
		MaskedWord:     maskedWord,
		AttemptsLeft:   10,
		GuessedLetters: "",
		CurrentStage:   0,
	}

	return nil
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gameState); err != nil {
		http.Error(w, "Erreur interne du serveur1", http.StatusInternalServerError)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	// Empêche les tentatives après la fin du jeu
	if gameState.Win || gameState.Lose {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	letter := r.FormValue("letter")
	if letter == "" || len(letter) > 1 {
		http.Error(w, "Veuillez fournir une lettre valide", http.StatusBadRequest)
		return
	}

	letter = strings.ToLower(letter)
	if strings.Contains(gameState.GuessedLetters, letter) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	gameState.GuessedLetters += letter
	lettresTrouvees := []rune(gameState.MaskedWord)

	if strings.Contains(gameState.Word, letter) {
		Hangman.RevelerLettresSimilaires(gameState.Word, letter, lettresTrouvees)
		gameState.MaskedWord = string(lettresTrouvees)
	} else {
		gameState.CurrentStage++
		gameState.AttemptsLeft--
	}

	if gameState.MaskedWord == gameState.Word {
		gameState.Win = true
	} else if gameState.AttemptsLeft <= 0 {
		gameState.Lose = true
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	// Détermine le chemin de l'image en fonction du nombre d'erreurs
	imagePath := fmt.Sprintf("/static/hangman_images/%d.jpeg", gameState.CurrentStage)

	// Chargement du template HTML
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template HTML : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		GameState GameState
		ImagePath string
	}{
		GameState: gameState,
		ImagePath: imagePath,
	}

	// Exécution du template avec les données
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template : "+err.Error(), http.StatusInternalServerError)
	}
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	// Réinitialiser l'état du jeu
	err := loadGame()
	if err != nil {
		http.Error(w, "Erreur lors du redémarrage du jeu : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page principale après le redémarrage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	err := loadGame()
	if err != nil {
		panic("Erreur lors du chargement du jeu : " + err.Error())
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/game", gameStateHandler)
	http.HandleFunc("/guess", guessHandler)
	http.HandleFunc("/restart", restartHandler)

	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur lors du démarrage du serveur : ", err)
	}
}

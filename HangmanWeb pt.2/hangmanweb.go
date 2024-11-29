package main

import (
	"HangmanWeb/Hangman"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

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
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	letter := r.URL.Query().Get("letter")
	if letter == "" || len(letter) > 1 {
		http.Error(w, "Veuillez fournir une lettre valide", http.StatusBadRequest)
		return
	}

	letter = strings.ToLower(letter)
	if strings.Contains(gameState.GuessedLetters, letter) {
		http.Error(w, "Lettre déjà devinée", http.StatusConflict)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gameState); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

func main() {
	err := loadGame()
	if err != nil {
		panic("Erreur lors du chargement du jeu : " + err.Error())
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/game", gameStateHandler)
	http.HandleFunc("/guess", guessHandler)

	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur lors du démarrage du serveur : ", err)
	}
}

package main

import (
	"HangmanWeb/Hangman"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"text/template"
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
	pseudo    string
)

// chargement du jeu et initialisation de l'état du jeu
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

// Fonction de réinitialisation du jeu
func restartHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	err := loadGame()
	if err != nil {
		http.Error(w, "Erreur lors du redémarrage du jeu : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page principale après le redémarrage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Gestion de la page principale (jeu)
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie si le pseudo est défini en mémoire
	if pseudo == "" {
		// Si aucun pseudo n'est défini, redirige vers la page de démarrage
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	maxStages := 10 // Limite du nombre d'images
	if gameState.CurrentStage >= maxStages {
		gameState.CurrentStage = maxStages - 1
	}
	imagePath := fmt.Sprintf("/static/hangman_images/%d.jpeg", gameState.CurrentStage)

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template HTML : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Structure de données à passer au template
	data := struct {
		Pseudo    string
		GameState GameState
		ImagePath string
	}{
		Pseudo:    pseudo,
		GameState: gameState,
		ImagePath: imagePath,
	}

	// Exécuter le template avec les données passées
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template : "+err.Error(), http.StatusInternalServerError)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	// empêche les tentatives après la fin du jeu
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

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Affiche le formulaire pour entrer le pseudo
		tmpl, err := template.ParseFiles("template/index2.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement du template : "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Récupère le pseudo soumis
		pseudo = r.FormValue("pseudo")
		if pseudo == "" {
			http.Error(w, "Veuillez entrer un pseudo valide", http.StatusBadRequest)
			return
		}

		// Redirige vers la page de jeu
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// vérifie si le pseudo est défini en mémoire
	if pseudo == "" {
		// si aucun pseudo n'est défini, redirige vers la page de démarrage
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	// si le pseudo est défini, afficher la page de jeu
	maxStages := 10 // limite du nombre d'images
	if gameState.CurrentStage >= maxStages {
		gameState.CurrentStage = maxStages - 1
	}
	imagePath := fmt.Sprintf("/static/hangman_images/%d.jpeg", gameState.CurrentStage)

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template HTML : "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Pseudo    string
		GameState GameState
		ImagePath string
	}{
		Pseudo:    pseudo,
		GameState: gameState,
		ImagePath: imagePath,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template : "+err.Error(), http.StatusInternalServerError)
	}
}

// démarre le serveur
func main() {
	err := loadGame()
	if err != nil {
		panic("Erreur lors du chargement du jeu : " + err.Error())
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// définition des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startGameHandler)
	http.HandleFunc("/game", gameStateHandler)
	http.HandleFunc("/guess", guessHandler)
	http.HandleFunc("/restart", restartHandler)

	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur lors du démarrage du serveur : ", err)
	}
}

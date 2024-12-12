package main

import (
	"HangmanWeb/Hangman"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
)

// Définition de l'état du jeu
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

var pseudo string

// Fonction de chargement du jeu
func loadGame() error {
	// Vérification des fichiers nécessaires
	_, err := os.Stat("Hangman/words.txt")
	if err != nil {
		return fmt.Errorf("le fichier des mots est introuvable : %v", err)
	}

	_, err = os.Stat("Hangman/hangman.txt")
	if err != nil {
		return fmt.Errorf("le fichier des positions hangman est introuvable : %v", err)
	}

	// Chargement des mots et des positions
	words, err := Hangman.LireMotsDepuisFichier("Hangman/words.txt")
	if err != nil {
		return err
	}

	hangs, err := Hangman.HangmanPositions("Hangman/hangman.txt")
	if err != nil {
		return err
	}
	fmt.Println("Positions hangman chargées : ", len(hangs))

	// Choix d'un mot aléatoire et création de l'état du jeu
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

// Route pour afficher l'état du jeu en JSON
func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gameState); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}

	// Charge la page du jeu si la méthode est POST
	if r.Method == http.MethodPost {
		// Récupère le pseudo soumis
		pseudo = r.FormValue("pseudo")
	}

	// Charge la page de jeu
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Erreur de chargement de la page", http.StatusInternalServerError)
		return
	}

	// Passe le pseudo à la page de jeu
	data := map[string]string{"Pseudo": pseudo}
	tmpl.Execute(w, data)
}

// Route pour gérer les devinettes
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

	// Mise à jour de l'état du jeu
	if strings.Contains(gameState.Word, letter) {
		Hangman.RevelerLettresSimilaires(gameState.Word, letter, lettresTrouvees)
		gameState.MaskedWord = string(lettresTrouvees)
	} else {
		gameState.CurrentStage++
		gameState.AttemptsLeft--
	}

	// Vérification de la victoire ou défaite
	if gameState.MaskedWord == gameState.Word {
		gameState.Win = true
	} else if gameState.AttemptsLeft <= 0 {
		gameState.Lose = true
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if pseudo == "" {
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	// Limite le nombre d'images pour le hangman
	maxStages := 10
	if gameState.CurrentStage >= maxStages {
		gameState.CurrentStage = maxStages - 1
	}
	imagePath := fmt.Sprintf("/static/hangman_images/%d.jpeg", gameState.CurrentStage)

	// Charge la page du jeu
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template HTML : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Structure les données pour la page
	data := struct {
		Pseudo    string
		GameState GameState
		ImagePath string
	}{
		Pseudo:    pseudo,
		GameState: gameState,
		ImagePath: imagePath,
	}

	// Exécute le template avec les données
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template : "+err.Error(), http.StatusInternalServerError)
	}
}

// Route pour redémarrer le jeu
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

// Route pour afficher le formulaire d'entrée du pseudo
func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Affiche le formulaire pour entrer le pseudo
		tmpl, err := template.ParseFiles("index2.html")
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

		// Redirige vers la page du jeu
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Point d'entrée principal
func main() {
	err := loadGame()
	if err != nil {
		panic("Erreur lors du chargement du jeu : " + err.Error())
	}

	// Sert les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Définit les routes
	http.HandleFunc("/", homeHandler)           // Page d'accueil
	http.HandleFunc("/start", startGameHandler) // Formulaire de pseudo
	http.HandleFunc("/game", gameStateHandler)  // Page principale du jeu
	http.HandleFunc("/guess", guessHandler)     // Gestion des devinettes
	http.HandleFunc("/restart", restartHandler) // Redémarrage du jeu

	// Démarre le serveur
	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur lors du démarrage du serveur : ", err)
	}
}

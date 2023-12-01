package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

var wordTab []string

func getRandomWord() string {
	data, err := os.ReadFile("words.txt") //sert à lire le fichier words.txt et le stocke dans data (tableau de byte)
	if err != nil {
		log.Fatal(err)
	}
	motRandom := strings.Split(string(data), "\n") // split le tableau de byte en tableau de string en fonction du \n (retour à la ligne)
	random := rand.Intn(len(motRandom))            // génère un nombre aléatoire entre 0 et la taille du tableau de string
	return motRandom[random]                       // retourne le mot aléatoire
}

func revealLetters(word string) string { // fonction qui permet de révéler la moitié des lettres du mot
	initialWord := ""
	for i := 0; i < len(word); i++ { // boucle qui permet de créer un mot avec des _ à la place des lettres
		initialWord = initialWord + "_" // initialWord = "______"
	}
	counter := 0
	lettersToReavel := len(word)/2 - 1 // permet de révéler la moitié des lettres du mot
	for counter <= lettersToReavel {   // boucle qui permet de révéler la moitié des lettres du mot
		randomNumber := rand.Intn(len(word))                  // génère un nombre aléatoire entre 0 et la taille du mot
		if string(initialWord[randomNumber]) == string("_") { // si la lettre n'est pas déjà révélée
			counter++
			initialWord = initialWord[:randomNumber] + string(word[randomNumber]) + initialWord[(randomNumber+1):] // on remplace le _ par la lettre du mot
		}
	}
	return initialWord
}

func printAttempt(initialWord string) {
	for i := 0; i < len(initialWord); i++ { // boucle qui permet d'afficher le mot avec les lettres révélées
		if string(initialWord[i]) == "_" { // si la lettre n'est pas révélée
			print(string(initialWord[i]), " ") // on affiche un _
		} else {
			print(strings.ToUpper(string(initialWord[i])), " ") // sinon on affiche la lettre en majuscule
		}
	}

}

func hangmanDrawing(attempts int) {
	f, err := os.Open("hangman.txt") // ouvre le fichier hangman.txt

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f) // permet de lire le fichier hangman.txt
	buf := make([]byte, 71)      // tableau de byte de taille 71
	hangmanlist := []string{}    // tableau de string vide
	for {
		n, err := reader.Read(buf) // lit le fichier hangman.txt et stocke les données dans buf n = nombre de byte lus

		if err != nil {

			if err != io.EOF {

				log.Fatal(err)
			}

			break
		}

		hangmanlist = append(hangmanlist, string(buf[0:n])) // ajoute les données de buf dans hangmanlist
	}
	print(hangmanlist[9-attempts]) // affiche le dessin du pendu en fonction du nombre d'essais restants
	return
}

func startGame(initialWord string, word string) { // fonction qui permet de jouer
	attempts := 10
	for attempts > 0 { // boucle qui permet de jouer tant qu'il reste des essais
		fmt.Println("Choose:")                 // demande à l'utilisateur de choisir une lettre
		reader := bufio.NewReader(os.Stdin)    // permet de lire ce que l'utilisateur écrit dans la console
		letter, err := reader.ReadString('\n') // stocke ce que l'utilisateur écrit dans letter
		if err != nil {
			log.Fatal(err)
		}

		underlines := 0
		foundLetter := false
		for i := 0; i < len(initialWord); i++ { // boucle qui permet de vérifier si la lettre choisie est dans le mot
			if word[i] == []byte(letter)[0] { // si la lettre choisie est dans le mot
				initialWord = initialWord[:i] + string(word[i]) + initialWord[(i+1):] // on remplace le _ par la lettre du mot
				foundLetter = true
			}
			if string(initialWord[i]) == "_" { // si la lettre n'est pas révélée
				underlines++
			}
		}
		if foundLetter == false { // si la lettre choisie n'est pas dans le mot
			attempts--
			fmt.Println("Not present in the word,", attempts, "attempts remaining") // on affiche le nombre d'essais restants
			hangmanDrawing(attempts)                                                // on affiche le dessin du pendu en fonction du nombre d'essais restants
		} else {
			printAttempt(initialWord) // on affiche le mot avec les lettres révélées
		}
		if attempts == 0 {
			fmt.Println("You lost !")
			fmt.Println("The word was:", word)
		}
		if underlines == 0 {
			fmt.Println("Congrats !")
			attempts = 0
		}
	}
}

func LancerPendu() {
	word := getRandomWord()                         // stocke le mot aléatoire dans word
	initialWord := revealLetters(word)              // stocke le mot avec les lettres révélées dans initialWord
	fmt.Println("Good Luck, you have 10 attempts.") // affiche le nombre d'essais restants
	printAttempt(initialWord)                       // affiche le mot avec les lettres révélées
	startGame(initialWord, word)                    // permet de jouer
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const wordslistFile = "words.txt"

func randomWord(words []string) string {
	rand.Seed(time.Now().Unix())
	return words[rand.Intn(len(words))]
}

func HangmanWord() string {
	f, err := os.Open(wordslistFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scannerlist := []string{}
	for scanner.Scan() {
		scannerlist = append(scannerlist, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return randomWord(scannerlist)
}

func HangmanP(attempts int) string {
	f, err := os.Open("hangman.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	buf := make([]byte, 71)
	hangmanlist := []string{}
	for {
		n, err := reader.Read(buf)

		if err != nil {

			if err != io.EOF {

				log.Fatal(err)
			}

			break
		}

		hangmanlist = append(hangmanlist, string(buf[0:n]))
	}
	return hangmanlist[attempts]
}

func main() {
	wordToGuess := HangmanWord()
	guesses := make(map[string]bool)
	maxAttempts := 10
	attempts := 0

	for n := len(wordToGuess)/2 - 1; n >= 0; n-- {
		letter := ReturnLetter(wordToGuess)
		guesses[letter] = true
	}

	for {
		displayWord(wordToGuess, guesses)
		if isWordGuessed(wordToGuess, guesses) {
			fmt.Println("Félicitations! Vous avez deviné le mot.")
			break
		}
		if attempts >= maxAttempts {
			fmt.Println("Vous avez épuisé vos tentatives. Le mot était:", wordToGuess)
			break
		}
		fmt.Print("Devinez une lettre: ")
		var guess string
		fmt.Scanln(&guess)
		guess = strings.ToLower(guess)
		if _, found := guesses[guess]; found {
			fmt.Println("Vous avez déjà deviné cette lettre.")
			continue
		}
		if !isLetterInWord(guess, wordToGuess) {
			fmt.Println("Il vous reste", 9-attempts, "tentatives")
			attempts++
		}
		fmt.Println(HangmanP(attempts))
		guesses[guess] = true
	}

}

func ReturnLetter(word string) string {
	num := rand.Intn(len(word))
	word2 := []byte(word)
	return string(word2[num])
}

func displayWord(word string, guesses map[string]bool) {
	for _, letter := range word {
		if _, found := guesses[string(letter)]; found {
			letter = rune(ToUpper(string(letter))[0])
			fmt.Printf("%c ", letter)
		} else {
			fmt.Print("_ ")
		}
	}
	fmt.Println()
}

func ToUpper(s string) string {
	h := []rune(s)
	result := ""
	for i := 0; i <= len(h)-1; i++ {
		if (h[i] >= 'a') && (h[i] <= 'z') {
			h[i] = h[i] - 32
		}
		result += string(h[i])
	}
	return result
}

func isWordGuessed(word string, guesses map[string]bool) bool {
	for _, letter := range word {
		if _, found := guesses[string(letter)]; !found {
			return false
		}
	}
	return true
}

func isLetterInWord(letter string, word string) bool {
	for _, l := range word {
		if string(l) == letter {
			return true
		}
	}
	return false
}
package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var listedemot = "Boss"
var deja = []string{}
var start = true
var mot string
var motcacher string
var vie int = 10
var endmessage string
var imagepath = "images/vie10.png"
var images = []string{
	"images/vie0.png",
	"images/vie1.png",
	"images/vie2.png",
	"images/vie3.png",
	"images/vie4.png",
	"images/vie5.png",
	"images/vie6.png",
	"images/vie7.png",
	"images/vie8.png",
	"images/vie9.png",
	"images/vie10.png",
}
var imagegif = "images/combat.gif"
var gifs = []string{
	"images/combat.gif",
	"images/combatv.gif",
	"images/combatg.gif",
	"images/combatd.gif",
	"images/combatp.gif",
}

type Hangman struct { // structure contenant les information a envoye au template
	Deja       []string
	Mot        string
	Vie        int
	Endmessage string
	Imagepath  string
	Imagegif   string
}

func Aleatoire(liste string) string {
	data, err := os.ReadFile(liste + ".txt") // Lecture du fichier "words.txt"
	if err != nil {
		log.Fatal(err) // En cas d'erreur, arrête le programme et affiche l'erreur
	}
	s := strings.Split(string(data), "\n") // Séparation des lignes du fichier en un tableau de chaînes de caractères
	random := rand.Intn(len(s))            // Génération d'un indice aléatoire dans la plage des lignes du fichier
	return ToUpper(s[random])              // Retourne le mot aléatoire
}

// revealLetter révèle certaines lettres du mot au début du jeu
func revealLetter(word string) string {
	initialWord := "" // Initialise une chaîne de caractères vide pour stocker le mot partiellement révélé
	for i := 0; i < len(word); i++ {
		initialWord = initialWord + "_" // Remplit initialWord avec des underscores pour chaque lettre du mot
	}
	aLettre := len(word)/2 - 1 // Calcule le nombre de lettres à révéler (moitié du mot - 1)
	compteur := 1              // Initialise un compteur pour suivre le nombre de lettres révélées
	for compteur <= aLettre {
		walid := rand.Intn(len(word)) // Génère un indice aléatoire pour choisir une lettre du mot
		// Vérifie si l'indice est valide et si la lettre correspondante dans initialWord est encore non révélée
		if walid >= 0 && walid < len(initialWord) && string(initialWord[walid]) == "_" {
			compteur++                                                                        // Incrémente le compteur de lettres révélées
			initialWord = initialWord[:walid] + string(word[walid]) + initialWord[(walid+1):] // Remplace l'underscore par la lettre dans initialWord
		}
	}
	return initialWord // Retourne le mot partiellement révélé
}
func restart() {
	mot = ""
	mot = Aleatoire(listedemot)
	print(mot) // les test iront plus vite
	motcacher = revealLetter(mot)
	vie = 10
	deja = []string{}
	endmessage = ""
	imagegif = gifs[0]
	start = false
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

func PasUtilise(l string) bool {
	for i := 0; i < len(deja); i++ {
		if deja[i] == l {
			return false
		}
	}
	return true
}

func main() {

	if start {
		restart()
	}
	fs := http.FileServer(http.Dir("images")) // recuperation du dossier images et son contenu
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	fsCss := http.FileServer(http.Dir("./css")) // recuperation du  dossier css et son contenu
	http.Handle("/css/", http.StripPrefix("/css/", fsCss))
	tmpl := template.Must(template.ParseFiles("HangmanWeb.html")) // recuperation du template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Hangman{
			Deja:       deja,
			Mot:        motcacher,
			Vie:        vie,
			Endmessage: endmessage,
			Imagepath:  imagepath,
			Imagegif:   imagegif,
		}
		liste := r.FormValue("liste") // choix de la liste de mots
		if liste != "" {
			listedemot = liste
		}
		nouvellepartie := r.FormValue("game") // reset
		if nouvellepartie == "Nouveau" {
			restart()
		}
		tLettre := false                         // Tlettre = True Lettre pour verifier si la lettre est dans le mot
		lettre := ToUpper(r.FormValue("lettre")) //Input de la lettre mis en maj
		if lettre != "" && PasUtilise(lettre) {  //Vérifie si une lettre est envoyée et si elle n'a pas déjà été utilisée
			deja = append(deja, lettre)
			for i := 0; i < len(motcacher); i++ { // Parcours du mot initial
				if mot[i] == []byte(lettre)[0] { // Vérifie si la lettre proposée est présente dans le mot
					motcacher = motcacher[:i] + string(mot[i]) + motcacher[(i+1):] // Met à jour le mot initial avec la lettre trouvée
					tLettre = true
					imagegif = gifs[1]
					if mot == motcacher {
						endmessage = "Vous avez vaincu"
						data.Endmessage = endmessage
						imagegif = gifs[2]
						data.Imagegif = imagegif
					}
				}
			}
			if tLettre == false { //perte de vie et changement d'images
				vie--
				imagepath = images[vie]
				imagegif = gifs[3]
			}

		}
		if vie == 0 { // arret de la partie
			endmessage = "Vous avez péri(e)"
			data.Endmessage = endmessage
			imagegif = gifs[4]
			data.Imagegif = imagegif
		}

		imagepath = images[vie] //Mise a jour des information de la struture
		data.Imagepath = imagepath
		data.Imagegif = imagegif
		data.Endmessage = endmessage
		data.Deja = deja
		data.Mot = motcacher
		data.Vie = vie
		tmpl.Execute(w, data) // execution des nouvelles information sur le template
	})
	http.ListenAndServe(":80", nil) // diffusion sur le localhost
}

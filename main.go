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
var imagecombat = "images/combat.gif"
var gifs = []string{
	"images/combat.gif",
	"images/combatv.gif",
	"images/combatg.gif",
	"images/combatd.gif",
	"images/combatp.gif",
}

type Hangman struct {
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
	return s[random]                       // Retourne le mot aléatoire
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
	mot = Aleatoire(listedemot)
	motcacher = revealLetter(mot)
	vie = 10
	deja = []string{}
	endmessage = ""
	imagecombat = gifs[0]
	start = false
}

func main() {

	if start {
		restart()
	}
	print(mot + "/")
	fs := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	fsCss := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", fsCss))
	tmpl := template.Must(template.ParseFiles("HangmanWeb.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Hangman{
			Deja:       deja,
			Mot:        motcacher,
			Vie:        vie,
			Endmessage: endmessage,
			Imagepath:  imagepath,
			Imagegif:   imagecombat,
		}
		liste := r.FormValue("liste")
		if liste != "" {
			listedemot = liste
		}
		nouvellepartie := r.FormValue("game")
		if nouvellepartie == "Nouveau" {
			restart()
		}
		tLettre := false
		lettre := r.FormValue("lettre")
		if lettre != "" {
			deja = append(deja, lettre)
			for i := 0; i < len(motcacher); i++ { // Parcours du mot initial
				if mot[i] == []byte(lettre)[0] { // Vérifie si la lettre proposée est présente dans le mot
					motcacher = motcacher[:i] + string(mot[i]) + motcacher[(i+1):] // Met à jour le mot initial avec la lettre trouvée
					tLettre = true
					imagecombat = gifs[1]
					if mot == motcacher {
						endmessage = "Vous avez vaincu"
						data.Endmessage = endmessage
						imagecombat = gifs[2]
						data.Imagegif = imagecombat
					}
				}
			}
			if tLettre == false {
				vie--
				imagepath = images[vie]
				imagecombat = gifs[3]
			}

		}
		if vie == 0 {
			endmessage = "Vous avez péri(e)"
			data.Endmessage = endmessage
			imagecombat = gifs[4]
			data.Imagegif = imagecombat
		}

		imagepath = images[vie]
		data.Imagepath = imagepath
		data.Imagegif = imagecombat
		data.Deja = deja
		data.Mot = motcacher
		data.Vie = vie
		tmpl.Execute(w, data)
		print(mot + "bou")
	})
	http.ListenAndServe(":80", nil)
}

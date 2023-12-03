package main

import (
	"fmt"
	"github.com/laspp/PS-2023/vaje/naloga-1/koda/xkcd"
	"sync"
	"time"
	"flag"
	"strings"
	"unicode"
)

func opravi (zacetek int, konec int) {
	
	defer wg.Done()

	for i := zacetek; i <= konec; i++ {
		poizvedba, napaka := xkcd.FetchComic(i)

		if napaka == nil {

			isPunctation := func(r rune) bool {
				return unicode.IsPunct(r)
			}

			title := strings.Map(func(r rune) rune {
				if isPunctation(r) {
					return -1
				}
				return r
			}, poizvedba.Title)
			transcript := strings.Map(func(r rune) rune {
				if isPunctation(r) {
					return -1
				}
				return r
			}, poizvedba.Transcript)
			tooltip := strings.Map(func(r rune) rune {
				if isPunctation(r) {
					return -1
				}
				return r
			}, poizvedba.Tooltip)

			if transcript != "" {
				tooltip = ""
			}

			besede := strings.Fields(title + " " + transcript + " " + tooltip)

			
			var izbraneBesede []string
			
			for _, beseda := range besede {
				if len(beseda) >= 4 {
					izbraneBesede = append(izbraneBesede, strings.ToLower(beseda))
				}
			}

			kanal <- izbraneBesede
		}

	}
}


/**
  * Funckija delavcem pravično dodeli stripe
  **/
func razdeli (steviloStripov, steviloDelavcev int) map[int]int {

	razdelitev := make(map[int]int)

	stStripovNaNit, preostanek := steviloStripov / steviloDelavcev, steviloStripov % steviloDelavcev

	for i := 1; i <= steviloDelavcev; i++ {
		razdelitev[i] = stStripovNaNit
		if preostanek > 0 {
			razdelitev[i]++
			preostanek--
		}
	}

	return razdelitev
}

/**
  * Funkcija vrne zaporedno številko najnovejšega stripa 
  * ker se številčenje začne s številko 1 je to hkrati tudi število stripov
  **/
func vrniSteviloStripov() (int) {
	
	st, err := xkcd.FetchComic(0)
	
	if err == nil {
		return st.Id
	}
	
	return -1
}

func izpisiNajpogostejse(seznam map[string]int) {

	// razvrscanje := time.Now()

	steviloZmagovalcev := 15	
	if len(rezultat) < steviloZmagovalcev {
		steviloZmagovalcev = len(rezultat)
	}
	
	for i := 0; i < steviloZmagovalcev; i++ {
		
		var najvecPojavitev int = 0
		var najboljsi string = ""

		for element := range seznam {
			if najvecPojavitev < seznam[element] {
				najboljsi = element
				najvecPojavitev = seznam[element]
			}
		}

		fmt.Printf("%s, %d\n", najboljsi, najvecPojavitev)
		delete(seznam, najboljsi)
	}

	// defer fmt.Println("Razporejanje: ", time.Now().Sub(razvrscanje))
}
	
var wg sync.WaitGroup
var kanal chan []string
var rezultat map[string]int

func main () {

	// Začni štoparico!
	startTime := time.Now()

	// Preberi število opravil, ki so vključena v program.
	steviloDelavcev := flag.Int("cpus-per-task", 1, "Number of CPUs per task")
	flag.Parse()
	

	// Ugotovimo koliko stripov obstaja
	var steviloStripov int = vrniSteviloStripov()
	// var steviloStripov int = 100

	// Uporabi funkcijo, ki delavcem dodeli stripe
	razdelitev := razdeli(steviloStripov, *steviloDelavcev)
	
	// Inicializacija kanala
	kanal = make(chan []string, steviloStripov)

	// Inicializacija slovarja za knočni rezultat
	rezultat = make(map[string]int)

	var indeks int = 1

	for i := 0; i < *steviloDelavcev; i++ {
		// Dodaj nit na seznam niti, ki se jih čaka
		wg.Add(1)

		// Z vsako nitjo posebaj kliči metodo
		go opravi(indeks, indeks + razdelitev[i + 1] - 1)
		
		// Premakni indeks ki označuje zaporedno številko stripa
		indeks = indeks + razdelitev[i + 1]
	}

	for i := 0; i < steviloStripov; i++ {
		
		for _, beseda := range <- kanal {
			rezultat[beseda]++
		}
	}

	wg.Wait()

	fmt.Println("Število stripov: ", steviloStripov)

	izpisiNajpogostejse(rezultat)
	
	defer fmt.Println("Pretekel čas: ", time.Now().Sub(startTime))
}

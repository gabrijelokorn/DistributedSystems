package main

import (
	"fmt"
	"flag"
	"time"
	"sync"
	"strings"
	"unicode"
	"github.com/laspp/PS-2023/vaje/naloga-2/koda/socialNetwork"
)

// Funkcija, ki jo bodo izvajali delavci.
// Delavec bo prejemal modrosti in jih indeksiral.
// Delavec bo prejemal zahtevke za iskanje in vrnil rezultate.
func prejmiSporocila(idDelavec int) {

	// Funkcija wg.Done() sporoči, da je delavec končal z delom.
	defer wg.Done()

	for {

		// Preveri, če je na voljo sporočilo.
		select {

			// Če je na voljo sporočilo, ga prejmi.
			case sporocilo := <- kanalIndex:
				var izbraneBesede []string

				// Indeksiraj modrost.
				func() {

					// Funkcija, ki preveri, ali je znak ločilo.
					isPunctation := func(r rune) bool {
						return unicode.IsPunct(r)
					}

					// Odstrani ločila iz modrosti.
					modrost := strings.Map(func(r rune) rune {
						if isPunctation(r) {
							return -1
						}
						return r
					}, sporocilo.Data)

					// Razdeli modrost na besede.
					// Funkcija strings.Fields() vrne rezultat kot rezino.
					besede := strings.Fields(modrost)

					// Izberi besede, ki so dolge vsaj 4 znake.
					for _, beseda := range besede {
						if len(beseda) >= 4 {
							// Funkcija strings.ToLower() vrne besedo v malem tisku.
							// Funkcija append() doda element na konec rezine.
							// Funkcija append() vrne rezino, ki je lahko različno dolga od vhodne rezine.
							// Funkcija append() lahko sprejme poljubno število argumentov.
							izbraneBesede = append(izbraneBesede, strings.ToLower(beseda))
						}
					}
					
				}()

				// Indeksiraj besede.
				// lock.Lock() blokira izvajanje, dokler ne pridobi zaklepa.
				lock.Lock()
				for _, beseda := range izbraneBesede {
					// Funkcija append() doda element na konec rezine.
					slovar[beseda] = append(slovar[beseda], sporocilo.Id)
				}
				// lock.Unlock() sprosti zaklep.
				lock.Unlock()
				
				// Izpiši ID delavca in število indeksiranih besed.
				fmt.Println(sporocilo.Data)

			default:
				// Če ni na voljo sporočilo, preveri, če je na voljo zahteva za iskanje.
				// Če je na voljo zahteva za iskanje, jo prejmi.
				select {
					case sporocilo := <- kanalSearch:
						// lock.Lock() blokira izvajanje, dokler ne pridobi zaklepa.
						lock.Lock()
						// Izpiši besedo in indekse modrosti, ki vsebujejo to besedo.
						fmt.Println(sporocilo.Data, ": ", slovar[sporocilo.Data])
						// lock.Unlock() sprosti zaklep.
						lock.Unlock()
					default:
						// Če ni na voljo zahteva za iskanje, preveri, če je na voljo zahteva za ustavitev.
						// Če je na voljo zahteva za ustavitev, prejmi sporočilo in se ustavi.
						case <- kanalOver:
							return

				}
		}
	}
}

// wg je objekt, ki ga uporabljamo za sinhronizacijo.
var wg sync.WaitGroup
// Kanali, ki jh uporabljamo za komunikacijo med delavci in glavno gorutino.
var kanalSearch chan socialNetwork.Task
var kanalIndex chan socialNetwork.Task
var kanalOver chan string

// Slovar, ki bo hranil besede v modrostih in njihove ID-je.
var slovar = make(map[string][]uint64)

// lock je objekt, ki ga uporabljamo za sinhronizacijo.
var lock sync.Mutex

func main () {
	// Uporabimo flag package, ki je del standardne knjižnice.
	// Omogoča nam, da programu podamo argumente iz ukazne vrstice.
	// V tem primeru bomo programu podali argumente -cpus-per-task in -rate.
	steviloDelavcev := flag.Int("cpus-per-task", 1, "Number of CPUs per task")
	steviloModrostiNaSekundo := flag.Int("rate", 100, "Number of CPUs per task")
	flag.Parse()

	// Inicializacija kanala
	// V našem primeru bomo uporabili kanal za komunikacijo med delavci in glavno gorutino.
	// Velikost kanala je število elementov, ki jih lahko kanal hkrati hrani.
	kanalSearch = make(chan socialNetwork.Task, 3 * *steviloDelavcev)
	kanalIndex = make(chan socialNetwork.Task, 3 * *steviloDelavcev)
	kanalOver = make(chan string, *steviloDelavcev)
	
	// Ustvari nov objekt tipa socialNetworkQ.
	var producer socialNetwork.Q
	// Razmerje med nizko prioriteto in visoko prioriteto.
	producer.New(0.5)

	// Začni štoparico
	startTime := time.Now()

	// Ustvari omejevalnik, ki bo omejeval število modrosti na sekundo.
	omejevalnik := make(chan time.Time, 0)
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(*steviloModrostiNaSekundo))
		// Omejevalnik bo vsako periodo poslal en element na kanal.
		for t := range ticker.C {
			omejevalnik <- t
		}	
	}()

	// Ustvari delavce.
	go func() {
		for {
			sporocilo := <- producer.TaskChan
			<- omejevalnik
			if (sporocilo.TaskType == "search") {
				kanalSearch <- sporocilo
			} else {
				kanalIndex <- sporocilo
			}
			
	}
	}()
		
	for i := 0; i < *steviloDelavcev; i++ {
		wg.Add(1)
			
		go prejmiSporocila(i)
	}
			
	// Producer.Run() bo ustvaril zahtevke in jih poslal na kanal.
	go producer.Run()
	time.Sleep(time.Millisecond * 3000)
	producer.Stop()
	
	// Pošlji zahtevek za ustavitev delavcem.
	for i := 0; i < *steviloDelavcev; i++ {
		kanalOver <- "Done"
	}

	// Počakaj, da se delavci ustavijo.
	wg.Wait()

	// Ustavi omejevalnik.
	elapsed := time.Since(startTime)
    // Izpišemo število generiranih zahtevkov na sekundo
	fmt.Printf("Spam rate: %f MReqs/s\n", float64(producer.N[socialNetwork.LowPriority]+producer.N[socialNetwork.HighPriority])/float64(elapsed.Seconds())/1000000.0)
	
	// Izognitev napaki "declared and not used"
}

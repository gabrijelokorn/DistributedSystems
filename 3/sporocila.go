package main

import (
	"fmt"
	"flag"
	"strconv"
	"sync"
	"math/rand"
	"net"
	"time"
)

func kreirajNaslov (naslov int) *net.UDPAddr {
	naslovString := fmt.Sprintf("localhost:%d", naslov)
	naslovUDP, err := net.ResolveUDPAddr("udp", naslovString)
	checkError(err)
	return naslovUDP
}

func checkError(err error) {
	if err != nil {
		panic(err)
		fmt.Println("---")
	}
}

// Funkcija, ki preveri, če je spročilo že v nabiralniku
func vsebujeSporocilo (sporocilo string, nabiralnik []string) bool {
	for _, element := range nabiralnik {
		if element == sporocilo {
			return true
		}
	}
	return false
}

// Funkcija, ki posluša na svojih vratih
func poslusaj (K int) {
	defer wg.Done()
	/*
	Vsak proces naj na svojih vratih posluša sporočila
	dokler ne preteče vnaprej omejena količina časa.
	Če je v tem času vmes prejel sporočilo, naj ga
	sporcesira in nadaljuje s poslušanjem.
	Če vmes ne prejme sporočila, naj se ustavi.
	*/

	// Določimo lasten naslov
	naslov := kreirajNaslov(naslov)

	// Odpremo povezavo
	conn, err := net.ListenUDP("udp", naslov)
	checkError(err)

	// Definiramo buffer
	buffer := make([]byte, 128)

	// Nastavimo timeout
	timeout := 500 * time.Millisecond
	conn.SetDeadline(time.Now().Add(timeout))

	// Preberemo sporocilo
	len, err := conn.Read(buffer)

	// Preverimo, če je prišlo do napake
	if err != nil {
		conn.Close()
		return
	}

	sporocilo := string(buffer[:len])
	// Preverimo ali je bilo sporocilo ze prejeto
	if vsebujeSporocilo(sporocilo, nabiralnik) {
		if komentiraj {
			fmt.Println("Sporocilo", sporocilo, "je ze bilo prejeto in ga bom ignoriral.")
		}
	} else {
		// Dodamo sporocilo v nabiralnik
		nabiralnik = append(nabiralnik, sporocilo)

		wg.Add(1)
		go govori(sporocilo, K)
	}

	conn.Close()
	wg.Add(1)
	poslusaj(K)

}

// Funkcija, ki preveri, če je dano število v danem slice-u
func vsebujeStevilo(arr []int, target int) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func posljiSporocilo (sporocilo string, prejemnik *net.UDPAddr) {

	// Določimo naslov pošiljatelja
	conn, err := net.DialUDP("udp", nil, prejemnik)
	checkError(err)

	// Po koncu funkcije zapremo povezavo
	defer conn.Close()

	// Pošljemo sporočilo
	_, err = conn.Write([]byte(sporocilo))
	checkError(err)
}

// Funkcija za pošiljanje sporočil
func govori (sporocilo string, K int) {
	defer wg.Done()

	// Slice, v katerem hranimo naslove, na katere smo že poslali sporočilo
	poslano := []int{}
	// Izberemo K naključnih naslovov
	for i := 0; i < K; i++ {
		// Izberemo naključno število
		nakljucniNaslov := rand.Intn(len(naslovi))
		
		// Preverimo, če smo že poslali sporočilo na ta naslov
		if vsebujeStevilo(poslano, nakljucniNaslov) {
			// Če smo, izberemo novo število
			i--
			continue
		}
		
		// Določimo prejemnika
		prejemnik:= kreirajNaslov(naslovi[nakljucniNaslov])
		
		// Pošlji sporočilo
		posljiSporocilo(sporocilo, prejemnik)
		if komentiraj {
			fmt.Println("Iz ", naslov, "Pošiljam sporočilo", sporocilo, "na naslov", naslovi[nakljucniNaslov])
		}

		// Dodaj naslov v seznam poslanih naslovov
		poslano = append(poslano, nakljucniNaslov)
	}
}

// Funkcija, ki razdeli naslove
func razdeliNaslove (id int, N int) {
	for i := 0; i < N; i++ {
		if id == i {
			continue
		}
		naslovi = append(naslovi, 9000 + i)
	}
}

// Funkcija, ki definira sporočila
func definirajSporocila (M int) {
	for i := 0; i < M; i++ {
		sporocila[i] = "<" + strconv.Itoa(i + 1) + ">"
	}

}

var komentiraj bool

var id int
var N int
var M int
var K int

var protokol string

var naslovi []int
var naslov int

var sporocila []string
var nabiralnik []string

var wg sync.WaitGroup


func main () {
	// Definiramo argumente
	idPtr := flag.Int("id", 3, "Identifikator procesa")
	NPtr := flag.Int("N", 1, "Število procesov")
	MPtr := flag.Int("M", 1,"Število sporočil, ki jih pošlje proces z identifikatorjem 0")
	KPtr := flag.Int("K", 1, "Število naslovnikov, ki jih vsak proces izbere")
	flag.Parse()
	id := *idPtr
	N := *NPtr
	M := *MPtr
	K := *KPtr
	fmt.Println("------ Proces", id, "------")

	komentiraj = false
	
	// Določimo kateri protokol se uporablja
	// if K < N - 1 {
	// 	protokol = "govorice"
	// } else {
	// 	protokol = "nestrpno"
	// }
		
	// Definiramo sporočila
	sporocila = make([]string, M)
	definirajSporocila(M)
	
	// Doloci lasten naslov 
	naslov = 9000 + id
	
	// Razdelimo naslove (vsi procesi morajo vedeti vse naslove)
	naslovi = make([]int, 0)
	razdeliNaslove(id, N)
	
	// Definiramo nabiralnik
	nabiralnik = make([]string, 0)

	
	if id != 0 {
		wg.Add(1)
		go poslusaj(K)
	} else {
		// time.Sleep(500 * time.Millisecond)
		for _, sporocilo := range sporocila {
			casSpanja := 100 * time.Millisecond
			time.Sleep(casSpanja)
			wg.Add(1)
			go govori(sporocilo, K)
		}
	}

	wg.Wait()

	for _, sporocilo := range nabiralnik {
		fmt.Println(sporocilo)
	}
		
	fmt.Println("------------")
}
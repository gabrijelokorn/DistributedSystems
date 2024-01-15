// Komunikacija po protokolu gRPC
//
// 		strežnik ustvari in vzdržuje shrambo nalog TodoStorage
//		odjemalec nad shrambo izvaja operacije CRUD
//		datoteka protobufStorage/protobufStorage.proto podaja strukturo sporočil gRPC
//			prevajalnik protoc iz datoteke *.proto ustvari strukture in metode, ki jih uporabljamo pri klicanju oddaljenih metod
//			navodila za prevajanje so v datoteki *.proto
//
// zaženemo strežnik
// 		go run *.go
// zaženemo enega ali več odjemalcev
//		go run *.go -s [ime strežnika] -p [vrata]
// za [ime strežnika] in [vrata] vpišemo vrednosti, ki jih izpiše strežnik ob zagonu
//
// pri uporabi SLURMa lahko s stikalom --nodelist=[vozlišče] določimo vozlišče, kjer naj se program zažene

package main

import (
	"flag"
	"fmt"
)

func main() {
	// Argument 
	hostPtr := flag.String("s", "", "server URL")
	// Argument, ki predstavlja ID procesa
	idPtr := flag.Int("id", 0, "process ID")
	// Argument, ki predstavlja vrata na katerih posluša strežnik
	pPtr := flag.Int("p", 8100, "port number")
	// Argument, ki predstavlja število procesov
	nPtr := flag.Int("n", 10, "number of processes")
	// preberemo argumente iz ukazne vrstice
	flag.Parse()

	ListenPORT := *pPtr + *idPtr
	GetPutPORT := *pPtr + *idPtr + 1
	CommitPORT := *pPtr + *idPtr - 1

	ListenURL := fmt.Sprintf("%v:%v", "", ListenPORT)
	GetPutURL := fmt.Sprintf("%v:%v", *hostPtr, GetPutPORT)
	CommitURL := fmt.Sprintf("%v:%v", *hostPtr, CommitPORT)

	if *idPtr == 0 {
		// glava
		CommitURL = fmt.Sprintf("glava")
	}
	if *idPtr >= *nPtr-1 {
		// rep
		GetPutURL = fmt.Sprintf("rep")
	}

	Proces(ListenURL, GetPutURL, CommitURL)

}

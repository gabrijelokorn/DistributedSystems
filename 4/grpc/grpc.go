package main

import (
	"flag"
	"fmt"
)

func main() {
	// preberemo argumente iz ukazne vrstice
	sPtr := flag.String("s", "", "server URL")
	idPtr := flag.Int("id", 0, "process ID")
	pPtr := flag.Int("p", 8100, "port number")
	nPtr := flag.Int("n", 10, "number of processes")
	flag.Parse()

	// zaženemo strežnik ali odjemalca
	send_to_URL := fmt.Sprintf("%v:%v", *sPtr, *pPtr + 1 + *idPtr)
	read_from_URL := fmt.Sprintf("%v:%v", "", *pPtr + *idPtr)
	send_back_URL := fmt.Sprintf("%v:%v", *sPtr, *pPtr - 1 + *idPtr)

	if *idPtr == 0 {
		glava := fmt.Sprintf("%v:%v", *sPtr, *pPtr + 1)
		rep := fmt.Sprintf("%v:%v", *sPtr, *pPtr + *nPtr - 1)

		Initiator(glava, rep, read_from_URL, send_back_URL, "unknown")
		return
	}
	
	if *idPtr >= *nPtr - 1 {
		Proces(read_from_URL, "unknown", send_back_URL)
	} else {
		Proces(read_from_URL, send_to_URL, send_back_URL)
	}

}
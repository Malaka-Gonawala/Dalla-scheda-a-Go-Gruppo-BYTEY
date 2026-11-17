package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"go.bug.st/serial"
)

func uso() {
	fmt.Println("uso:")
	fmt.Println("  go run . porte")
	fmt.Println("  go run . ascolta COM5")
	fmt.Println("  go run . set COM5 APPROVED")
	fmt.Println("  go run . get COM5")
	os.Exit(1)
}

// apre la porta a 115200, oppure esce con un errore
func apri(nome string) serial.Port {
	porta, err := serial.Open(nome, &serial.Mode{BaudRate: 115200})
	if err != nil {
		fmt.Println("errore: non riesco ad aprire", nome, ":", err)
		os.Exit(1)
	}
	return porta
}

// legge una riga dalla scheda, ma non oltre la scadenza
func leggiRiga(porta serial.Port, scadenza time.Time) (string, error) {
	porta.SetReadTimeout(100 * time.Millisecond)
	var riga []byte
	buf := make([]byte, 1)
	for time.Now().Before(scadenza) {
		n, err := porta.Read(buf)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue // timeout breve, riprovo finché c'è tempo
		}
		if buf[0] == '\n' {
			return strings.TrimSpace(string(riga)), nil
		}
		riga = append(riga, buf[0])
	}
	return "", fmt.Errorf("la scheda non ha risposto entro 2 secondi")
}

// chiede GET_LED e restituisce la risposta, es. {"led":"APPROVED"}
func chiediStato(porta serial.Port) (string, error) {
	porta.ResetInputBuffer()
	_, err := porta.Write([]byte("{\"command\":\"GET_LED\"}\n"))
	if err != nil {
		return "", err
	}
	scadenza := time.Now().Add(2 * time.Second)
	for {
		riga, err := leggiRiga(porta, scadenza)
		if err != nil {
			return "", err
		}
		if strings.Contains(riga, "\"led\"") {
			return riga, nil
		}
	}
}

func porte() {
	lista, err := serial.GetPortsList()
	if err != nil {
		fmt.Println("errore:", err)
		os.Exit(1)
	}
	if len(lista) == 0 {
		fmt.Println("nessuna porta trovata")
	}
	for _, p := range lista {
		fmt.Println(p)
	}
}

func ascolta(nome string) {
	porta := apri(nome)
	defer porta.Close()
	fmt.Println("in ascolto su", nome, "(Ctrl+C per uscire)")
	scanner := bufio.NewScanner(porta)
	for scanner.Scan() {
		riga := scanner.Text()
		if strings.Contains(riga, "REQUEST") {
			fmt.Println("richiesta di turno")
		} else if strings.Contains(riga, "RELEASE") {
			fmt.Println("rilascio")
		}
	}
}

func set(nome, stato string) {
	// controllo PRIMA di aprire la porta: con BLU non si scrive niente
	if stato != "APPROVED" {
		fmt.Println("errore: stato non valido:", stato, "(ammesso: APPROVED)")
		os.Exit(1)
	}
	porta := apri(nome)
	defer porta.Close()
	comando := "{\"command\":\"SET_LED\",\"state\":\"" + stato + "\"}\n"
	if _, err := porta.Write([]byte(comando)); err != nil {
		fmt.Println("errore di scrittura:", err)
		os.Exit(1)
	}
	fmt.Println("inviato:", strings.TrimSpace(comando))

	time.Sleep(50 * time.Millisecond) // lascio alla scheda il tempo di eseguire SET_LED

	risposta, err := chiediStato(porta)
	if err != nil {
		fmt.Println("errore nella verifica:", err)
		os.Exit(1)
	}
	if !strings.Contains(risposta, stato) {
		fmt.Println("errore: la scheda dice", risposta, "invece di", stato)
		os.Exit(1)
	}
	fmt.Println("verificato:", risposta)
}

func get(nome string) {
	porta := apri(nome)
	defer porta.Close()
	risposta, err := chiediStato(porta)
	if err != nil {
		fmt.Println("errore:", err)
		os.Exit(1)
	}
	fmt.Println(risposta)
}

func main() {
	if len(os.Args) < 2 {
		uso()
	}
	switch os.Args[1] {
	case "porte":
		porte()
	case "ascolta":
		if len(os.Args) < 3 {
			uso()
		}
		ascolta(os.Args[2])
	case "set":
		if len(os.Args) < 4 {
			uso()
		}
		set(os.Args[2], os.Args[3])
	case "get":
		if len(os.Args) < 3 {
			uso()
		}
		get(os.Args[2])
	default:
		uso()
	}
}

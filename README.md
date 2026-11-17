# Dalla scheda a Go

Progetto di laboratorio di TPSIT del Gruppo BYTEY.

Una scheda STM32 NUCLEO-L476RG gestisce una richiesta di turno con il bottone e il LED. Un programma scritto in Go, sul PC, comunica con la scheda attraverso la porta seriale USB.

## Struttura del repository

- Firmware della scheda (Tappa 1): progetto STM32CubeIDE in C.
- `go-scheda/` (Tappa 2): programma Go che parla con la scheda.

## Tappa 1: il firmware della scheda

La scheda usa la USART2 (porta seriale virtuale dell'ST-Link) a 115200 baud, 8 bit, nessuna parità, 1 bit di stop.

Il LED LD2 ha tre stati:

| Stato      | LED                  |
|------------|----------------------|
| `OFF`      | spento               |
| `PENDING`  | lampeggiante         |
| `APPROVED` | acceso fisso         |

Il bottone B1:

- da `OFF` porta a `PENDING` e invia `{"event":"REQUEST"}`;
- da `PENDING` o `APPROVED` porta a `OFF` e invia `{"event":"RELEASE"}`.

### Protocollo

Ogni messaggio è una riga JSON terminata da `\n`.

| Direzione      | Messaggio                                     | Significato                         |
|----------------|-----------------------------------------------|-------------------------------------|
| scheda -> PC   | `{"event":"REQUEST"}`                         | richiesta di turno (bottone)        |
| scheda -> PC   | `{"event":"RELEASE"}`                         | rilascio del turno (bottone)        |
| PC -> scheda   | `{"command":"SET_LED","state":"APPROVED"}`    | approva il turno, LED fisso         |
| PC -> scheda   | `{"command":"GET_LED"}`                       | chiede lo stato del LED             |
| scheda -> PC   | `{"led":"OFF"}`, `{"led":"PENDING"}`, `{"led":"APPROVED"}` | risposta a `GET_LED`   |

## Tappa 2: il programma Go

### Requisiti

- Go installato
- Libreria `go.bug.st/serial` versione 1.6.2 (già indicata in `go.mod`)
- Scheda collegata via USB (nel nostro caso sulla porta `COM5`)

### Come si usa

Dalla cartella `go-scheda`:

```powershell
go run . porte
go run . ascolta COM5
go run . set COM5 APPROVED
go run . get COM5
```

| Comando                 | Cosa fa                                                                                   |
|-------------------------|-------------------------------------------------------------------------------------------|
| `porte`                 | elenca le porte seriali disponibili                                                       |
| `ascolta COM5`          | stampa «richiesta di turno» o «rilascio» ogni volta che si preme il bottone (Ctrl+C per uscire) |
| `set COM5 APPROVED`     | invia `SET_LED` e poi verifica con `GET_LED` che la scheda abbia cambiato stato           |
| `get COM5`              | stampa lo stato del LED; se la scheda non risponde entro 2 secondi stampa un errore       |

Con uno stato non valido, per esempio `go run . set COM5 BLU`, il programma stampa un errore e non scrive niente sulla scheda.

### Esempio

```
> go run . set COM5 APPROVED
inviato: {"command":"SET_LED","state":"APPROVED"}
verificato: {"led":"APPROVED"}

> go run . get COM5
{"led":"APPROVED"}
```

## Licenza

Questo progetto è distribuito con licenza MIT. Il testo completo è nel file [LICENSE](LICENSE).

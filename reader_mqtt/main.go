package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type Event struct {
	UUID        string `bigquery:"uuid"`
	Impianto    string `bigquery:"impianto"`
	Codice      int    `bigquery:"codice"`
	Descrizione string `bigquery:"descrizione"`
	Valore      string `bigquery:"valore"`
	Data        string `bigquery:"data"`
}

func main() {
	// Carico il certificato blueteam
	caCert, err := ioutil.ReadFile("blueteam.crt")
	if err != nil {
		fmt.Println("Non sono riuscito a leggere il file.", err)
		os.Exit(1)
	}

	// Carico il certificato dell'utente reader
	cert, err := tls.LoadX509KeyPair("reader.crt", "reader.key")
	if err != nil {
		fmt.Println("Errore durante il caricamento del certificato e della chiave.", err)
		os.Exit(1)
	}

	// Creo un pool per i certificati e metto al suo interno il certificato blueteam.
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Imposto la configurazione TLS utilizzando la pool precedentamente creata.
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		RootCAs:            caCertPool,
	}

	// Impostazioni per il Client MQTT
	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.34.154.53.250.nip.io:8883")
	opts.SetClientID("reader")
	opts.SetTLSConfig(tlsConfig)

	// Creo e avvio il Client MQTT
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Imposto le credenziali per l'autenticarmi su Google Cloud Platform.
	// Il seguente utente ha i permessi di editor per il servizio di BigQuery.
	ctx := context.Background()
	creds := []byte(`
	{
		// Compilami :(
	}
	  
	`)

	// Creo il Client per BigQuery passando le credenziali in formato JSON
	bqClient, err := bigquery.NewClient(ctx, "blueteam-vargroup", option.WithCredentialsJSON(creds))
	if err != nil {
		fmt.Println("Non è stato possibile creare il Client BigQuery.", err)
		os.Exit(1)
	}

	// Creo la mappa per memorizzare i messaggi
	messageMap := make(map[string]map[string]string)

	// Definisco una funziona per la gestione dei messaggi.
	messageHandler := func(client MQTT.Client, msg MQTT.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())

		// Estraggo i dati dall'argomento
		impianto, codice, _, _ := extractDataFromTopic(topic, payload)

		// Controllo se l'evento esiste
		if _, ok := messageMap[impianto]; !ok {
			messageMap[impianto] = make(map[string]string)
		}

		// Salvo i messaggi nella mappa
		messageMap[impianto][topic] = payload

		// Verifico se tutti i messaggi sono presenti
		if len(messageMap[impianto]) == 3 {
			codiceVal, hasCodiceVal := messageMap[impianto]["/eventi/"+impianto+"/codice"]
			descrizione, hasDescrizione := messageMap[impianto]["/eventi/"+impianto+"/descrizione"]
			valore, hasValore := messageMap[impianto]["/eventi/"+impianto+"/tipologia"]

			if !hasCodiceVal || codiceVal == "" {
				// Se codiceVal contiene una stringa vuota, imposto codice a -1 come valore di default.
				codice = -1
			} else {
				// Converto codiceVal da stringa a intero.
				// Se la conversione fallisce, stampo un errore e interrompo la funzione.
				var err error
				codice, err = strconv.Atoi(codiceVal)
				if err != nil {
					fmt.Println("Non sono riuscito a convertire la variabile codice in un intero.", err)
					return
				}
			}
			if !hasDescrizione {
				descrizione = ""
			}
			if !hasValore {
				valore = ""
			}

			// Cancello
			delete(messageMap, impianto)

			// Creo la struttura con i dati ricavati.
			eventData := Event{
				UUID:        uuid.NewString(),
				Impianto:    impianto,
				Codice:      codice,
				Descrizione: descrizione,
				Valore:      valore,
				Data:        time.Now().UTC().Format("2006-01-02 15:04:05"),
			}

			// Lancio la query
			if err := insertEventoIntoBigQuery(ctx, bqClient, eventData); err != nil {
				fmt.Println("Errore durante la query:", err)
			}
		}
	}

	// Mi iscrivo a tutti gli argomenti dentro eventi
	if token := client.Subscribe("/eventi/#", 0, messageHandler); token.Wait() && token.Error() != nil {
		fmt.Println("Errore durante la sottoscrizione al topic /eventi/#:", token.Error())
		os.Exit(1)
	}

	// Mantengo in esecuzione.
	select {}
}

func extractDataFromTopic(topic string, payload string) (string, int, string, string) {
	// Divido la stringa topic usando / come delimitatore
	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return "", 0, "", ""
	}

	// Prelevo il nome impianto
	impianto := parts[2]
	// Prelevo il topic
	topicType := parts[3]

	codice := 0
	descrizione := ""
	valore := ""

	if topicType == "codice" {
		// Converto in intero
		codice, _ = strconv.Atoi(payload)
	} else if topicType == "descrizione" {
		descrizione = payload
	} else if topicType == "tipologia" {
		valore = payload
	}

	return impianto, codice, descrizione, valore
}

func insertEventoIntoBigQuery(ctx context.Context, client *bigquery.Client, event Event) error {
	// Imposto il data-set e la tabella
	datasetID := "blueteam_mqtt"
	tableID := "eventi"

	// Imposto i riferimenti
	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)

	// Creo un oggetto inserter
	inserter := table.Inserter()

	// Inserisco i dati dentro BigQuery
	if err := inserter.Put(ctx, event); err != nil {
		return err
	}

	return nil
}

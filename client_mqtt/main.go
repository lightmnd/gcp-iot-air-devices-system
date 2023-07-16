package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ErrorMessage struct {
	Codice      int         `json:"codice"`
	Descrizione string      `json:"descrizione"`
	Tipologia   interface{} `json:"tipologia"`
}

func jsonGroup() []ErrorMessage {
	group1 := []ErrorMessage{
		{Codice: 1, Descrizione: "On / Off", Tipologia: "TRUE"},
		{Codice: 2, Descrizione: "On effettivo", Tipologia: "TRUE"},
		{Codice: 3, Descrizione: "Stagione", Tipologia: "TRUE"},
		{Codice: 4, Descrizione: "Stagione effettiva", Tipologia: "TRUE"},
		{Codice: 5, Descrizione: "Abilitazione fasce orarie su display macchina", Tipologia: "TRUE"},
		{Codice: 6, Descrizione: "Deumidifica attiva", Tipologia: "TRUE"},
		{Codice: 7, Descrizione: "Richiesta acqua per trattamento aria", Tipologia: "TRUE"},
		{Codice: 8, Descrizione: "Warning filtri sporchi", Tipologia: "TRUE"},
		{Codice: 11, Descrizione: "Abilitazione deumidifica", Tipologia: "TRUE"},
		{Codice: 12, Descrizione: "Abilitazione riscaldamento", Tipologia: "TRUE"},
		{Codice: 13, Descrizione: "Abilitazione raffreddamento", Tipologia: "TRUE"},
		{Codice: 14, Descrizione: "Abilitazione free-cooling\\heating", Tipologia: "TRUE"},
		{Codice: 15, Descrizione: "Presenza riscaldamento dell’aria", Tipologia: "TRUE"},
		{Codice: 16, Descrizione: "Presenza raffreddamento dell’aria", Tipologia: "TRUE"},
		{Codice: 17, Descrizione: "Presenza recuperatore", Tipologia: "TRUE"},
		{Codice: 18, Descrizione: "Presenza free-cooling\\heating", Tipologia: "TRUE"},
		{Codice: 19, Descrizione: "Riscaldamento attivo", Tipologia: "TRUE"},
		{Codice: 20, Descrizione: "Raffreddamento attivo", Tipologia: "TRUE"},
		{Codice: 21, Descrizione: "Ricambio attivo", Tipologia: "TRUE"},
		{Codice: 22, Descrizione: "Free-cooling\\heating attivo", Tipologia: "TRUE"},
		{Codice: 23, Descrizione: "Sbrinamento attivo", Tipologia: "TRUE"},
		{Codice: 24, Descrizione: "Abilitazione riduzione ventilazione", Tipologia: "TRUE"},
		{Codice: 25, Descrizione: "Abilitazione umidifica", Tipologia: "TRUE"},
		{Codice: 26, Descrizione: "Sbrinamento recuperatore attivo", Tipologia: "TRUE"},
		{Codice: 27, Descrizione: "Presenza condensatore remoto", Tipologia: "TRUE"},
		{Codice: 28, Descrizione: "Presenza valvola acqua", Tipologia: "TRUE"},
		{Codice: 29, Descrizione: "Presenza valvola acqua on-off", Tipologia: "TRUE"},
		{Codice: 30, Descrizione: "Presenza valvola acqua modulante", Tipologia: "TRUE"},
		{Codice: 31, Descrizione: "Presenza e abilitazione batteria acqua calda", Tipologia: "TRUE"},
		{Codice: 32, Descrizione: "Presenza e abilitazione batteria acqua fredda", Tipologia: "TRUE"},
		{Codice: 33, Descrizione: "Presenza controllo temperatura", Tipologia: "TRUE"},
		{Codice: 34, Descrizione: "On compressore 1", Tipologia: "TRUE"},
		{Codice: 35, Descrizione: "On compressore 2", Tipologia: "TRUE"},
		{Codice: 36, Descrizione: "On ventilatore mandata", Tipologia: "TRUE"},
		{Codice: 37, Descrizione: "On resistenza elettriche", Tipologia: "TRUE"},
		{Codice: 1, Descrizione: "Unità ON", Tipologia: "TRUE"},
		{Codice: 10, Descrizione: "Stato compressore", Tipologia: "TRUE"},
		{Codice: 11, Descrizione: "Stato valvola acqua", Tipologia: "TRUE"},
		{Codice: 12, Descrizione: "Stato resistenza elettrica", Tipologia: "TRUE"},
		{Codice: 13, Descrizione: "Presenza valvola acqua", Tipologia: "TRUE"},
		{Codice: 14, Descrizione: "Presenza resistenza elettrica", Tipologia: "TRUE"},
		{Codice: 15, Descrizione: "Presenza allarme", Tipologia: "TRUE"},
		{Codice: 16, Descrizione: "Filtri da pulire", Tipologia: "TRUE"},
		{Codice: 17, Descrizione: "Presenza ventilatori elettronici", Tipologia: "TRUE"},
		{Codice: 18, Descrizione: "Presenza opzione sbrinamento gas caldo", Tipologia: "TRUE"},
		{Codice: 27, Descrizione: "Sbrinamento attivo", Tipologia: "TRUE"},
		{Codice: 28, Descrizione: "Richiesta deumidifica", Tipologia: "TRUE"},
		{Codice: 29, Descrizione: "Richiesta riscaldamento", Tipologia: "TRUE"},
		{Codice: 30, Descrizione: "Richiesta raffreddamento", Tipologia: "TRUE"},
		{Codice: 1, Descrizione: "temperatura ambiente", Tipologia: "20"},
		{Codice: 2, Descrizione: "temperatura esterna", Tipologia: "34"},
		{Codice: 3, Descrizione: "umidità relativa ambiente", Tipologia: "2.2"},
		{Codice: 4, Descrizione: "set umidità relativa", Tipologia: "3.7"},
		{Codice: 5, Descrizione: "set umidità relativa effettiva", Tipologia: "1.9"},
		{Codice: 6, Descrizione: "set temperatura / set temperatura inverno", Tipologia: "5.4"},
		{Codice: 7, Descrizione: "set temperatura estate", Tipologia: "4.9"},
		{Codice: 8, Descrizione: "set temperatura effettivo", Tipologia: "3.1"},
		{Codice: 9, Descrizione: "Percentuale ventilatore mandata", Tipologia: "82"},
		{Codice: 10, Descrizione: "Percentuale ventilatore estrazione", Tipologia: "73"},
		{Codice: 11, Descrizione: "Percentuale valvola acqua", Tipologia: "41"},
		{Codice: 12, Descrizione: "Percentuale umidificatore", Tipologia: "65"},
		{Codice: 13, Descrizione: "Percentuale valvola gas", Tipologia: "28"},
		{Codice: 14, Descrizione: "Percentuale serranda free-cooling", Tipologia: "47"},
		{Codice: 15, Descrizione: "Percentuale serranda ricircolo", Tipologia: "91"},
		{Codice: 16, Descrizione: "Set sbrinamento statico", Tipologia: "1.8"},
		{Codice: 17, Descrizione: "Differenziale sbrinamento statico", Tipologia: "0.6"},
		{Codice: 18, Descrizione: "Tempo sgocciolamento sbrinamento statico", Tipologia: "12"},
		{Codice: 19, Descrizione: "Temperatura mandata in ambiente", Tipologia: "17"},
		{Codice: 20, Descrizione: "Versione software", Tipologia: "3.5"},
		{Codice: 21, Descrizione: "Percentuale ricambio, step fisso e minimo di 5%", Tipologia: "9"},
		{Codice: 22, Descrizione: "Percentuale ricambio effettivo", Tipologia: "14"},
		{Codice: 23, Descrizione: "Temperatura protezione batteria acqua", Tipologia: "42.7"},
		{Codice: 24, Descrizione: "Temperatura ingresso batteria acqua", Tipologia: "39.2"},
		{Codice: 25, Descrizione: "temperatura ambiente", Tipologia: "24.9"},
		{Codice: 1, Descrizione: "Forzatura umidifica", Tipologia: "37"},
		{Codice: 2, Descrizione: "Set temperatura", Tipologia: "21.5"},
		{Codice: 3, Descrizione: "Set umidità relativa", Tipologia: "58.3"},
		{Codice: 4, Descrizione: "Ventilatore di ricircolo in standby", Tipologia: "67"},
		{Codice: 5, Descrizione: "Differenziale on raffreddamento", Tipologia: "1.2"},
		{Codice: 6, Descrizione: "Differenziale off raffreddamento", Tipologia: "0.9"},
		{Codice: 7, Descrizione: "Differenziale on deumidifica", Tipologia: "1.6"},
		{Codice: 8, Descrizione: "Differenziale off deumidifica", Tipologia: "1.1"},
		{Codice: 9, Descrizione: "Differenziale on riscaldamento", Tipologia: "1.4"},
		{Codice: 10, Descrizione: "Differenziale off riscaldamento", Tipologia: "0.8"},
		{Codice: 11, Descrizione: "Inizio rampa umidifica", Tipologia: "1.3"},
		{Codice: 12, Descrizione: "Fine rampa umidifica", Tipologia: "2.7"},
		{Codice: 13, Descrizione: "Offset temperatura ambiente", Tipologia: "0.5"},
		{Codice: 14, Descrizione: "Offset umidità ambiente", Tipologia: "0.9"},
		{Codice: 15, Descrizione: "Ore di attesa promemoria pulizia filtri", Tipologia: "24"},
		{Codice: 16, Descrizione: "Taratura fase 1 - ventilatore mandata", Tipologia: "75"},
		{Codice: 17, Descrizione: "Taratura fase 2 - ventilatore mandata", Tipologia: "68"},
		{Codice: 18, Descrizione: "Taratura fase 2 - ventilatore estrazione", Tipologia: "54"},
		{Codice: 19, Descrizione: "Taratura fase 3 - ventilatore mandata", Tipologia: "79"},
		{Codice: 20, Descrizione: "Taratura fase 3 - serranda ricircolo", Tipologia: "92"},
		{Codice: 21, Descrizione: "Taratura fase 1 - ventilatore mandata", Tipologia: "51"},
		{Codice: 22, Descrizione: "Taratura fase 2 - ventilatore mandata", Tipologia: "62"},
		{Codice: 23, Descrizione: "Taratura fase 2 - ventilatore estrazione", Tipologia: "47"},
		{Codice: 24, Descrizione: "Taratura fase 3 - ventilatore mandata", Tipologia: "58"},
		{Codice: 25, Descrizione: "Taratura fase 3 - serranda ricircolo", Tipologia: "71"},
	}

	return group1
}

func main() {
	// Carico il certificato blueteam
	caCert, err := ioutil.ReadFile("blueteam.crt")
	if err != nil {
		fmt.Println("Non sono riuscito a leggere il file.", err)
		os.Exit(1)
	}

	// Carico il certificato dell'impianto
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
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
	opts.SetClientID("client1")
	opts.SetTLSConfig(tlsConfig)

	// Creo e avvio il Client MQTT
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	startTime := time.Now()
	duration := 10 * time.Second
	for time.Since(startTime) < duration {
		group := jsonGroup()
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(group))
		message := group[randomIndex]

		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte(fmt.Sprintf("%d", message.Codice)))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte(message.Descrizione))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte(fmt.Sprintf("%v", message.Tipologia)))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		time.Sleep(1 * time.Second)
	}

	// Disconnetto l'impianto
	client.Disconnect(250)
}

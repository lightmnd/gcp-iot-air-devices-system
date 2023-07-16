package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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

	token := client.Subscribe("/remotecontrol/client1/command-remote", 1, handleCommandRemote)
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Mi sono iscritto al topic: /remotecontrol/client1/command-remote")

	// Mantengo in esecuzione.
	select {}
}

// Funzione di callback per la gestione dei comandi ricevuti tramite topic
func handleCommandRemote(client MQTT.Client, msg MQTT.Message) {
	if string(msg.Payload()) == "forzatura-deumedifica" {
		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte("2"))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte("Forzatura deumidifica"))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte("TRUE"))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		fmt.Println("Forzatura deumidifica")
	}
	if string(msg.Payload()) == "halt" {
		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte("1"))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte("Off"))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte("TRUE"))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		fmt.Println("Ho ricevuto un comando di power-off")
	}
	if string(msg.Payload()) == "reset --allarmi" {
		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte("10"))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte("Reset allarmi"))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte("TRUE"))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		fmt.Println("Ho ricevuto un comando per effetturare il reset degli allarmi")
	}
	if string(msg.Payload()) == "start-allarme-generico" {
		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte("500"))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte("allarme generico"))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte("TRUE"))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		fmt.Println("ERRORE 500: ALLARME GENERICO")
	}
	if string(msg.Payload()) == "stop-allarme-generico" {
		tokenCodice := client.Publish("/eventi/client1/codice", 1, false, []byte("500"))
		tokenDescrizione := client.Publish("/eventi/client1/descrizione", 1, false, []byte("allarme generico"))
		tokenTipologia := client.Publish("/eventi/client1/tipologia", 1, false, []byte("FALSE"))

		tokenCodice.Wait()
		tokenDescrizione.Wait()
		tokenTipologia.Wait()
		fmt.Println("OPERATIVO")
	}
}

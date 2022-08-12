package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}

func run(args []string, stdout io.Writer) error {
	var (
		flags    = flag.NewFlagSet(args[0], flag.ExitOnError)
		host     = flags.String("host", "", "the host/endpoint")
		keyFile  = flags.String("keyfile", "", "The private key file in PEM format")
		certFile = flags.String("certfile", "", "The certificate file in PEM format")
		caFile   = flags.String("cafile", "", "The CA certificate file in PEM format")
		topic    = flags.String("topic", "topic_1", "The MQTT topic for pub/sub")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if *host == "" || *keyFile == "" || *certFile == "" {
		return errors.New("host, keyFile and certFile are required")
	}

	fmt.Printf("host=%s, keyFile=%s, certFile=%s, caFile=%s\n", *host, *keyFile, *certFile, *caFile)

	// MQTT specific stuff
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println("Connected")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Printf("Connect lost: %v", err)
	}

	tlsConfig, err := NewTLSConfig(*caFile, *certFile, *keyFile)
	if err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", *host, 8883))
	opts.SetClientID("sdk-java")
	opts.SetTLSConfig(tlsConfig)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// Start the connection.
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	publish(client, *topic)
	client.Disconnect(250)
	return nil
}

func publish(client mqtt.Client, topic string) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func NewTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certpool.AppendCertsFromPEM(pemCerts)

	// Import client certificate/key pair.
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	config := &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
	return config, nil
}

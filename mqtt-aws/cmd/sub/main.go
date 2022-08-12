package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic = "topic_1"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func main() {
	var (
		flags    = flag.NewFlagSet("mqtt", flag.ExitOnError)
		host     = flag.String("host", "", "the host/endpoint")
		keyFile  = flags.String("keyfile", "", "The private key file in PEM format")
		certFile = flags.String("certfile", "", "The certificate file in PEM format")
		caFile   = flags.String("cafile", "", "The CA certificate file in PEM format")
	)

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}

	if *keyFile == "" || *certFile == "" {
		log.Fatalln(errors.New("both keyFile and certFile are required"))
	}

	fmt.Printf("keyFile = %s, certFile = %s, caFile = %s\n", *keyFile, *certFile, *caFile)
	tlsConfig, err := NewTLSConfig(*caFile, *certFile, *keyFile)
	if err != nil {
		log.Fatalln(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", *host, 8883))
	opts.SetClientID("sdk-java")
	opts.SetTLSConfig(tlsConfig)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	sub(client)

	client.Disconnect(250)
}

func sub(client mqtt.Client) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
	time.Sleep(9 * time.Second)
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

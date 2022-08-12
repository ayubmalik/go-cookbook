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
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {

	var (
		flags    = flag.NewFlagSet(args[0], flag.ExitOnError)
		host     = flag.String("host", "", "the host/endpoint")
		keyFile  = flags.String("keyfile", "", "The private key file in PEM format")
		certFile = flags.String("certfile", "", "The certificate file in PEM format")
		caFile   = flags.String("cafile", "", "The CA certificate file in PEM format")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if *keyFile == "" || *certFile == "" {
		return errors.New("both keyFile and certFile are required")
	}

	fmt.Printf("keyFile = %s, certFile = %s, caFile = %s\n", *keyFile, *certFile, *caFile)

	// MQTT specific stuff
	tlsConfig, err := NewTLSConfig(*caFile, *certFile, *keyFile)
	if err != nil {
		return err
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", *host, 8883))
	opts.SetClientID("sdk-java")
	opts.SetTLSConfig(tlsConfig)

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		log.Printf("connecting to %s\n", broker.String())
		return tlsConfig
	})

	// Start the connection.
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Printf("connected to MQTT = %t\n", c.IsConnected())

	done := make(chan bool)
	go publish(done, c, "topic_1")
	<-done

	defer c.Disconnect(50)

	return nil
}

func publish(done chan bool, conn mqtt.Client, topic string) {
	tmpl, err := template.New("msg").Parse(`{ "timestamp": "{{.}}", "frequency": "50.0"}`)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 100; i++ {
		msg := &strings.Builder{}
		err = tmpl.Execute(msg, time.Now().UTC())
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("sending message: %s", msg)
		if token := conn.Publish(topic, 0, false, msg.String()); token.Wait() && token.Error() != nil {
			log.Fatalf("failed to send update: %v", token.Error())
		}
		time.Sleep(time.Second)
	}
	done <- true
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

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	print      = fmt.Println
	configFlag = flag.String("c", "", "-c Path to config file.")
	isTls      = false
)

func main() {

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = checkConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(config.PubDirectory)))

	if isTls {
		go log.Fatal(http.ListenAndServeTLS(":"+config.ListenPort, config.CertFullChainPath, config.CertPrivateKeyPath, nil))
	} else {
		go log.Fatal(http.ListenAndServe(":"+config.ListenPort, nil))
	}
}

// Boilerplate stuffs below
func checkFlags() {
	flag.Parse()
	if *configFlag == "" {
		print("[ERROR] Must provide a config file as an argument")
		printUsage()
		os.Exit(0)
	}
}

func printUsage() {
	print("TODO")
}

func getConfig() (Config, error) {
	var config Config

	checkFlags()

	bytez, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytez, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func checkConfig(config Config) error {
	if config.ListenPort == "" {
		return errors.New("[ERROR] Must specify a port to listen on")
	}

	if config.CertFullChainPath != "" {
		_, err := os.Stat(config.CertFullChainPath)
		if err != nil {
			return errors.New("[ERROR] Cert chain path invalid")
		}
	}

	if config.CertPrivateKeyPath != "" {
		_, err := os.Stat(config.CertPrivateKeyPath)
		if err != nil {
			return errors.New("[ERROR] Cert private key path invalid")
		}
	}

	if config.ListenPort == "443" && (config.CertFullChainPath == "" || config.CertPrivateKeyPath == "") {
		return errors.New("[ERROR] Provided port 443 but no certificate!")
	}

	if config.CertFullChainPath != "" && config.CertPrivateKeyPath != "" {
		isTls = true
	}

	return nil
}

type Config struct {
	ListenPort         string `json:"listenPort"`
	PubDirectory       string `json:"pubDirectory"`
	CertFullChainPath  string `json:"certFullChainPath"`
	CertPrivateKeyPath string `json:"certPrivateKeyPath"`
}

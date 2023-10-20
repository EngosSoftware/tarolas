// Package main is the storage server's starting point.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/wisbery/tarolas/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var defaultConfigurationFileName = "./config/config.json"

// showUsageAndExit shows the usage of this service.
func showUsage() {
}

// readConfiguration reads server's configuration from configuration file specified in command line.
// To correctly start tarolas server, there must be given one command line parameter named '--config'.
// This parameter must specify existing JSON file with configuration parameters.
// When something goes wrong, the error message and usage is printed,
// and then the file server terminates with exit code 1.
func readConfiguration() (*Configuration, error) {
	// prepare configuration flag
	cfgFileName := flag.String("config", "", "path to configuration file")
	flag.Parse()
	// check if configuration file name is given
	if *cfgFileName == "" {
		cfgFileName = &defaultConfigurationFileName
	}
	// try to read the configuration file
	data, err := os.ReadFile(*cfgFileName)
	if err != nil {
		fmt.Printf("error occured while reading tarolas configuration file\n%v\n\n", err)
		showUsage()
		return nil, err
	}
	// try to load JSON data read from configuration file
	var cfg Configuration
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("error occured while parsing tarolas configuration file:\n%v\n\n", err)
		showUsage()
		return nil, err
	}
	// TODO check here if root directory exists and the server has read/write/delete access in this directory
	// TODO add read/write/delete checks
	return &cfg, nil
}

// Function main starts the server.
func main() {
	var httpServer *http.Server
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	if configuration, err := readConfiguration(); err == nil {
		httpServer = StartServer(configuration)
		sig := <-osSignals
		fmt.Printf("\nTarolas received signal [%v], gracefully exiting...\n", sig)
		StopServer(httpServer)
	}
}

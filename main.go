package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

const (
	stateDir  = ".piaport"
	stateFile = "state"
)

var (
	piaUsername = flag.String("pia.username", "", "PIA username")
	piaPassword = flag.String("pia.password", "", "PIA password")

	transmissionUsername = flag.String("transmission.username", "", "Transmission RPC username")
	transmissionPassword = flag.String("transmission.password", "", "Transmission RPC password")
	transmissionHost     = flag.String("transmission.host", "127.0.0.1", "Transmission RPC host")

	ifaceName = flag.String("vpn.iface", "tun0", "VPN network interface")
)

type UpdateTarget interface {
	UpdateOpenPort(oldPort, newPort int) error
	IsUpdateNecessary() (bool, error)
}

type PiaState struct {
	OldPort int `json:"old_port"`
}

func main() {
	flag.Parse()
	piaClientId := RandString(32)
	userDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Can't determine homedir: %#v", err)
	}
	oldPort, err := getOldPort(userDir)
	if err != nil {
		log.Fatalf("Failed to get old port: %#v", err)
	}

	transmissionUpdater, err := NewTransmissionClient(*transmissionHost, *transmissionUsername, *transmissionPassword)
	if err != nil {
		log.Fatalf("Failed to create transmission RPC client: %#v", err)
	}

	if update, err := transmissionUpdater.IsUpdateNecessary(); err != nil {
		log.Fatalf("Can't determine if Transmission needs an update: %#v", err)
	} else if !update {
		log.Println("Transmission peer port still seems to be open, aborting")
		return
	}
	log.Printf("Transmission port %d is closed, update needed\n", oldPort)
	if oldPort != -1 {
		log.Printf("Deleting rules for port %d", oldPort)
		if err := DeleteOldRules(oldPort); err != nil {
			log.Printf("[WARN] Failed to delete old firewall rules for port %d: %#v", oldPort, err)
		}
	}

	newPort, err := RequestOpenPort(*ifaceName, PIACredentials{
		ClientID: piaClientId,
		Username: *piaUsername,
		Password: *piaPassword,
	})
	if err != nil {
		log.Fatalf("Failed to request new open port from PIA: %s %#v", err, err)
	}
	log.Printf("Got port %d from PIA to forward\n", newPort)

	if err := saveOldPort(userDir, newPort); err != nil {
		log.Printf("Failed to save requested port %d as old port. Take care to remove old firewall rules manually. %#v", newPort, err)
	}

	log.Printf("Creating iptable rules for port %d\n", newPort)
	if err := CreateNewRule(newPort); err != nil {
		log.Fatalf("Failed to create new firewall rules for open port %d: %s %#v", newPort, err, err)
	}

	log.Printf("Updating transmission to use port %d instead of old port %d\n", newPort, oldPort)
	if err := transmissionUpdater.UpdateOpenPort(oldPort, newPort); err != nil {
		log.Printf("Failed to update port on transmission: %s %#v", err, err)
	}

}

func saveOldPort(userDir string, port int) error {
	stateFileDir := path.Join(userDir, stateDir)
	os.MkdirAll(stateFileDir, 0744)
	stateFilePath := path.Join(stateFileDir, stateFile)
	state := PiaState{
		OldPort: port,
	}

	fileBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(stateFilePath, fileBytes, 0644)
}

func getOldPort(userDir string) (int, error) {
	stateFileDir := path.Join(userDir, stateDir)
	os.MkdirAll(stateFileDir, 0744)
	stateFilePath := path.Join(stateFileDir, stateFile)
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return -1, nil
	}

	fileBytes, err := ioutil.ReadFile(stateFilePath)
	if err != nil {
		return -1, err
	}

	var state PiaState

	if err := json.Unmarshal(fileBytes, &state); err != nil {
		return -1, err
	}

	return state.OldPort, nil
}

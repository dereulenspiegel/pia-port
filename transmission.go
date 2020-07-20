package main

import (
	trans "github.com/hekmon/transmissionrpc"
)

type TransmissionUpdateClient struct {
	client *trans.Client
}

func NewTransmissionClient(host, username, password string) (*TransmissionUpdateClient, error) {
	client, err := trans.New(host, username, password, nil)
	if err != nil {
		return nil, err
	}
	return &TransmissionUpdateClient{
		client: client,
	}, nil
}

func (t *TransmissionUpdateClient) UpdateOpenPort(oldPort, port int) error {
	peerPort := int64(port)
	return t.client.SessionArgumentsSet(&trans.SessionArguments{
		PeerPort: &peerPort,
	})
}

func (t *TransmissionUpdateClient) IsUpdateNecessary() (bool, error) {
	open, err := t.client.CheckPort()
	return !open, err
}

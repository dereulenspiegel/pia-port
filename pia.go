package main

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	piaEndpoint = "https://www.privateinternetaccess.com/vpninfo/port_forward_assignment"
)

type PIACredentials struct {
	Username string
	Password string
	ClientID string
}

type piaPortResponse struct {
	Port int `json:"port"`
}

func RequestOpenPort(ifaceName string, creds PIACredentials) (int, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return -1, err
	}

	addr, err := selectNonLocalV4Addr(iface)
	if err != nil {
		return -1, err
	}

	return requestPortFromPia(creds, addr)
}

func requestPortFromPia(creds PIACredentials, localAddr *net.IP) (int, error) {
	postData := url.Values{
		"user":      []string{creds.Username},
		"pass":      []string{creds.Password},
		"client_id": []string{creds.ClientID},
		"local_ip":  []string{localAddr.String()},
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.PostForm(piaEndpoint, postData)
	if err != nil {
		return -1, err
	}

	var piaResp piaPortResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&piaResp); err != nil {
		return -1, err
	}
	return piaResp.Port, nil
}

func selectNonLocalV4Addr(iface *net.Interface) (*net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				// This is most probably an IPv4 address
				return &ipnet.IP, nil
			}
		} else {
			return nil, errors.New("Unknown Addr type on interface")
		}
	}
	return nil, errors.New("No suitable address found on interface")
}

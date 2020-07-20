package main

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterfaceAddr(t *testing.T) {
	iface, err := net.InterfaceByName("en0")
	require.Nil(t, err)
	addrs, err := iface.Addrs()
	require.Nil(t, err)

	for _, addr := range addrs {
		switch t := addr.(type) {
		case *net.IPAddr:
			fmt.Printf("Got an net.IPAddr: %#v \n", t)
		default:
			fmt.Printf("Got something different: %T\n", t)
		}
		fmt.Printf("%s \n", addr)
	}
}

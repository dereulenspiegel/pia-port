package main

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

var (
	defaultTableName = "filter"
	defaultChainName = "INPUT"
)

func ruleSpec(proto string, port int) []string {
	return []string{"-p", proto, "--dport", fmt.Sprintf("%d", port), "-j", "ACCEPT"}
}

func DeleteOldRules(oldPort int) error {
	table, err := iptables.New()
	if err != nil {
		return err
	}

	err = table.Delete(defaultTableName, defaultChainName, ruleSpec("tcp", oldPort)...)
	err = table.Delete(defaultTableName, defaultChainName, ruleSpec("udp", oldPort)...)
	return err
}

func CreateNewRule(newPort int) error {
	table, err := iptables.New()
	if err != nil {
		return err
	}

	err = table.Insert(defaultTableName, defaultChainName, 1, ruleSpec("tcp", newPort)...)
	err = table.Insert(defaultTableName, defaultChainName, 2, ruleSpec("udp", newPort)...)
	return err
}

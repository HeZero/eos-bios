package bios

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	yaml2json "github.com/bronze1man/go-yaml2json"
	"github.com/eoscanada/eos-bios/bios/disco"
)

func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}

func LoadDiscoveryFromFile(fileName string) (discovery *disco.Discovery, err error) {
	cnt, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	err = yamlUnmarshal(cnt, &discovery)
	if err != nil {
		return
	}

	err = ValidateDiscovery(discovery)
	return
}

func ValidateDiscoveryFile(filename string) error {
	discovery, err := LoadDiscoveryFromFile(filename)
	if err != nil {
		return err
	}
	ValidateDiscovery(discovery)
	return nil
}

func ValidateDiscovery(discovery *disco.Discovery) error {
	for _, peer := range discovery.SeedNetworkPeers {
		if peer.Weight > 100 || peer.Weight < 0 {
			return fmt.Errorf("peer %q weight must be between 0 and 100", peer.Account)
		}
	}

	if discovery.TargetAccountName == "" {
		return errors.New("target_account_name is missing")
	}

	if !strings.Contains(discovery.TargetP2PAddress, ":") {
		return errors.New("target_p2p_address doesn't contain port number")
	}

	// if strings.Contains(discovery.TargetP2PAddress, "example.com") {
	// 	return errors.New("target_p2p_address contains an example.com domain, are you sure about that?")
	// }

	if strings.Contains(discovery.TargetP2PAddress, "://") {
		return fmt.Errorf("target_p2p_address should be of format ip:port, not prefixed with a protocol")
	}

	if !strings.Contains(discovery.TargetHTTPAddress, "://") {
		return fmt.Errorf("target_http_address should include the protocol (like http:// or https://)")
	}

	// TODO: make sure it's not the DEFAULT key

	// TODO: boot node should orchestrate the PUBLICLY accessible thing..
	//       so we don't need to have two `target_api_address`

	return nil
}
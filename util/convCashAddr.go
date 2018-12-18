package util

import (
	"strings"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
)

func ConvToCashAddr(addresses []string, param *chaincfg.Params) ([]string, error) {
	for idx, addr := range addresses {
		// to guarantee all addresses is valid
		cashAddr, err := cashutil.DecodeAddress(addr, param)
		if err != nil {
			return nil, err
		}

		addresses[idx] = cashAddr.EncodeAddress(true)
	}

	return addresses, nil
}

func ConvToCashAddrCopy(addresses []string, param *chaincfg.Params) ([]string, error) {
	ret := make([]string, len(addresses))
	for idx, addr := range addresses {
		if !IsCashAddrFormat(addr) {
			cashAddr, err := cashutil.DecodeAddress(addr, param)
			if err != nil {
				return nil, err
			}

			ret[idx] = cashAddr.EncodeAddress(true)
		} else {
			ret[idx] = addr
		}
	}

	return ret, nil
}

func IsCashAddrFormat(address string) bool {
	// simple judge for bitcoin cash symbol(`:`),
	// but this symbol is optional.
	if strings.Contains(address, ":") {
		return true
	}

	// 1. mainnet:  bitcoincash: + 42 bytes = 12 + 42
	// 2. testnet3: bchtest:     + 42 bytes = 8  + 42
	// 3. regtest:  bchreg:      + 42 bytes = 7  + 42
	// if the address does not have the prefix: length = 42
	// if the address have the prefix: max length = 42 + 12
	if len(address) == 42 {
		return true
	}

	return false
}

func ConvTolegacyAddr(addr string, param *chaincfg.Params) (string, error) {
	address, err := cashutil.DecodeAddress(addr, param)
	if err != nil {
		return "", err
	}

	return address.EncodeAddress(false), nil
}

func MustConvToLegacyAddr(addr string, param *chaincfg.Params) string {
	address, _ := cashutil.DecodeAddress(addr, param)

	return address.EncodeAddress(false)
}

func ConvSingleCashAddr(addr string, param *chaincfg.Params) (string, error) {
	address, err := cashutil.DecodeAddress(addr, param)
	if err != nil {
		return "", err
	}

	return address.EncodeAddress(true), nil
}

func MustConvSingleCashAddr(addr string, param *chaincfg.Params) string {
	address, _ := cashutil.DecodeAddress(addr, param)

	return address.EncodeAddress(true)
}

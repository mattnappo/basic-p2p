package p2p

import (
	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

// AddressList is an array of multiaddrs (for the flag parser).
type AddressList []maddr.Multiaddr

// Set sets the value of a mutliaddress in the address list.
func (al *AddressList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

// StringsToAddrs converts an array of strings to an AddressList.
func StringsToAddrs(addrStrings []string) (AddressList, error) {
	var maddrs AddressList
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return maddrs, nil
}

// String stringifys the address list.
func (al *AddressList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

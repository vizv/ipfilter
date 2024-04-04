package iprange

import (
	"math/big"
	"net/netip"
)

type IP struct {
	netip.Addr
}

func ParseIP(str string) (*IP, error) {
	ip, err := netip.ParseAddr(str)
	if err != nil {
		return nil, err
	}

	return &IP{ip}, nil
}

func (i *IP) ToInt() *Int {
	ret := big.NewInt(0)
	// fmt.Println(i.AsSlice())
	ret.SetBytes(i.AsSlice())

	return &Int{ret}
}

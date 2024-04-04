package iprange

import (
	"math/big"
	"net/netip"
)

type Int struct {
	*big.Int
}

func (i *Int) ToIP(isIPv6 bool) *IP {
	bufLen := 4
	if isIPv6 {
		bufLen = 16
	}

	buf := make([]byte, bufLen)
	bytes := i.Bytes()
	if len(bytes) > bufLen {
		return nil
	}

	// fmt.Println(bytes)
	copy(buf[bufLen-len(bytes):], bytes)
	// fmt.Println(buf)

	addr, ok := netip.AddrFromSlice(buf)
	if !ok {
		return nil
	}

	return &IP{addr}
}

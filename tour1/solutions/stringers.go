//go:build exclude

package solutions

import "fmt"

type IPAddr [4]byte

func (ia IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ia[0], ia[1], ia[2], ia[3])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}

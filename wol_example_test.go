package wol_test

import "github.com/mkch/wol"

func ExampleWake() {
	wol.Wake("11:22:33:44:55:66") // MAC address
}

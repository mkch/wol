/*
Package wol implements the creation and sending of Wake-on-LAN (WoL or WOL) magic packet.

Wake-on-LAN (WoL or WOL) is a standard that allows a computer to be turned on
or awakened by a network message(aka. magic packet).
See: https://en.wikipedia.org/wiki/Wake-on-LAN

A WoL message is usually broadcasted to all attached devices on a given network, using the network broadcast address.
It is typically sent as a UDP datagram to port 0 (reserved port number), 7 (Echo Protocol) or 9 (Discard Protocol),
or directly over Ethernet as EtherType 0x0842. A connection-oriented transport-layer protocol like TCP is less suited
for this task as it requires establishing an active connection before sending user data.

The simplest way to ues this package is

	wol.Wake("MAC address")
*/
package wol

import "net"

// NewPacket creates a WOL magic packet whose target is macAddr
// and with an optional SecureOn password.
func NewPacket(macAddr net.HardwareAddr, password *[6]byte) (packet []byte) {
	// https://en.wikipedia.org/wiki/Wake-on-LAN
	// 6 bytes of all 255 (FF FF FF FF FF FF in hexadecimal)
	packet = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	for i := 0; i < 16; i++ {
		packet = append(packet, macAddr...)
	}
	if password != nil {
		packet = append(packet, password[:]...)
	}
	return
}

// SendUDP sends a WOL magic packet to addr through UDP network.
// macAddr and password are passed to NewPacket to create a WOL
// magic packet.
// A WOL packet is typically sent as a UDP datagram to port 0 (reserved port number),
// 7 (Echo Protocol) or 9 (Discard Protocol)
func SendUDP(addr string, macAddr net.HardwareAddr, password *[6]byte) error {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(NewPacket(macAddr, password))
	return err
}

// Wake sends a WOL magic packet to 255.255.255.255:6 through UDP network.
// macAddr is the MAC address of the destination device.
func Wake(macAddress string) error {
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return err
	}
	return SendUDP(net.IPv4bcast.String()+":6", mac, nil)
}

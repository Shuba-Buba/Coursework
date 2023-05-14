package connectors

import (
	"log"
	"net"
)

func MakeUDPConnector(IP string, port uint) *net.UDPConn {
	addr := net.UDPAddr{
		Port: int(port),
		IP:   net.ParseIP(IP),
	}
	UDPconn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Panicf("Error creating UDP connector %v\n", err)
	}
	return UDPconn
}

package postman

import (
	"encoding/binary"
	"fmt"
	"net"
	"test/connectors"
	"test/contracts"
)

var start_port = 10000

type Postman struct {
	ch       <-chan contracts.Contract
	connects []*connectors.Connector
}

func MakePostman(receiver chan contracts.Contract) *Postman {
	return &Postman{ch: receiver}
}

func (this *Postman) Run() {
	cur_free_port := 0
	for {
		select {
		case new_contract := <-this.ch:
			if new_contract.Remote_port == 0 {
				cur_real_port := start_port + cur_free_port
				this.connects = append(this.connects, &connectors.Connector{
					Port: cur_real_port},
				)
				socket_address := "wss://stream.binance.com:9443/ws/" + new_contract.Symbol + "@depth@100ms"
				go this.connects[cur_free_port].Connect(socket_address)

				addr := fmt.Sprintf("127.0.0.1:%d", new_contract.Port)
				conn, _ := net.Dial("udp", addr)

				bs := make([]byte, 4)
				binary.LittleEndian.PutUint32(bs, uint32(cur_real_port))
				conn.Write(bs)
				cur_free_port += 1
			}
		default:
			// TODO
		}
	}

}

package postman

import (
	"encoding/json"
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

	for new_contract := range this.ch {

		// new port
		new_contract.Remote_port = start_port + cur_free_port

		current_connector := connectors.MakeConnector(new_contract)
		this.connects = append(this.connects, current_connector)

		go current_connector.Connect()

		addr := fmt.Sprintf("127.0.0.1:%d", new_contract.Port)
		conn, _ := net.Dial("udp", addr)

		bs, err := json.Marshal(*current_connector)

		if err != nil {
			panic("Bad try")
		}
		conn.Write(bs)
		cur_free_port += 1
	}

}

package contracts

type Contract struct {
	Symbol      string `json:"symbol"`
	Port        int
	Remote_port int
}

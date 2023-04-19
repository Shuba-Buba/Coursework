package contracts

type Contract struct {
	Symbol      string `json:"symbol"`
	FileName    string `json:"file_name"`
	Port        int
	Remote_port int
}

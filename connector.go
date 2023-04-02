

type Connector struct {

	chan conn; // Поле в котором будет хранится websocket соединение
	Orderbook* ; // Orderbook к которому он будет применять обновления
}

func (con *Connector) Connect() string {	
}

// обработчик сообщения из websocket канала
func callback () {}

// обработчик ошибки(закрывает канал)
func error callback () {}



/*
Примерное использование:

// создаём Connector
var conn Connector;

// стартуем коннектор и указываем Orderbook, к которому будем применять websocket сообщения
// внутри будет посылать запрос на connection а потом 
conn.Connect("BTCUSDT@futures@100ms", Orderbook* book)

// ждем 10 секунд

// закрываем текущее соединение
conn.Close()

print(book.asks[0], book.bids[0])

*/
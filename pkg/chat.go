package pkg

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// Параметры обновления http соединения в WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Массив для хранения Client
var clients = make(map[*Client]bool)

// Канал для рассылки сообщений всем клиентам
var broadcast = make(chan []byte)

// Client Структура для хранения
type Client struct {
	conn *websocket.Conn // Создание WebSocket
	send chan []byte     // Канал для передачи сообщения Client
}

// WebSocketHandler Новый Handler с ServeHTTP для правильной обработки WebSocket соединения
type WebSocketHandler struct{}

func (ws WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleWebSocket(w, r)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Перевод http соединения в WebSocket соединение
	fmt.Println("Client connected")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}

	// Ссылаемся на существующую структуру при инициации переменной client
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	clients[client] = true

	// Запуск go-routine для чтения и записи сообщений
	go client.write()
	go client.read()
}

// Метод для чтения сообщений
func (cl *Client) read() {
	// Закрытие соединения и удаление Client из массива после завершения метода
	defer func() {
		cl.conn.Close()
		delete(clients, cl)
	}()

	// Чтение сообщения от клиента
	for {
		_, message, err := cl.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Reading Message Error:", err)
			} // В случае неожиданного закрытия  WebSocket-подключения выдаст ошибку
			return // Если нет ошибки, выходим из метода
		} // Передаёт в переменную message сообщение client отправленное через WebSocket

		// Отправка сообщения в канал рассылки
		broadcast <- message
	}
}

// Метод для записи сообщений Client
func (cl *Client) write() {
	defer func() {
		cl.conn.Close()
		delete(clients, cl)
	}()

	// Отправка сообщения Client в канал через WebSocket-соединение
	for message := range cl.send {
		err := cl.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

// go-routine для рассылки сообщений всем подключенным Client
func handleMessages() {
	for {
		// Получаем сообщение из канала broadcast и записывает в message
		message := <-broadcast
		// Рассылаем сообщение всем подключенным Client
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			} // Если есть message, то будет отправлено сообщение в канал client.send
			// Если message нет, то канал закроется и client будет удалён из массива
		}
	}
}

// Запуск go-routine для отправки сообщений
func init() {
	go handleMessages()
}

package websocket

import (
	"encoding/json"
	"net/http"

	"livetail/app/consumer/kafkaConsumer"
	log "livetail/app/utility/logger"

	// "github.com/Shopify/sarama"
	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
)

type ClientRequest struct {
	StreamId string `json:"streamId"`
}

type Client struct {
	req        *http.Request
	connection *websocket.Conn
	// manager    *Manager
	egress     chan interface{}
	topicChan  chan string
	switchChan chan bool
}

func NewClient(conn *websocket.Conn, req *http.Request) *Client {

	return &Client{
		req:        req,
		connection: conn,
		// manager:    manager,
		// egress:     make(chan model.Message),
		egress:     make(chan interface{}),
		topicChan:  make(chan string),
		switchChan: make(chan bool),
	}
}

func (c *Client) removeClient(client *Client) {

	// Redis store remove inactive user
	// profileId := client.req.URL.Query().Get("id")
	// err := redis.PopUserFromActiveSet(profileId)
	// if err != nil {
	// 	log.ErrorLogger.Println(err)
	// }

	// remove from manager list and close the connection from server
	log.WarningLogger.Println("closing connection")
	// close websocket connection to the client
	// if client.connection.
	err := client.connection.Close()
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}

	// close(client.egress)
	// close(client.switchChan)
	// close(client.topicChan)
	// return

}
func (c *Client) ProcessMessage(mt int, p []byte) {
	log.DebugLogger.Printf("got message type: %v, with message: %v", mt, string(p))

	if mt == websocket.CloseMessage {
		log.DebugLogger.Println(mt, "got close connection message from client")
		return
	}

	// TODO
	if mt == websocket.BinaryMessage {
		// parse and execute
		log.WarningLogger.Println("got binary message on websocket. not implemented yet")
		return
		// var clientReq ClientRequest
		// err := json.Unmarshal(p, &clientReq)
		// if err != nil {
		// 	log.ErrorLogger.Println(err)
		// 	return
		// }
		// log.DebugLogger.Println(clientReq)
		// topic := clientReq.StreamId
		// c.topicChan <- topic
	}

	if mt == websocket.TextMessage {

		// data := &Incoming{}
		log.InfoLogger.Println(string(p))
		data := Incoming{}
		err := json.Unmarshal(p, &data)
		// err = c.connection.ReadJSON(data)
		if err != nil {
			log.ErrorLogger.Println(err)
			return
		}
		log.InfoLogger.Println(data)
		if data.Event == "switchstream" {
			log.InfoLogger.Println("switching stream")
			topic := ""
			if data.StreamId == c.req.URL.Query().Get("id") {
				topic = "user_topic_" + data.StreamId
				log.InfoLogger.Println(topic)
			} else {
				topic = "filter_topic_" + data.StreamId
				log.InfoLogger.Println(topic)
			}
			c.switchChan <- true
			c.topicChan <- topic
		} else {
			// TODO
			log.WarningLogger.Println("not implemented yet. the event type: ", data.Event)
			return
		}
	}
}

type Incoming struct {
	Event    string `json:"event"`
	StreamId string `json:"streamId"`
}

func (c *Client) ReadMessage() {
	// to close gracefully if this function exits
	defer func() {
		close(c.switchChan)
		// c.removeClient(c)
	}()

	for {
		messageType, p, err := c.connection.ReadMessage()
		log.DebugLogger.Println(messageType)
		if err != nil {
			log.WarningLogger.Println(err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// log.WarningLogger.Printf("error: %v", err)
				// return
				log.WarningLogger.Printf("error: %v, user-agent: %v", err, c.req.Header.Get("User-Agent"))
				return
			}
			// else {
			log.WarningLogger.Println("closing websocket connection due to above error")
			// 	// log.WarningLogger.Printf("error: %v", err)
			return
			// 	log.ErrorLogger.Printf("error: %v, user-agent: %v", err, c.req.Header.Get("User-Agent"))
			// 	// return
			// }
		}
		c.ProcessMessage(messageType, p)
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// ok will be false if egress channel is closed
			if !ok {
				// tell frontend that manager has closed this channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					// log the connection is closed
					log.InfoLogger.Println("error in writing-close message: ", err)
					return
				}
				return
			}
			//  write the message to the connection
			if err := c.connection.WriteJSON(message); err != nil {
				log.DebugLogger.Println("error in writing to the websocket due to: ", err)
				return
			}
		}
	}
}

func (c *Client) StreamKafka() {

	defer func() {
		// c.removeClient(c)
		close(c.egress)
	}()

	clientID := "go-kafka-consumer"

	// Create new consumer
	master, err := sarama.NewConsumer(kafkaConsumer.Brokers, kafkaConsumer.NewKafkaConfiguration(clientID))
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}

	defer func() {
		if err := master.Close(); err != nil {
			log.ErrorLogger.Println(err)
			return
		}
	}()
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)

	defer close(consumers)
	defer close(errors)

	doneThisChan := make(chan bool)
	defaultTopic := "user_topic_" + c.req.URL.Query().Get("id")
	// startSaramaConsume(c, &defaultTopic, master, c.switchChan, consumers, errors)
	startSaramaConsume(c, &defaultTopic, master, doneThisChan, consumers, errors)

	for {

		topic, ok := <-c.topicChan
		if !ok {
			log.DebugLogger.Println("2227")
			return
		}
		// log.InfoLogger.Println(topic)
		log.DebugLogger.Println("2228")
		doneThisChan := make(chan bool)
		// startSaramaConsume(c, &topic, master, c.switchChan, consumers, errors)
		startSaramaConsume(c, &topic, master, doneThisChan, consumers, errors)
		log.DebugLogger.Println("2222")
	}
}

func startSaramaConsume(
	c *Client,
	topic *string,
	master sarama.Consumer,
	doneChan chan bool,
	consumeChan chan *sarama.ConsumerMessage,
	errorChan chan *sarama.ConsumerError,
) {
	go func() {
		for {
			select {
			case msg, ok := <-consumeChan:
				if !ok {
					log.DebugLogger.Println("Sarama consumer Channel down")
					return
				}

				var message interface{}
				err := json.Unmarshal(msg.Value, &message)
				if err != nil {
					log.ErrorLogger.Println(err)
					return
				}
				log.DebugLogger.Println("writing to egress chan")
				c.egress <- message

			case consumerError, ok := <-errorChan:
				if !ok {
					log.DebugLogger.Println("Sarama error Channel down")
					return
				}

				log.ErrorLogger.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				log.DebugLogger.Println("2225")
				return
			case <-c.switchChan:
				log.DebugLogger.Println("2235")
				doneChan <- true
				return

			}
		}
	}()
	go kafkaConsumer.Consume(topic, master, doneChan, consumeChan, errorChan)
}

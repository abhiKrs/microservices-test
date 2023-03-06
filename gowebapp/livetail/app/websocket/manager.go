package websocket

import (
	// "bytes"
	"encoding/json"
	// "io"
	log "livetail/app/utility/logger"
	"net/http"
	// "livetail/app/utility/redis"
	// "github.com/Shopify/sarama"
)

func ServeClient(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("New Kafka Connection")

	// var req *Filter
	// // err := json.NewDecoder(r.Body).Decode(r)
	// err := req.Bind(r.Body)
	// if err != nil {
	// 	log.DebugLogger.Println(err)
	// 	res := Response{IsSuccessful: false, Message: []string{err.Error()}}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	byteRes, err := json.Marshal(res)
	// 	if err != nil {
	// 		// w.Write([]bytes([]string{err.Error()}))
	// 		return
	// 	}
	// 	w.Write(byteRes)
	// 	return
	// }

	// upgrade to websocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		res := Response{IsSuccessful: false, Message: []string{err.Error()}}
		byteRes, _ := json.Marshal(res)
		w.Write(byteRes)
		return
	}
	// log.DebugLogger.Println("111")

	// create new client
	client := NewClient(conn, r)
	// err = m.addClient(client)
	// if err != nil {
	// 	log.ErrorLogger.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	res := Response{IsSuccessful: false, Message: []string{err.Error()}}
	// 	byteRes, _ := json.Marshal(res)
	// 	w.Write(byteRes)
	// 	return
	// }
	go client.ReadMessage()
	go client.WriteMessage()
	go client.StreamKafka()

}

// func Stream() {
// 	log.InfoLogger.Println(len(m.clients))
// 	for len(m.clients) > 0 {
// 		log.InfoLogger.Println("streaming")
// 		// file, ferr := os.Open("dummy_logs.txt")
// 		file, ferr := os.Open("logs_message.txt")
// 		if ferr != nil {
// 			log.ErrorLogger.Panic(ferr)
// 		}
// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() && (len(m.clients) > 0) {
// 			line := scanner.Text()
// 			data := utility.GenerateData(line)
// 			message := model.Message{Type: 3, Data: data}
// 			for i := range m.clients {
// 				i.egress <- message
// 			}
// 			time.Sleep(time.Second / 2)
// 		}

// 		log.InfoLogger.Printf("Message Received: %+v\n", "logs_message.txt")
// 	}
// }

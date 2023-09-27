package websocket

import (
	"context"
	"encoding/base64"
	"fmt"
	. "go-fiber-auth/database"
	. "go-fiber-auth/database/schemas"
	"log"
	"os"
	"time"

	// Import the encoding/json package

	"go-fiber-auth/utilities"

	"github.com/gofiber/websocket/v2"
)

// Upgraded websocket request
func handleWebsocket(c *websocket.Conn) {

	var videoFile *os.File
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if string(msg) == "start_recording" {
			// userName := strings.Replace(string(msg), "start_recording", "", -1)
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			// userName := "admin"
			filePath := fmt.Sprintf("./static/video/%s.webm", "khpdev"+"-"+string(timestamp))
			file, err := os.Create(filePath)
			if err == nil {
				videoFile = file
				c.WriteMessage(websocket.TextMessage, []byte("accepted_recording"))
				fmt.Println("start_recording_after")
			} else {
				fmt.Println("Error creating file:", err)
				// Handle the error, perhaps by sending an error response to the client
				return
			}

			now := utilities.MakeTimestamp()

			VideoFileLogCollection := Instance.Database.Collection("Videofilelog")
			VideoFileLog := new(Videofilelog)

			VideoFileLog.FileName = fmt.Sprintf(filePath)
			VideoFileLog.Created = now
			VideoFileLog.Updated = now
			_, insertionError := VideoFileLogCollection.InsertOne(context.Background(), VideoFileLog)
			if insertionError != nil {
				fmt.Println("insertionError")
			}

		} else if string(msg) == "stop_recording" {
			if videoFile != nil {
				videoFile.Close()
			} else {
				fmt.Println("No active recording to stop")
				// Handle the error, perhaps by sending an error response to the client
				return
			}

			now := utilities.MakeTimestamp()
			VideoFileLogCollection := Instance.Database.Collection("Videofilelog")
			VideoFileLog := new(Videofilelog)
			VideoFileLog.EndTime = now
			_, insertionError := VideoFileLogCollection.InsertOne(context.Background(), VideoFileLog)
			if insertionError != nil {
				fmt.Println("insertionError")
			}
		} else {
			decodedBytes, err := base64.StdEncoding.DecodeString(string(msg))
			if err != nil {
				fmt.Println("Error decoding message:", err)
				// Handle the error, perhaps by sending an error response to the client
				return
			}

			if videoFile != nil {
				videoFile.Write(decodedBytes)
			} else {
				fmt.Println("No active recording to write to")
				// Handle the error, perhaps by sending an error response to the client
				return
			}
		}
	}
}

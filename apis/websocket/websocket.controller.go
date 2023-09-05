package websocket

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"
	"context"


	. "go-fiber-auth/database"
	. "go-fiber-auth/database/schemas"

	"github.com/gofiber/websocket/v2"
	"go-fiber-auth/utilities"

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
			file, err := os.Create(fmt.Sprintf("./static/video/%d.webm", time.Now().Unix()))
			if err == nil {
				videoFile = file
				c.WriteMessage(websocket.TextMessage, []byte("accepted_recording"))
				fmt.Println("start_recording_after")
			}

			now := utilities.MakeTimestamp()
			
			VideoFileLogCollection := Instance.Database.Collection("Videofilelog")
			VideoFileLog := new(Videofilelog)
			VideoFileLog.FileName = fmt.Sprintf("./static/video/%d.webm", time.Now().Unix())
			VideoFileLog.Created = now
			VideoFileLog.Updated = now
			_, insertionError := VideoFileLogCollection.InsertOne(context.Background(), VideoFileLog)
			if insertionError != nil {
				fmt.Println("insertionError")
			}

		} else if string(msg) == "stop_recording" {
			videoFile.Close()
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
			if err == nil {
				videoFile.Write(decodedBytes)
			}
		}
	}
}



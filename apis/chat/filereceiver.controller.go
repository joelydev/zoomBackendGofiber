package chat

import (

	"os"
	"fmt"
	"time"
	"strings"

	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"go-fiber-auth/configuration"
	"go-fiber-auth/utilities"

	. "go-fiber-auth/database"
	. "go-fiber-auth/database/schemas"
)

// Handle msg
func fileReceiver(ctx *fiber.Ctx) error {
	fmt.Println("fileReceiver:Call")
	var body FileTransferRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	fmt.Println("body.fileReceiver:%+v", body)

	currentTime := time.Now()
	fmt.Println("body.currentTime:%+v", currentTime)

	authHeader := ctx.Get("Authorization") // Get the token from the request header
	// Check if the header is not empty and starts with "Bearer"
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
    }

	// Extract the token value after "Bearer "
    authToken := strings.Split(authHeader, " ")[1]
	// parse JWT
	trimmedToken := strings.TrimSpace(authToken)
	claims, parsingError := utilities.ParseClaims(trimmedToken)
	fmt.Println("claims.UserId:%+v", claims.UserId)

	if parsingError != nil {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token Analyze Error"})
    }
	// Decode Base64 data
	decodedFileData, err := base64.StdEncoding.DecodeString(string(body.XmppMsgData))
	decodedFileName, err := base64.StdEncoding.DecodeString(body.FileName)
	// decodedFileObj, err := base64.StdEncoding.DecodeString(body.FileObj)
	// decodedFileSize, err := base64.StdEncoding.DecodeString(body.FileSize)
	// fmt.Println("base64.StdEncoding.DecodeString(body.FileSize):%+v", string(base64.StdEncoding.DecodeString(body.FileSize)))
	fmt.Println("body.decodedFileData:%+v", decodedFileData)
	fmt.Println("body.decodedFileName:%+v",string(decodedFileName))
	// fmt.Println("body.decodedFileObj:%+v", string(decodedFileObj))
	// fmt.Println("body.decodedFileSize:%+v", string(decodedFileSize))
	
	if err != nil {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Decode Error"})
    }

	// "eyJzdiI6IjAwMDAwMiIsImFsZyI6IkhTNTEyIiwidiI6IjIuMCIsImtpZCI6ImY2NmU0NzY0LTBhMGYtNDU4Mi1hYTBmLWEzNWM2ODhmZGY1OSJ9.eyJhdWQiOiJtdHUiLCJ1aWQiOiJpLTMwRVFQMlFyRzQtLXNVZTZ3b05RIiwibmJmIjoxNjkzMTIxMzA5LCJjc3IiOiIiLCJqaWQiOiJpLTMwZXFwMnFyZzQtLXN1ZTZ3b25xQHhtcHAuem9vbS51cyIsImxvZ2dlZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly93ZWIuem9vbS51cyIsIm1pZCI6IkZEKzJPS2VZUzN1RTRwbEtFNzZlOFE9PSIsImV4cCI6MTY5MzIyOTMwOSwiaWF0IjoxNjkzMTIxMzA5LCJqdGkiOiI2ZjdkNTQ1ZC04OWRkLTRhMDctOGY5OC0xMTkwY2ZlNTE2OGIiLCJtbm8iOjg4NDg1MjU5NTUyfQ.usP80GGfMrDugJroHem8sfYNzAjPgdNZ6Za4x4YouY2G9I0nrQmeB1IbjXEt-bcOe7tUMD0X_e-MwAbz_eUk5A"
	// Save the file
	decodedFileData1, err := base64.StdEncoding.DecodeString("eyJzdiI6IjAwM")

	filePath := "static/chat/upload/" + string(decodedFileName)
	fmt.Println("body.decodedFilepathe:%+v", filePath)

	fmt.Println("----------------body.file_data--------------:%+v", decodedFileData)
	err = os.WriteFile(filePath, decodedFileData1, 0644)
	if err != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// Insert upload file log into MongoDB
	// load Upload File Log Schema
	UploadFileLogCollection := Instance.Database.Collection("Uploadfilelog")

	// create a new Chat Message record, insert it
	now := utilities.MakeTimestamp()
	NewUploadFileLog := new(Uploadfilelog)

	NewUploadFileLog.FileName 		= string(decodedFileName)
	NewUploadFileLog.FileObj 		= ""
	// NewUploadFileLog.FileSize 		= string(decodedFileSize)
	// NewUploadFileLog.FileType 		= body.FileType
	// NewUploadFileLog.ReceiverType	= body.ReceiverType
	// NewUploadFileLog.TransferType 	= body.TransferType
	NewUploadFileLog.ID 			= ""
	NewUploadFileLog.UserId 		= claims.UserId
	NewUploadFileLog.Updated 		= now
	NewUploadFileLog.Created 		= now
	insertionChatResult, insertionError := UploadFileLogCollection.InsertOne(ctx.Context(), NewUploadFileLog)

	fmt.Println("body.insertionChatResult:%+v", insertionChatResult)
	fmt.Println("body.insertionError:%+v", insertionError)
	
	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"state": "success",
		},
	})
}
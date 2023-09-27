package chat

import (
	"os"
	// "strconv"
	"context"
	"fmt"
	"go-fiber-auth/configuration"
	"log"
	"strings"
	"time"

	// "go-fiber-auth/api/chat/excel"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "go-fiber-auth/database"
	. "go-fiber-auth/database/schemas"

	"go-fiber-auth/utilities"

	"github.com/xuri/excelize/v2"
)

// Handle msg
func chatMsg(ctx *fiber.Ctx) error {
	// check data
	// var reqBody = ctx.Body()
	var body ChatMsgRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	fmt.Println("body.Message:%+v", body)

	currentTime := time.Now()

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

	// load User schema
	UserCollection := Instance.Database.Collection("User")

	// Parse the ID into an ObjectID
	objectID, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Find the document using the ObjectID
	var userResult bson.M
	err = UserCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&userResult)
	if err != nil {
		return ctx.Status(404).SendString("Document not found")
	}
	fmt.Println("ctx.JSON(result):%+v", userResult["name"].(string))

	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// Insert message log into MongoDB
	// load User schema
	ChatMsgCollection := Instance.Database.Collection("Chatmsg")

	// create a new Chat Message record, insert it and get back the ID
	now := utilities.MakeTimestamp()
	NewChatMsg := new(Chatmsg)
	NewChatMsg.Created = now
	NewChatMsg.Message = body.Message
	NewChatMsg.MessageType = body.MessageType
	NewChatMsg.From = body.Sender
	NewChatMsg.To = body.Receiver
	NewChatMsg.ID = ""
	NewChatMsg.UserId = claims.UserId
	NewChatMsg.UserName = userResult["name"].(string)
	NewChatMsg.Updated = now
	insertionChatResult, insertionError := ChatMsgCollection.InsertOne(ctx.Context(), NewChatMsg)

	// // Create a document to insert
	// document := bson.M{	"message": body.Message,
	// 					"messagetype": body.MessageType,
	// 					"filepath": "",
	// 					"from": body.Sender,
	// 					"to": body.Receiver,
	// 					"userId": "value2",
	// 					"username": "value2"
	// 					"created": "value2",
	// 					"updated": "value2",
	// 				}

	// // Insert the document
	// insertResult, err := collection.InsertOne(context.Background(), document)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Error inserting record")
	// }

	// // Retrieve the inserted ID from the insert result
	// insertedID := insertResult.InsertedID

	fmt.Println("insertionResult:%+v", insertionChatResult)
	fmt.Println("insertionError:%+v", insertionError)
	fmt.Println("body.Order:%+v", body.Order)

	CurrentTimeFormat := string(currentTime.Format("2006-01-02-15-04-05"))
	ExcelFileLogCollection := Instance.Database.Collection("Excelfilelog")
	if body.Order == 0 {
		// Create Excel file

		fmt.Println("Create Excel file")
		newFile := excelize.NewFile()
		defer func() {
			if err := newFile.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		for idx, row := range [][]interface{}{
			{"Order_id",
				"user_name",
				"message",
				"from",
				"to",
				"time",
			},
			{body.Order, userResult["name"].(string), body.Message, body.Sender, body.Receiver, CurrentTimeFormat},
		} {
			cell, err := excelize.CoordinatesToCellName(1, idx+1)
			if err != nil {
				fmt.Println(err)
				// return
			}
			newFile.SetSheetRow("Sheet1", cell, &row)
		}

		// Save spreadsheet by the given path.
		createFilePath := fmt.Sprintf("%s", "static/chat/"+userResult["name"].(string)+"-"+CurrentTimeFormat+".xlsx")

		if err := newFile.SaveAs(createFilePath); err != nil {
			fmt.Println(err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			return nil
		}
		// Replace backslashes with slashes
		cwd = strings.Replace(cwd, "\\", "/", -1)

		fmt.Println("os.Getwd:%+v", cwd)

		// create a new Chat Message record, insert it and get back the ID
		now := utilities.MakeTimestamp()
		NewExcelfilelog := new(Excelfilelog)
		NewExcelfilelog.FileName = fmt.Sprintf("%s", userResult["name"].(string)+"-"+CurrentTimeFormat+".xlsx")
		NewExcelfilelog.FilePath = createFilePath
		NewExcelfilelog.Creator = userResult["name"].(string)
		NewExcelfilelog.CreatorId = claims.UserId
		NewExcelfilelog.ID = ""
		NewExcelfilelog.Updated = now
		NewExcelfilelog.Created = now
		insertionResult, insertionError := ExcelFileLogCollection.InsertOne(ctx.Context(), NewExcelfilelog)

		fmt.Println("ExcelFileLogCollectioninsertionResult:%+v", insertionResult.InsertedID)
		fmt.Println("ExcelFileLogCollectioninsertionError:%+v", insertionError)

		// Define update criteria
		filter := bson.M{"_id": insertionChatResult.InsertedID}

		// Define update operation
		update := bson.M{"$set": bson.M{"filepath": NewExcelfilelog.FilePath}}

		// Perform the update
		_, err = ChatMsgCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error updating record")
		}

	} else {

		fmt.Println("ExistingFIle")
		// Get xls file paht from database
		// Define your specific condition using the bson.M map
		condition := bson.M{"creatorid": claims.UserId}
		fmt.Println("condition:%+v", condition)
		// Sort by _id in descending order to get the latest document
		options := options.FindOne().SetSort(bson.D{{"_id", -1}})
		var excelResult bson.M
		err = ExcelFileLogCollection.FindOne(context.Background(), condition, options).Decode(&excelResult)
		fmt.Println("excelResultCollection:%+v", excelResult)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error retrieving last inserted record by condition")
		}

		fmt.Println("find.result:%+v", ctx.JSON(excelResult))
		fmt.Println("excelResult.(string):%+v", excelResult["filename"].(string))
		// Open an existing Excel file
		filePath := fmt.Sprintf("%s", "static/chat/"+excelResult["filename"].(string))
		file, err := excelize.OpenFile(filePath)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error opening the Excel file")
		}

		// Get the name of the first sheet
		for sheetIndex, sheetName := range file.GetSheetMap() {
			fmt.Printf("sheetIndex: %d\n", sheetIndex)

			rows, err := file.GetRows(sheetName)
			if err != nil {
				fmt.Println("Error getting rows from sheet:", err)
				// Handle the error, possibly by returning or reporting it
				continue // Skip this sheet and continue with the next one
			}

			nextRowIndex := len(rows) + 1

			newData := []interface{}{body.Order, userResult["name"].(string), body.Message, body.Sender, body.Receiver, string(currentTime.Format("2006-01-02 15-04-05"))}

			for colIdx, cellValue := range newData {
				cellName, _ := excelize.CoordinatesToCellName(colIdx+1, nextRowIndex)
				if err := file.SetCellValue(sheetName, cellName, cellValue); err != nil {
					fmt.Println("Error setting cell value:", err)
					// Handle the error, possibly by returning or reporting it
					continue // Continue with the next cell or row
				}
			}
		}

		// Save changes to the Excel file
		if err := file.Save(); err != nil {
			log.Fatal(err)
		}
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"token": "token",
			"user":  "user",
		},
	})
}

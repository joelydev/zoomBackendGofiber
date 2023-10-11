package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"go-fiber-auth/configuration"
	. "go-fiber-auth/database"
	. "go-fiber-auth/database/schemas"
	"go-fiber-auth/utilities"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"encoding/base64"

	"io"
)

// Handle signing in
func signIn(ctx *fiber.Ctx) error {
	// check data
	var body SignInUserRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	fmt.Println("body.Email:", body.Email)

	email := body.Email
	password := body.Password
	if email == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedEmail := strings.TrimSpace(email)
	trimmedPassword := strings.TrimSpace(password)
	if trimmedEmail == "" || trimmedPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// load User schema
	UserCollection := Instance.Database.Collection("User")

	// find a user
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: trimmedEmail}},
	)
	userRecord := &User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// load Password schema
	PasswordCollection := Instance.Database.Collection("Password")

	// find a password
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
	)
	passwordRecord := &Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// compare hashes
	passwordIsValid := utilities.CompareHashes(trimmedPassword, passwordRecord.Hash)
	if !passwordIsValid {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	accessExpiration, expirationError := strconv.Atoi(os.Getenv("TOKENS_ACCESS_EXPIRATION"))
	if expirationError != nil {
		accessExpiration = 24
	}
	token, tokenError := utilities.GenerateJWT(utilities.GenerateJWTParams{
		ExpiresIn: int64(accessExpiration),
		UserId:    userRecord.ID,
	})
	if tokenError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// requestUrl := fmt.Sprintf("http://192.168.5.251:5000/register?ip=%s", ctx.IP())
	requestUrl := fmt.Sprintf("%s://%s:%s/register?ip=%s", os.Getenv("PROXY_SCHEME"), os.Getenv("PROXY_SERVER_IP"), os.Getenv("PROXY_NODE_SERVICE_PORT"), ctx.IP())

	proxyState, proxyStateError := proxyServerAllowIpRequest(requestUrl)
	if proxyStateError != nil {
		return utilities.Response(utilities.ResponseParams{
			// Proxy Node Service is not running on Linux.
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.ProxyNodeServerError,
			Status: fiber.StatusServiceUnavailable,
		})
	}

	log.Printf("proxyState response from Node.js: %v", proxyState)

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"token":         token,
			"user":          userRecord,
			"proxyUsername": os.Getenv("PROXY_USER_NAME"),
			"proxyPassword": os.Getenv("PROXY_USER_PASS"),
			"proxyScheme":   os.Getenv("PROXY_SCHEME"),
			"proxyServerIp": os.Getenv("PROXY_SERVER_IP"),
			"proxyPort":     os.Getenv("PROXY_PORT"),
			// "proxyState":    proxyState,
		},
	})
}

func encryptData(dataToEncrypt string, keyString string) (string, error) {
	key := []byte(keyString)
	plaintext := []byte(dataToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func proxyServerAllowIpRequest(ipAddress string) (string, error) {
	resp, err := http.Get(ipAddress)
	log.Printf("---------resp------------: %v", resp)
	log.Printf("---------resp_err------------: %v", err)

	if err != nil {
		log.Printf("Error sending GET request to Node.js server: %v", err)
		return "", err // Return the error object
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code from Node.js server: %d", resp.StatusCode)
		return "", fmt.Errorf("Internal Server Error")
	}

	// Read and return the response from the Node.js server
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from Node.js server: %v", err)
		return "", err
	}

	return string(body), nil // Return the response and a nil error
}

package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"go-fiber-auth/configuration"
	"go-fiber-auth/utilities"
	"os"
)

// Handle signing in
func unregister(ctx *fiber.Ctx) error {
	// check data
	log.Printf("unregister_request_ip: %v", ctx.IP())
	requestUrl := fmt.Sprintf("%s://%s:%s/unregister?ip=%s", os.Getenv("PROXY_SCHEME"), os.Getenv("PROXY_SERVER_IP"), os.Getenv("PROXY_NODE_SERVICE_PORT"), ctx.IP())
	log.Printf("requestUrl_: %v", requestUrl)
	proxyState, proxyStateError := proxyServerAllowIpRequest(requestUrl)
	log.Printf("requestUrl_proxyState: %v", proxyState)
	log.Printf("proxyStateError: %v", proxyStateError)
	if proxyStateError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.ErrInternalServerError.Code,
		})
	}

	log.Printf("proxyStateError_after")
	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"proxyState": proxyState,
		},
	})
}

func proxyServerAllowIpRequest(ipAddress string) (string, error) {
	resp, err := http.Get(ipAddress)
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

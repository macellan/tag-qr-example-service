package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
)

// OpenGateRequest When Macellan SuperApp comes to you to perform the action, it comes with the following values.
type OpenGateRequest struct {
	Price              string `json:"price"`                // The Price net amount to be deducted from the user balance when you call the success sendCallback
	Point              string `json:"point"`                // The Point net amount to be deducted from the user point when you call the success sendCallback
	RefCode            string `json:"ref_code"`             // RefCode is the Tag-QR reference code.
	UserId             string `json:"user_id"`              // UserId of the user who read the QR
	OrderId            string `json:"order_id"`             // OrderId is the unique ID of the Macellan SuperApp side of the payment transaction.
	CallbackSuccessUrl string `json:"callback_success_url"` // The CallbackSuccessUrl to be called if the transactions have been completed successfully
	CallbackFailUrl    string `json:"callback_fail_url"`    // CallbackFailUrl to call when something goes wrong
	Hash               string `json:"hash"`                 // Hash required to verify that the request was sent by the Macellan SuperApp
}

// OpenGateResponse SuperApp does not look at the body of successful requests. But in unsuccessful requests, you should return as follows.
type OpenGateResponse struct {
	Message string `json:"message"`
}

// CallbackRequest You should use the following structure to send a request back to SuperApp.
type CallbackRequest struct {
	Hash string `json:"hash"` // The Hash required to verify that the acknowledgment was sent by you
}

func OpenGate(c *fiber.Ctx) error {
	request := new(OpenGateRequest)

	if err := c.BodyParser(request); err != nil {
		return err
	}

	if !verifyRequestHash(request) {
		// All operations other than code 200 are considered as failed requests by Macellan SuperApp.
		return c.Status(fiber.StatusForbidden).
			JSON(OpenGateResponse{
				Message: "Hash is invalid",
			})
	}

	// Here you can do everything necessary to open the door.
	// If everything is ok and the door is opened, you should call success.
	go sendCallback(true, request)

	// if something went wrong call fail
	// defer sendCallback(false, request)

	// If you have successfully received the request, you should return status code 200.
	// It is not a value depending on whether you open the door or not. You simply verify that you have received the request.
	// You should notify with sendCallback requests whether the door is opened or not.
	return c.Status(fiber.StatusOK).
		JSON(OpenGateResponse{
			Message: "Success",
		})
}

// Verifies that the request to the service was sent by Macellan SuperApp.
func verifyRequestHash(request *OpenGateRequest) bool {
	s := []string{
		os.Getenv("SALT"),
		request.CallbackFailUrl,
		request.CallbackSuccessUrl,
		request.Price,
	}

	calculatedHash := hash(s)

	// println("Request Hash => ", request.Hash)
	// println("Calculated Hash => ", calculatedHash)

	return request.Hash == calculatedHash
}

func sendCallback(isSuccess bool, request *OpenGateRequest) {
	jsonData := prepareCallbackData(request)
	callbackURL := request.CallbackSuccessUrl

	if !isSuccess {
		callbackURL = request.CallbackFailUrl
	}

	_, _ = http.Post(callbackURL, "application/json", bytes.NewBuffer(jsonData))
}

func prepareCallbackData(request *OpenGateRequest) []byte {
	s := []string{
		request.Price,
		request.CallbackSuccessUrl,
		request.CallbackFailUrl,
		os.Getenv("SALT"),
	}

	calculatedHash := hash(s)

	jsonData, _ := json.Marshal(CallbackRequest{
		Hash: calculatedHash,
	})

	return jsonData
}

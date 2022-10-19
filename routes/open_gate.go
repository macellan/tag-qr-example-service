package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
)

// OpenGateRequest When Alternatif SuperApp comes to you to perform the action, it comes with the following values.
type OpenGateRequest struct {
	Price              string `json:"price"`                // The Price net amount to be deducted from the user balance when you call the success callback
	Point              string `json:"point"`                // The Point net amount to be deducted from the user point when you call the success callback
	RefCode            string `json:"ref_code"`             // RefCode is the Tag-QR reference code.
	UserId             string `json:"user_id"`              // UserId of the user who read the QR
	OrderId            string `json:"order_id"`             // OrderId is the unique ID of the Alternate SuperApp side of the payment transaction.
	CallbackSuccessUrl string `json:"callback_success_url"` // The CallbackSuccessUrl to be called if the transactions have been completed successfully
	CallbackFailUrl    string `json:"callback_fail_url"`    // CallbackFailUrl to call when something goes wrong
	Hash               string `json:"hash"`                 // Hash required to verify that the request was sent by the Alternatif SuperApp
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
	r := new(OpenGateRequest)

	if err := c.BodyParser(r); err != nil {
		return err
	}

	if verifyRequestHash(r) == false {
		// All operations other than code 200 are considered as failed requests by Alternatif SuperApp.
		return c.Status(fiber.StatusForbidden).
			JSON(OpenGateResponse{
				Message: "Hash is invalid",
			})
	}

	// Here you can do everything necessary to open the door.
	// If everything is ok and the door is opened, you should call success.
	go pingCallbackSuccess(r)

	// if something went wrong call fail
	// defer pingCallbackSuccess(r)

	// If you have successfully received the request, you should return status code 200.
	// It is not a value depending on whether you open the door or not. You simply verify that you have received the request.
	// You should notify with callback requests whether the door is opened or not.
	return c.Status(fiber.StatusOK).
		JSON(OpenGateResponse{
			Message: "Success",
		})
}

// Verifies that the request to the service was sent by Alternatif SuperApp.
func verifyRequestHash(r *OpenGateRequest) bool {
	s := []string{
		os.Getenv("SALT"),
		r.CallbackFailUrl,
		r.CallbackSuccessUrl,
		r.Price,
	}

	cHash := hash(s)

	// println("Request Hash => ", r.Hash)
	// println("Calculated Hash => ", cHash)

	return r.Hash == cHash
}

func callbackPostData(r *OpenGateRequest) []byte {
	s := []string{
		r.Price,
		r.CallbackSuccessUrl,
		r.CallbackFailUrl,
		os.Getenv("SALT"),
	}

	cHash := hash(s)

	jsonData, _ := json.Marshal(CallbackRequest{
		Hash: cHash,
	})

	return jsonData
}

func pingCallbackSuccess(r *OpenGateRequest) {
	jsonData := callbackPostData(r)

	_, _ = http.Post(r.CallbackSuccessUrl, "application/json", bytes.NewBuffer(jsonData))
}

//goland:noinspection GoUnusedFunction
func pingCallbackFail(r *OpenGateRequest) {
	jsonData := callbackPostData(r)

	_, _ = http.Post(r.CallbackFailUrl, "application/json", bytes.NewBuffer(jsonData))
}

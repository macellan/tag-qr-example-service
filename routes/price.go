package routes

import (
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"os"
)

// PriceRequest Macellan SuperApp will come up with these values to query the amount.
type PriceRequest struct {
	RefCode string `json:"ref_code"` // RefCode is the Tag-QR reference code.
	Hash    string `json:"hash"`     // Hash required to verify that the request was sent by the Macellan SuperApp
	UserId  string `json:"user_id"`  // UserId of the user who read the QR
	OrderId string `json:"order_id"` // OrderId is the unique ID of the Macellan SuperApp side of the payment transaction.
}

// PriceSuccessResponse Amount response should be given with the following structure.
type PriceSuccessResponse struct {
	Price float32 `json:"price"` // Price indicates how much will be deducted from the balance.
}

// PriceFailResponse If you want to return the answer as an error, it must be answered with the following values.
type PriceFailResponse struct {
	Message     string `json:"message,omitempty"`      // This is the Message you will send to Macellan SuperApp API.
	UserMessage string `json:"user_message,omitempty"` // UserMessage is the message you want to show the user. Can be null
}

func GetPrice(c *fiber.Ctx) error {
	request := new(PriceRequest)

	if err := c.BodyParser(request); err != nil {
		return err
	}

	if !verifyPriceRequestHash(request) {
		// All operations other than code 200 are considered as failed requests by Macellan SuperApp.
		return c.Status(fiber.StatusForbidden).
			JSON(PriceFailResponse{
				Message: "Hash is invalid",
			})
	}

	// Here you can calculate the amount to be deducted from the balance.
	price := calculatePrice()

	return c.Status(fiber.StatusOK).
		JSON(PriceSuccessResponse{
			Price: price,
		})
}

// verifyPriceRequestHash Verifies that the request to the service was sent by Macellan SuperApp.
func verifyPriceRequestHash(request *PriceRequest) bool {
	s := []string{
		request.RefCode,
		request.UserId,
		request.OrderId,
		os.Getenv("SALT"),
	}

	calculatedHash := hash(s)

	// println("Request Hash => ", request.Hash)
	// println("Calculated Hash => ", calculatedHash)

	return request.Hash == calculatedHash
}

// This is a sample code written so that it doesn't always look the same.
func calculatePrice() float32 {
	return 1 + rand.Float32()*(10-1)
}

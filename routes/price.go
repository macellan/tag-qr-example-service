package routes

import (
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"os"
)

// PriceRequest Alternatif SuperApp will come up with these values to query the amount.
type PriceRequest struct {
	RefCode string `json:"ref_code"` // RefCode is the Tag-QR reference code.
	Hash    string `json:"hash"`     // Hash required to verify that the request was sent by the Alternatif SuperApp
	UserId  string `json:"user_id"`  // UserId of the user who read the QR
	OrderId string `json:"order_id"` // OrderId is the unique ID of the Alternate SuperApp side of the payment transaction.
}

// PriceSuccessResponse Amount response should be given with the following structure.
type PriceSuccessResponse struct {
	Price float32 `json:"price"` // Price indicates how much will be deducted from the balance.
}

// PriceFailResponse If you want to return the answer as an error, it must be answered with the following values.
type PriceFailResponse struct {
	Message     string `json:"message,omitempty"`      // This is the Message you will send to Alternatif SuperApp API.
	UserMessage string `json:"user_message,omitempty"` // UserMessage is the message you want to show the user. Can be null
}

func GetPrice(c *fiber.Ctx) error {
	r := new(PriceRequest)

	if err := c.BodyParser(r); err != nil {
		return err
	}

	if verifyHash(r) == false {
		// All operations other than code 200 are considered as failed requests by Alternatif SuperApp.
		return c.Status(fiber.StatusForbidden).
			JSON(PriceFailResponse{
				Message: "Hash is invalid",
			})
	}

	// Here you can calculate the amount to be deducted from the balance.
	price := calcPrice()

	return c.Status(fiber.StatusOK).
		JSON(PriceSuccessResponse{
			Price: price,
		})
}

// verifyHash Verifies that the request to the service was sent by Alternatif SuperApp.
func verifyHash(r *PriceRequest) bool {
	s := []string{
		r.RefCode,
		r.UserId,
		r.OrderId,
		os.Getenv("SALT"),
	}

	cHash := hash(s)

	// println("Request Hash => ", r.Hash)
	// println("Calculated Hash => ", cHash)

	return r.Hash == cHash
}

// This is a sample code written so that it doesn't always look the same.
func calcPrice() float32 {
	return 1 + rand.Float32()*(10-1)
}

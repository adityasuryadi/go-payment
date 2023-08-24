package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment/config"
	"payment/controller"
	"payment/exception"
	"payment/repository"
	"payment/service"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: exception.ErrorHandler})
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	paymentRepository := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepository, db)
	paymentController := controller.NewPaymentController(paymentService)
	paymentController.Route(app)
	return app
}

var app = Setup()

func TestCreatePaymentEmptyField(t *testing.T) {
	t.Run("failed", func(t *testing.T) {
		requestBody := strings.NewReader(`{
			"cust_name": "",
			"email": "",
			"phone": "",
			"service_name": "",
			"service_id": "",
			"booking_date": ""
		  }`)
		request := httptest.NewRequest(http.MethodPost, "/payment", requestBody)
		request.Header.Add("Content-Type", "application/json")
		res, _ := app.Test(request)
		body, _ := ioutil.ReadAll(res.Body)
		response := make(map[string]interface{})
		json.Unmarshal(body, &response)
		parse := response
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "BAD_REQUEST", parse["status"])
		assert.Equal(t, float64(400), parse["code"])
	})
}

func TestCreatePaymentEmptyName(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "",
			"email": "adit@mail.com",
			"phone": "0896123132123",
			"service_name": "Service Iphone",
			"service_id": "30",
			"booking_date": "` + bookingDate + `"
		  }`)
	request := httptest.NewRequest(http.MethodPost, "/payment", requestBody)
	request.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(request)
	body, _ := ioutil.ReadAll(res.Body)
	response := make(map[string]interface{})
	json.Unmarshal(body, &response)
	parse := response
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "BAD_REQUEST", parse["status"])
	assert.Equal(t, float64(400), parse["code"])

	data := response["data"].([]interface{})
	for _, val := range data {
		value := val.(map[string]interface{})
		assert.Equal(t, "cust_name", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

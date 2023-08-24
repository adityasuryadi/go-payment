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

func TestCreatePaymentEmptyEmail(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "",
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
		assert.Equal(t, "email", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestCreatePaymentWrongEmail(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@",
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
		assert.Equal(t, "email", value["field"])
		assert.Equal(t, "format email salah", value["message"])
	}
}

func TestCreatePaymentEmptyPhone(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "",
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
		assert.Equal(t, "phone", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestCreatePaymentPhoneNotNumber(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "0812499ff111",
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
		assert.Equal(t, "phone", value["field"])
		assert.Equal(t, "harus numeric", value["message"])
	}
}

func TestCreatePaymentPhoneMin(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "0812499",
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
		assert.Equal(t, "phone", value["field"])
		assert.Equal(t, "minimal 8 karakter", value["message"])
	}
}

func TestCreatePaymentEmptyServiceName(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "08121234211",
			"service_name": "",
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
		assert.Equal(t, "service_name", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestCreatePaymentEmptyServiceId(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "08121234211",
			"service_name": "Service",
			"service_id": "",
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
		assert.Equal(t, "service_id", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestCreatePaymentEmptybookingDate(t *testing.T) {
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "08121234211",
			"service_name": "iphone",
			"service_id": "30",
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

	data := response["data"].([]interface{})
	for _, val := range data {
		value := val.(map[string]interface{})
		assert.Equal(t, "booking_date", value["field"])
		assert.Equal(t, "field tidak boleh kosong", value["message"])
	}
}

func TestCreatePaymentBookingDateLowerThanNow(t *testing.T) {
	bookingDate := time.Now().Add(time.Hour * -24).Format("2006-01-02")
	requestBody := strings.NewReader(`{
			"cust_name": "Adit",
			"email": "adi@mail.com",
			"phone": "08121234211",
			"service_name": "iphone",
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
		assert.Equal(t, "booking_date", value["field"])
		assert.Equal(t, "harus lebih besar dari tanggal hari ini", value["message"])
	}
}

package service

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"payment/model"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func NewFaspayService() FaspayService {
	return &FaspayServiceImpl{}
}

type FaspayServiceImpl struct {
}

// CreatePaymentExpress implements FaspayService
func (*FaspayServiceImpl) CreatePaymentExpress(request model.CreateFaspayPaymentRequest) (*model.CreatePaymentFaspayResponse, error) {
	merchantId := os.Getenv("FASPAY_MERCHANT_ID")
	userId := os.Getenv("FASPAY_USER_ID")
	password := os.Getenv("FASPAY_PASSWORD")

	url := "https://xpress.faspay.co.id/v4/post"
	if os.Getenv("APP_ENV") == "dev" {
		url = "https://xpress-sandbox.faspay.co.id/v4/post"
	}
	shaEncrypt := sha1.New()
	md5Encrypt := md5.New()

	plainSignature := userId + password + request.BillNo + strconv.Itoa(int(price))

	md5Encrypt.Write([]byte(plainSignature))
	md5Signature := md5Encrypt.Sum(nil)

	shaEncrypt.Write([]byte(string(fmt.Sprintf("%x", md5Signature))))
	signature := shaEncrypt.Sum(nil)

	request.Signature = string(fmt.Sprintf("%x", signature))
	request.ReturnUrl = os.Getenv("FASPAY_CALLBACK_URL")
	request.MerchantId = merchantId

	client := resty.New()
	resp, err := client.R().
		SetBody(request).
		Post(url)
	response := make(map[string]interface{})
	json.Unmarshal(resp.Body(), &response)
	var faspayResponse model.CreatePaymentFaspayResponse
	if response["response_code"] == "00" {
		err = nil
		faspayResponse = model.CreatePaymentFaspayResponse{
			BillNo:       response["bill_no"].(string),
			MerchantId:   merchantId,
			Merchant:     response["merchant"].(string),
			ResponseCode: response["response_code"].(string),
			ResponseDesc: response["response_desc"].(string),
			RedirectUrl:  response["redirect_url"].(string),
		}
	}

	if response["response_code"] != "00" {
		return nil, err
	}
	return &faspayResponse, nil
}

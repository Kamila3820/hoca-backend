package helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/helper/model"
	"github.com/sirupsen/logrus"
)

func scbGetAccessToken() string {
	cfg := config.ConfigGetting()
	url := cfg.Scb.ScbUrl + "/v1/oauth/token"

	getTokenBody := *&model.ScbGetTokenRequest{
		ApplicationKey:    cfg.Scb.ScbAppKey,
		ApplicationSecret: cfg.Scb.ScbAppSecret,
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(getTokenBody)
	if err != nil {
		logrus.Error("Unable to marshal JSON ", err)
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Error("Unable to construct request ", err)
		panic(err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept-language", "EN")
	req.Header.Add("resourceOwnerId", cfg.Scb.ScbAppKey)
	req.Header.Add("requestUId", "HOCA_PAYMENT_SYSTEM")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Unable to send request ", err)
		panic(err)
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logrus.Error("Unable to read response body ", err)
		panic(err)
	}

	var tokenResponse *model.ScbGetTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		logrus.Error("Unable to parse response body ", err)
		panic(err)
	}

	return tokenResponse.Data.AccessToken
}

func ScbCreateQrPayment(amount uint, transactionID string) model.ScbCreateQrDataResponse {
	cfg := config.ConfigGetting()
	accessToken := scbGetAccessToken()

	url := cfg.Scb.ScbUrl + "/v1/payment/qrcode/create"
	createQrBody := model.ScbCreateQrPaymentRequest{
		QrType: "PP",
		Amount: strconv.Itoa(int(amount)) + ".00",
		PpType: "BILLERID",
		PpId:   cfg.Scb.ScbBillerId,
		Ref1:   transactionID,
		Ref2:   "HOCA01",
		Ref3:   "HOCA01",
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(createQrBody)
	if err != nil {
		logrus.Error("Unable to marshal JSON ", err)
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Error("Unable to construct request ", err)
		panic(err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept-language", "EN")
	req.Header.Add("resourceOwnerId", cfg.Scb.ScbAppKey)
	req.Header.Add("requestUId", "1234567890")
	req.Header.Add("authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Unable to send request ", err)
		panic(err)
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logrus.Error("Unable to read response body ", err)
		panic(err)
	}

	var qrResponse *model.ScbCreateQrResponse
	if err := json.Unmarshal(body, &qrResponse); err != nil {
		logrus.Error("Unable to parse response body ", err)
		panic(err)
	}

	return qrResponse.Data
}

func ScbInquiryPayment(transactionID string) (*model.ScbInquiryQrDataResponse, error) {
	cfg := config.ConfigGetting()
	accessToken := scbGetAccessToken()

	currentDate := time.Now().Format("2006-01-02")
	url := cfg.Scb.ScbUrl + "/v1/payment/billpayment/inquiry?eventCode=00300100&transactionDate=" + currentDate + "&billerId=" + cfg.Scb.ScbBillerId + "&reference1=" + transactionID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Unable to construct request ", err)
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept-language", "EN")
	req.Header.Add("resourceOwnerId", cfg.Scb.ScbAppKey)
	req.Header.Add("requestUId", "1234567890")
	req.Header.Add("authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Unable to send request ", err)
		return nil, err
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logrus.Error("Unable to read response body ", err)
		return nil, err
	}

	var qrResponse *model.ScbInquiryQrResponse
	if err := json.Unmarshal(body, &qrResponse); err != nil {
		logrus.Error("Unable to parse response body ", err)
		return nil, err
	}

	if qrResponse.Data == nil || len(qrResponse.Data) == 0 {
		return nil, nil
	}

	return &qrResponse.Data[0], nil
}

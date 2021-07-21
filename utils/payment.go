package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/final-project/database"
	"github.com/final-project/models"
	"github.com/sony/sonyflake"
	"log"
	"net/http"
	"strconv"
)

func MomoPayment(amount string) string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	//random orderID and requestID
	b, err := flake.NextID()

	var id int

	db := database.Connect()
	defer db.Close()
	row, _ := db.Query("SELECT MAX(id) FROM orders")
	if row.Next() {
		row.Scan(&id)
	}
	var orderId = strconv.Itoa(id + 1)
	var requestId = strconv.FormatUint(b, 16)
	var endpoint = "https://test-payment.momo.vn/gw_payment/transactionProcessor"
	var partnerCode = "MOMOPSMB20210716"
	var accessKey = "dxcoKDjTD7VLcwMB"
	var serectkey = "w7c1jJSJok65DZybUH0FPsheFsWmSBMZ"
	var orderInfo = "Mua hang tai shop Linh Kien MH"
	var returnUrl = "http://localhost:8080/redirect"
	var notifyurl = "http://localhost:10000/api/v1/ipn-momo"
	var requestType = "captureMoMoWallet"
	var extraData = "merchantName=;merchantId=MOMOPSMB20210716"

	var rawSignature bytes.Buffer
	rawSignature.WriteString("partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&returnUrl=")
	rawSignature.WriteString(returnUrl)
	rawSignature.WriteString("&notifyUrl=")
	rawSignature.WriteString(notifyurl)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(serectkey))
	hmac.Write(rawSignature.Bytes()) // Write Data to it

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))

	var payload = models.Payload{
		PartnerCode: partnerCode,
		AccessKey:   accessKey,
		RequestID:   requestId,
		Amount:      amount,
		OrderID:     orderId,
		OrderInfo:   orderInfo,
		ReturnURL:   returnUrl,
		NotifyURL:   notifyurl,
		ExtraData:   extraData,
		RequestType: requestType,
		Signature:   signature,
	}

	var jsonPayload []byte
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	//send HTTP to momo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalln(err)
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	payUrl := fmt.Sprintf("%v", result["payUrl"])
	return payUrl
}

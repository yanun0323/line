package linepay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type LinePay struct {
	isProduction  bool
	channelID     string
	channelSecret string
}

func NewLinePay(isProduction bool, channelID string, channelSecret string) *LinePay {
	return &LinePay{
		isProduction:  isProduction,
		channelID:     channelID,
		channelSecret: channelSecret,
	}
}

func (lp *LinePay) signKey(clientKey string, msg string) string {
	h := hmac.New(sha256.New, []byte(clientKey))
	h.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type APIRequest struct {
	Method      string
	APIPath     string
	QueryString string
	Data        interface{}
}

func (lp *LinePay) requestOnlineAPI(ctx context.Context, req *APIRequest, parser func([]byte) (any, error)) (any, error) {
	baseURL := "https://sandbox-api-pay.line.me"
	if lp.isProduction {
		baseURL = "https://api-pay.line.me"
	}

	nonce := uuid.New().String()
	var signature string

	if req.Method == "GET" {
		signature = lp.signKey(
			lp.channelSecret,
			lp.channelSecret+req.APIPath+req.QueryString+nonce,
		)
	} else if req.Method == "POST" {
		dataJSON, err := json.Marshal(req.Data)
		if err != nil {
			return "", err
		}
		signature = lp.signKey(
			lp.channelSecret,
			lp.channelSecret+req.APIPath+string(dataJSON)+nonce,
		)
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s%s", baseURL, req.APIPath)
	if req.QueryString != "" {
		url += "&" + req.QueryString
	}

	var reqBody *strings.Reader
	if req.Data != nil {
		dataJSON, err := json.Marshal(req.Data)
		if err != nil {
			return "", err
		}
		reqBody = strings.NewReader(string(dataJSON))
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, reqBody)
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-LINE-ChannelId", lp.channelID)
	httpReq.Header.Set("X-LINE-Authorization", signature)
	httpReq.Header.Set("X-LINE-Authorization-Nonce", nonce)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res, err := parser(body)
	if err != nil {
		return "", errors.Errorf("parse response, err: %+v", err)
	}

	return res, nil
}

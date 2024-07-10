package service

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	vision "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/twilio/twilio-go"
	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	errGettingImageFromUrl = "error occurred while getting image from url"
)

type TwilioCreds struct {
	TwilioSID   string
	TwilioToken string
}

type Client struct {
	logger        logger.Logger
	httpClient    http.Client
	textExtractor ITextExtractor
	twilioCreds   twilio.ClientParams
}

func NewHttpClient(logger logger.Logger, textExtractor ITextExtractor, creds twilio.ClientParams) *Client {
	return &Client{
		logger:        logger,
		httpClient:    http.Client{},
		textExtractor: textExtractor,
		twilioCreds:   creds,
	}
}

func (c *Client) FetchMedia(ctx context.Context, mediaUrl string) (*vision.Image, error) {

	req, err := http.NewRequest(http.MethodGet, mediaUrl, nil)
	if err != nil {
		return nil, errs.WrapError(errGettingImageFromUrl, err)
	}

	req.Header.Add("Authorization", c.basicAuth())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.WrapError(errGettingImageFromUrl, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapError(errGettingImageFromUrl, err)
	}

	imgEncoding := base64.StdEncoding.EncodeToString(body)
	img, err := c.textExtractor.ConvertToVisionImage(ctx, imgEncoding)
	if err != nil {
		return nil, errs.WrapError(errGettingImageFromUrl, err)
	}

	return img, nil
}

func (c *Client) basicAuth() string {
	auth := c.twilioCreds.Username + ":" + c.twilioCreds.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

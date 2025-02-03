package service

import (
	"context"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	visionApi "cloud.google.com/go/vision/apiv1"
	vision "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	maxResults = 20
	pngFormat  = "png"
	jpegFormat = "jpeg"

	errCreatingVisionClient = "failed to get client"
	errExtractingText       = "failed to extract text"
	errConvertingImage      = "failed to convert image"
	errNoTextExtracted      = "no text extracted from image"
)

type TextExtractor struct {
	logger logger.Logger
	client *visionApi.ImageAnnotatorClient
}

func NewTextExtractor(logger logger.Logger, client *visionApi.ImageAnnotatorClient) *TextExtractor {
	return &TextExtractor{
		logger: logger,
		client: client,
	}
}

// ExtractTextFromImage uses gcloud vision api to extract text from provided image
func (s *TextExtractor) ExtractTextFromImage(ctx context.Context, image *vision.Image) (string, error) {

	extractedText, err := s.client.DetectTexts(ctx, image, nil, maxResults)
	if err != nil {
		return "", errs.WrapError(errExtractingText, err)
	}
	if len(extractedText) == 0 {
		return "", errs.WrapError(errExtractingText, errors.New(errNoTextExtracted))
	}

	return extractedText[0].Description, nil
}

// ConvertToVisionImage converts a b64 encoded image to a type vision.Image
func (s *TextExtractor) ConvertToVisionImage(ctx context.Context, b64Str string) (*vision.Image, error) {

	imageData, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return nil, errs.WrapError(errConvertingImage, err)
	}

	reader := strings.NewReader(string(imageData))
	image, format, err := image.Decode(reader)
	if err != nil {
		return nil, errs.WrapError(errConvertingImage, err)
	}

	var imageBytes []byte
	switch format {
	case pngFormat:
		imageBytes, err = encodeToPNG(image)
	case jpegFormat:
		imageBytes, err = encodeToJPEG(image)
	default:
		return nil, errs.WrapError(errConvertingImage, errors.New("unsupported image format"))
	}

	if err != nil {
		return nil, errs.WrapError(errConvertingImage, err)
	}

	convertedImage, err := visionApi.NewImageFromReader(strings.NewReader(string(imageBytes)))
	if err != nil {
		return nil, errs.WrapError(errConvertingImage, err)
	}

	return convertedImage, nil
}

func encodeToJPEG(img image.Image) ([]byte, error) {
	buf := new(strings.Builder)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func encodeToPNG(img image.Image) ([]byte, error) {
	buf := new(strings.Builder)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

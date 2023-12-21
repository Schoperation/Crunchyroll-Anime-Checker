package core

import "encoding/base64"

// ImageEncoder is a domain service used to encode image data to some format.
type ImageEncoder struct{}

func NewImageEncoder() ImageEncoder {
	return ImageEncoder{}
}

func (encoder ImageEncoder) Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

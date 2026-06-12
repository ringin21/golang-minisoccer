package constants

import "net/textproto"

var (
	XServiceName = textproto.CanonicalMIMEHeaderKey(s: "x-service-name")
	XApiKey = textproto.CanonicalMIMEHeaderKey(s: "x-api-key")
	XRequestAt = textproto.CanonicalMIMEHeaderKey(s: "x-request-at")
	Authorization = textproto.CanonicalMIMEHeaderKey(s: "authorization")
)
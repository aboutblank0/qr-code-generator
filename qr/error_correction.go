package qr

type ErrorCorrectionLevel int

const (
	EC_L ErrorCorrectionLevel = iota
	EC_M
	EC_Q
	EC_H
)

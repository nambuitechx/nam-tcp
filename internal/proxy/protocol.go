package proxy

import "time"

const (
	ResponseOK  = "OK\n"
	ResponseERR = "ERR\n"
)

const AuthReadTimeout = 5 * time.Second

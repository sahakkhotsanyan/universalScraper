package config

import "os"

var DEBUG bool = false

func IsDebug() bool {
	if len(os.Args) > 1 {
		if os.Args[1] == "--debug" {
			DEBUG = true
		}
	}
	return DEBUG
}

const DebugID = 123456789
const StandardID = 123456789
const Token = "123456789:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

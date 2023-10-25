package utils

import (
	"os"
)

const TRANSKRIBUS_REST_API_URL = "https://transkribus.eu/TrpServer/rest"
const TRANSKRIBUS_AUTH_API_URL = "https://account.readcoop.eu/auth/realms/readcoop/protocol/openid-connect/token"

func OpenLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

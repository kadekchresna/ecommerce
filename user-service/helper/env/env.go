package env

import (
	"os"
	"strconv"
)

func GetEnvInt(s string) int {
	n, err := strconv.Atoi(os.Getenv(s))
	if err != nil {
		return 0
	}

	return n
}

package env

import (

	"os"

)

func GetEnv() string {
		envVar := os.Getenv("IFACES")

		if envVar == "" {
			envVar = "eth0"
		}

		return envVar

	}

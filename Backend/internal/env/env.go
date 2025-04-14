package env

import (

	"os"

)

func GetEnv() string {
		envVar := os.Getenv("IFACES")

		if envVar == "" {
			envVar = "wlp0s20f3"
		}

		return envVar

	}

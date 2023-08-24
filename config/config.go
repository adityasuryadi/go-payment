package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config interface {
	Get(key string) string
}

type configImpl struct {
}

func (config *configImpl) Get(key string) string {
	return os.Getenv(key)
}

func New(filenames ...string) Config {
	const projectDirName = "payment"
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + filenames[0])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cause": err,
			"cwd":   cwd,
		}).Fatal("Problem loading .env file")

		os.Exit(-1)
	}
	return &configImpl{}
}

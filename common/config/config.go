package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	if err := initViper(); err != nil {
		logrus.Panicln(err)
	}
}

func initViper() error {
	configPath, err := getRelativePathFromCaller()
	if err != nil {
		return fmt.Errorf("initViper | Failed to get relative path from caller | %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("initViper | Failed to read config | %v", err)
	}

	return nil
}

// getRelativePathFromCaller 获得从主调模块到配置文件的相对路径
func getRelativePathFromCaller() (string, error) {
	callerDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, here, _, _ := runtime.Caller(0)
	relativePath, err := filepath.Rel(callerDir, filepath.Dir(here))
	if err != nil {
		return "", err
	}

	return relativePath, nil
}

package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewViperConfig() error {
	// 加载 .env 文件（从根目录）
	_ = godotenv.Load("../../.env")

	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 替换配置值中的环境变量引用（支持 $VAR 和 ${VAR} 格式）
	replaceEnvVars(viper.AllSettings())

	return nil
}

func replaceEnvVars(settings map[string]interface{}) {
	for key, value := range settings {
		switch v := value.(type) {
		case string:
			// 替换 $VAR 或 ${VAR} 格式的环境变量
			expanded := os.ExpandEnv(v)
			viper.Set(key, expanded)
		case map[string]interface{}:
			// 递归处理嵌套的 map
			replaceEnvVars(v)
		}
	}
}

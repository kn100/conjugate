package configmanager

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ConfigManager struct {
	configFor string
	viper     *viper.Viper
}

func (cfgm *ConfigManager) Init(conjugatorName string) error {
	cfgm.configFor = conjugatorName

	hp, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	root := filepath.Join(hp, ".config", "conjugate", conjugatorName)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		err := os.MkdirAll(root, 0700)
		if err != nil {
			return err
		}
	}

	cfgm.viper = viper.New()
	cfgm.viper.SetConfigName("conjconfig")
	cfgm.viper.SetConfigType("json")
	cfgm.viper.AddConfigPath(root)
	cfgm.viper.SafeWriteConfig()
	return cfgm.viper.ReadInConfig()
}

func (cfgm *ConfigManager) Set(key string, value string) error {
	cfgm.viper.Set(key, value)
	return cfgm.viper.WriteConfig()
}

func (cfgm *ConfigManager) Get(key string) string {
	return cfgm.viper.GetString(key)
}

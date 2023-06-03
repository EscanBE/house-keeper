package config

import (
	"fmt"
	logtypes "github.com/EscanBE/go-lib/logging/types"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
)

// ApplicationConfig is the structure representation of configuration from `config.yaml` file
type ApplicationConfig struct {
	Logging        logtypes.LoggingConfig    `mapstructure:"logging"`
	SecretConfig   ApplicationSecretConfig   `mapstructure:"secrets"`
	TelegramConfig ApplicationTelegramConfig `mapstructure:"telegram"`
}

// ApplicationSecretConfig is the structure representation of configuration from `config.yaml` file, at `secret` section.
// Secret keys, tokens,... can be putted here
type ApplicationSecretConfig struct {
	TelegramToken string `mapstructure:"telegram-token"`
}

// ApplicationTelegramConfig is the structure representation of configuration from `config.yaml` file, at `telegram` section.
// It holds configuration of Telegram bot
type ApplicationTelegramConfig struct {
	LogChannelID int64 `mapstructure:"log-channel-id"`
	ErrChannelID int64 `mapstructure:"error-channel-id"`
}

// LoadConfig load the configuration from `config.yaml` file within the specified application's home directory
func LoadConfig(homeDir string) (*ApplicationConfig, error) {
	cfgFile := path.Join(homeDir, constants.DEFAULT_CONFIG_FILE_NAME)

	fileStats, err := os.Stat(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("conf file %s could not be found", cfgFile)
		}

		return nil, err
	}

	if fileStats.Mode().Perm() != constants.FILE_PERMISSION && fileStats.Mode().Perm() != 0o700 {
		//goland:noinspection GoBoolExpressions
		if constants.FILE_PERMISSION == 0o700 {
			panic(fmt.Errorf("incorrect permission of %s, must be %s", constants.DEFAULT_CONFIG_FILE_NAME, constants.FILE_PERMISSION_STR))
		} else {
			panic(fmt.Errorf("incorrect permission of %s, must be %s or 700", constants.DEFAULT_CONFIG_FILE_NAME, constants.FILE_PERMISSION_STR))
		}
	}

	viper.SetConfigType(constants.CONFIG_TYPE)
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "unable to read conf file")
	}

	conf := &ApplicationConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, errors.Wrap(err, "unable to deserialize conf file")
	}

	return conf, nil
}

// PrintOptions prints the configuration in the `config.yaml` in a nice way, human-readable
func (c ApplicationConfig) PrintOptions() {
	headerPrintln("- Tokens configuration:")
	if len(c.SecretConfig.TelegramToken) > 0 {
		headerPrintln("  + Telegram bot token has set")

		if len(c.SecretConfig.TelegramToken) > 0 {
			if c.TelegramConfig.LogChannelID != 0 {
				headerPrintf("  + Telegram log channel ID: %s\n", c.TelegramConfig.LogChannelID)
			} else {
				headerPrintln("  + Missing configuration for log channel ID")
			}
			if c.TelegramConfig.ErrChannelID != 0 {
				headerPrintf("  + Telegram error channel ID: %s\n", c.TelegramConfig.ErrChannelID)
			} else {
				headerPrintln("  + Missing configuration for error channel ID")
			}
		}
	} else {
		headerPrintln("  + Telegram function was disabled because token has not been set")
	}

	headerPrintln("- Logging:")
	if len(c.Logging.Level) < 1 {
		headerPrintf("  + Level: %s\n", logtypes.LOG_LEVEL_DEFAULT)
	} else {
		headerPrintf("  + Level: %s\n", c.Logging.Level)
	}

	if len(c.Logging.Format) < 1 {
		headerPrintf("  + Format: %s\n", logtypes.LOG_FORMAT_DEFAULT)
	} else {
		headerPrintf("  + Format: %s\n", c.Logging.Format)
	}
}

// headerPrintf prints text with prefix
func headerPrintf(format string, a ...any) {
	fmt.Printf("[HCFG]"+format, a...)
}

// headerPrintln prints text with prefix
func headerPrintln(a string) {
	fmt.Println("[HCFG]" + a)
}

// Validate performs validation on the configuration specified in the `config.yaml` within application's home directory
func (c ApplicationConfig) Validate() error {
	if len(c.SecretConfig.TelegramToken) > 0 {
		if c.TelegramConfig.LogChannelID == 0 {
			return fmt.Errorf("missing telegram log channel ID")
		}
		if c.TelegramConfig.ErrChannelID == 0 {
			return fmt.Errorf("missing telegram error channel ID")
		}
	}

	// validate Logging section
	errLogCfg := c.Logging.Validate()
	if errLogCfg != nil {
		return errLogCfg
	}

	return nil
}

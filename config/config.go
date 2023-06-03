package config

import (
	"fmt"
	"github.com/EscanBE/go-app-name/constants"
	libdbtypes "github.com/EscanBE/go-lib/database/types"
	logtypes "github.com/EscanBE/go-lib/logging/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
)

// ApplicationConfig is the structure representation of configuration from `config.yaml` file
type ApplicationConfig struct {
	Logging        logtypes.LoggingConfig     `mapstructure:"logging"`
	WorkerConfig   ApplicationWorkerConfig    `mapstructure:"worker"`
	SecretConfig   ApplicationSecretConfig    `mapstructure:"secrets"`
	Endpoints      ApplicationEndpointsConfig `mapstructure:"endpoints"`
	TelegramConfig ApplicationTelegramConfig  `mapstructure:"telegram"`

	// TODO Sample of template
	// TODO remove this sample field
	SampleDictionary map[string]SampleDictionaryValueConfig `mapstructure:"sample-dict"`

	// TODO Sample of template
	// TODO remove this sample field
	SampleArray []int `mapstructure:"sample-array"`
}

// ApplicationWorkerConfig is the structure representation of configuration from `config.yaml` file, at `worker` section.
// It holds configuration related to how the process would work, like how many parallel go routines,...
type ApplicationWorkerConfig struct {
	// Count is the number of worker you want to launch in parallel
	Count int `mapstructure:"count"`
}

// ApplicationSecretConfig is the structure representation of configuration from `config.yaml` file, at `secret` section.
// Secret keys, tokens,... can be putted here
type ApplicationSecretConfig struct {
	TelegramToken string `mapstructure:"telegram-token"`

	// TODO Sample of template
	// TODO remove this sample field
	SampleToken1 string `mapstructure:"sample-token1"`

	// TODO Sample of template
	// TODO remove this sample field
	SampleToken2 string `mapstructure:"sample-token2"`
}

// ApplicationEndpointsConfig holds nested configurations relates to remote endpoints
type ApplicationEndpointsConfig struct {
	Database libdbtypes.PostgresDatabaseConfig `mapstructure:"db"`
}

// ApplicationTelegramConfig is the structure representation of configuration from `config.yaml` file, at `telegram` section.
// It holds configuration of Telegram bot
type ApplicationTelegramConfig struct {
	LogChannelID int64 `mapstructure:"log-channel-id"`
	ErrChannelID int64 `mapstructure:"error-channel-id"`
}

// SampleDictionaryValueConfig is a sample struct
// TODO Sample of template
// TODO remove this sample struct
type SampleDictionaryValueConfig struct {
	SampleValueElement1 string `mapstructure:"sample-value-element-1"`
	SampleValueElement2 string `mapstructure:"sample-value-element-2"`
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
	// TODO Sample of template
	// TODO BEGIN printing sample config
	// TODO remove the following sample print
	headerPrintf("  + Sample token 1 has set: %t\n", len(c.SecretConfig.SampleToken1) > 0)
	if len(c.SecretConfig.SampleToken2) > 0 {
		headerPrintln("  + Sample token 2 has set")
	} else {
		headerPrintln("  + Sample token 2 has not been set")
	}

	headerPrintf("- Sample array size: %d\n", len(c.SampleArray))
	headerPrintf("  + Values: [")
	for i, v := range c.SampleArray {
		if i > 0 {
			fmt.Printf(",")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println("]")
	// TODO END printing sample config

	// TODO Sample of template
	// TODO show configuration here

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

	headerPrintln("- Worker's behavior:")
	headerPrintf("  + Number of workers: %d\n", c.WorkerConfig.Count)

	headerPrintln("- Database:")
	headerPrintf("  + Host: %s\n", c.Endpoints.Database.Host)
	headerPrintf("  + Port: %d\n", c.Endpoints.Database.Port)
	headerPrintf("  + Username: %s\n", c.Endpoints.Database.Username)
	headerPrintf("  + DB name: %s\n", c.Endpoints.Database.Name)
	headerPrintf("  + Schema name: %s\n", c.Endpoints.Database.Schema)
	headerPrintf("  + Enable SSL: %t\n", c.Endpoints.Database.EnableSsl)
	headerPrintf("  + Max open connections: %d\n", c.Endpoints.Database.MaxOpenConnectionCount)
	headerPrintf("  + Max idle connections: %d\n", c.Endpoints.Database.MaxIdleConnectionCount)
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

	// TODO Sample of template
	// TODO validate application configuration here, return error if any problem

	// validate Logging section
	errLogCfg := c.Logging.Validate()
	if errLogCfg != nil {
		return errLogCfg
	}

	// validator Worker section
	if c.WorkerConfig.Count < 1 {
		return fmt.Errorf("number of worker can not be lower than 1")
	}

	// Validate Endpoints-DB section
	errDbCfg := c.Endpoints.Database.Validate()
	if errDbCfg != nil {
		return errDbCfg
	}

	return nil
}

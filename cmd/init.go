package cmd

import (
	"fmt"
	cmdutils "github.com/EscanBE/go-app-name/cmd/utils"
	"github.com/EscanBE/go-app-name/constants"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// initCmd represents the init command, to be used to generate home directory with configuration file named `config.yaml`
var initCmd = &cobra.Command{
	Use:   "init",
	Short: fmt.Sprintf("Init home directory & configuration files for %s at %s", constants.APP_NAME, cmdutils.GetDefaultHomeDirectory()),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(homeDir)

		if err != nil && os.IsNotExist(err) {
			fmt.Printf("Require home dir '%s' does not exists, going to create new home dir\n", homeDir)
			err := os.Mkdir(homeDir, 0o750)
			libutils.ExitIfErr(err, fmt.Sprintf("Unable to create home dir %s", homeDir))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		cfgFile := path.Join(homeDir, constants.DEFAULT_CONFIG_FILE_NAME)

		_, err = os.Stat(cfgFile)
		if err != nil && os.IsNotExist(err) {
			fmt.Printf("Config file '%s' does not exists, going to create new file with permission %s\n", cfgFile, constants.FILE_PERMISSION_STR)
			file, err := os.Create(cfgFile)
			libutils.ExitIfErr(err, fmt.Sprintf("Unable to create config file %s", cfgFile))
			err = file.Chmod(constants.FILE_PERMISSION)
			libutils.ExitIfErr(err, fmt.Sprintf("Unable to set permission for new config file %s to %s", cfgFile, constants.FILE_PERMISSION_STR))
			_, err = file.WriteString(
				// TODO Sample of template
				// trailing style: 2 spaces
				fmt.Sprintf(`# %s's configuration file
logging:
  level: info # debug || info || error
  format: text # text || json
worker: # configure worker's behavior or remove this section if no need worker concept
  count: 1 # number of parallel workers
secrets:
  telegram-token: # leave it empty to disable telegram, but it will crash if you invoke function to send message
  sample-token1: XXX
  sample-token2: XXX
endpoints:
  db: # database connection
    name: postgres
    host: 127.0.0.1
    port: 5432
    username: postgres
    password: 1234567
    schema: public
    enable-ssl: false
    max-open-connection-count: 20
    max-idle-connection-count: 20
telegram:
  log-channel-id: -1001738777107
  error-channel-id: -1001500282644
sample-dict:
  sample-key-1:
    sample-value-element-1: XXX
    sample-value-element-2: XXX
  sample-key-2:
    sample-value-element-1: XXX
    sample-value-element-2: XXX
sample-array: [1,2,3,4,5]
`, constants.APP_NAME))
			libutils.ExitIfErr(err, fmt.Sprintf("Unable to write content for new config file %s", cfgFile))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

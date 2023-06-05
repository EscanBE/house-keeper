package files

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var defaultRsyncOptions = []string{"--human-readable", "--compress", "--progress", "--stats"}

// RsyncCommands registers a sub-tree of commands
func RsyncCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rsync [src] [dest]",
		Short: fmt.Sprintf("Remotely/Locally transfer file: %s rsync src dest", constants.BINARY_NAME),
		Long: fmt.Sprintf(`Remotely/Locally transfer file using rsync.
- Send file:
> %s rsync '/var/log/*.log' 'backup@192.168.0.2:/mnt/md0/backup/logs'
- Receive file:
> %s rsync 'load-balancer:/var/log/*.log' '/mnt/md0/backup/logs'

Note:
- This is just a wrapper of rsync, you must know how to use rsync and got rsync installed in order to use this.
  Actual translated rsync command would look similar to:
  > /usr/bin/rsync -hz --progress --stats -e ssh "server:/var/logs/*.log" "/mnt/md0/backup/logs"
- When transfer from/to remote server, you must connect to that remote server at least one time before to perform host key verification (one time action) because the transfer will be performed via ssh.
`, constants.BINARY_NAME, constants.BINARY_NAME),
		Args: cobra.ExactArgs(2),
		Run:  remoteTransferFile,
	}

	cmd.PersistentFlags().Bool(
		constants.FLAG_REMOTE_TO_LOCAL,
		false,
		"ensure the transfer direction is from remote server to local",
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_LOCAL_TO_REMOTE,
		false,
		"ensure the transfer direction is from local to remote server",
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_LOCAL_TO_LOCAL,
		false,
		"ensure the transfer direction is from local to local",
	)

	cmd.PersistentFlags().StringSlice(
		constants.FLAG_TOOL_OPTIONS,
		defaultRsyncOptions,
		"supply options passes to rsync, comma separated",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_TOOL_FILE,
		"",
		"custom rsync file path (absolute)",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_PASSWORD_FILE,
		"",
		"file path which store password to access remote server",
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_NO_PASSWORD,
		false,
		"force connect remote server without password (when remote user does not have password or identity key does not protected by password)",
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_SSHPASS_PASSPHRASE,
		false,
		"by default sshpass (if sshpass exists) passes password. If you are authenticating using passphrase, program will be hang (search phrase not found), supply this flag to indicate and would fix it",
	)

	return cmd
}

func remoteTransferFile(cmd *cobra.Command, args []string) {
	src := strings.TrimSpace(args[0])
	if len(src) < 1 {
		panic("source file/dir is empty")
	}

	dest := strings.TrimSpace(args[1])
	if len(dest) < 1 {
		panic("destination file/dir is empty")
	}

	isSrcRemote := strings.Contains(src, ":")
	isDestRemote := strings.Contains(dest, ":")

	if isSrcRemote && isDestRemote {
		panic("not support transfer direction from remote to remote")
	} else if isSrcRemote && !isDestRemote {
		confirm, _ := cmd.Flags().GetBool(constants.FLAG_REMOTE_TO_LOCAL)
		if !confirm {
			panic(fmt.Errorf("detected transfer direction is from remote to local so flag --%s is required to confirm", constants.FLAG_REMOTE_TO_LOCAL))
		}
	} else if !isSrcRemote && isDestRemote {
		confirm, _ := cmd.Flags().GetBool(constants.FLAG_LOCAL_TO_REMOTE)
		if !confirm {
			panic(fmt.Errorf("detected transfer direction is from local to remote so flag --%s is required to confirm", constants.FLAG_LOCAL_TO_REMOTE))
		}
	} else if !isSrcRemote && !isDestRemote {
		confirm, _ := cmd.Flags().GetBool(constants.FLAG_LOCAL_TO_LOCAL)
		if !confirm {
			panic(fmt.Errorf("detected local transfer so flag --%s is required to confirm", constants.FLAG_LOCAL_TO_LOCAL))
		}
	}

	if !isSrcRemote {
		_, err := os.Stat(src)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("local source file/dir does not exists: %s", src))
		}
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("problem while checking local source file %s", src)))
		}
	}

	if !isDestRemote {
		fi, err := os.Stat(dest)
		if err == nil {
			if fi.IsDir() {
				// ok
			} else {
				panic(fmt.Errorf("local destination file/dir already exists: %s", dest))
			}
		} else {
			if os.IsNotExist(err) {
				// ok
			} else {
				panic(errors.Wrap(err, fmt.Sprintf("problem while checking local destination file %s", dest)))
			}
		}
	}

	toolName := "rsync"
	customToolName, _ := cmd.Flags().GetString(constants.FLAG_TOOL_FILE)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err := os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
		}

		toolName = customToolName
	}

	options, _ := cmd.Flags().GetStringSlice(constants.FLAG_TOOL_OPTIONS)
	if len(options) < 1 {
		options = defaultRsyncOptions
	}

	if !isSrcRemote && !isDestRemote {
		run(toolName, append(options, src, dest))
		return
	}

	noPassword, _ := cmd.Flags().GetBool(constants.FLAG_NO_PASSWORD)
	if noPassword {
		run(toolName, append(options, "-e", "ssh", src, dest))
		return
	}

	sshPassPhrase, _ := cmd.Flags().GetBool(constants.FLAG_SSHPASS_PASSPHRASE)

	passwordFile, _ := cmd.Flags().GetString(constants.FLAG_PASSWORD_FILE)
	if len(passwordFile) > 0 {
		fip, err := os.Stat(passwordFile)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("supplied password file does not exists: %s", passwordFile))
		}

		bz, err := os.ReadFile(passwordFile)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("failed to read password file: %s", passwordFile)))
		}
		password := strings.TrimSpace(string(bz))
		if len(password) < 1 {
			panic(fmt.Errorf("password file is empty: %s", passwordFile))
		}

		if fip.Mode().Perm() != constants.FILE_PERMISSION_400 &&
			fip.Mode().Perm() != constants.FILE_PERMISSION_600 &&
			fip.Mode().Perm() != constants.FILE_PERMISSION_700 {
			//goland:noinspection GoBoolExpressions
			panic(fmt.Errorf("incorrect permission of password file, must be %s or %s or %s", constants.FILE_PERMISSION_400_STR, constants.FILE_PERMISSION_600_STR, constants.FILE_PERMISSION_700_STR))
		}

		if utils.HasToolSshPass() {
			fmt.Println("Using sshpass to passing password file")

			var cmdArgs []string
			if sshPassPhrase {
				//goland:noinspection SpellCheckingInspection
				cmdArgs = []string{"-P", "assphrase", "-f", passwordFile}
			} else {
				cmdArgs = []string{"-f", passwordFile}
			}
			cmdArgs = append(cmdArgs, toolName)
			cmdArgs = append(cmdArgs, options...)
			cmdArgs = append(cmdArgs, "-e", "ssh", src, dest)

			run("sshpass", cmdArgs)
			return
		}

		fmt.Println("Using environment variable", constants.ENV_RSYNC_PASSWORD, "to passing password from password file to rsync")
		fmt.Println("**WARNING: if remote machine does not have rsync service running, password prompt still appears")
		run(toolName, append(options, "-e", "ssh", src, dest), fmt.Sprintf("%s=%s", constants.ENV_RSYNC_PASSWORD, password))
		return
	}

	rsyncPassword := strings.TrimSpace(os.Getenv(constants.ENV_RSYNC_PASSWORD))
	sshPassword := strings.TrimSpace(os.Getenv(constants.ENV_SSHPASS))

	var password string
	if len(rsyncPassword) > 0 && len(sshPassword) > 0 {
		if rsyncPassword != sshPassword {
			panic(fmt.Errorf("both environment variables %s and %s are set but mis-match, consider remove one to take the rest", constants.ENV_RSYNC_PASSWORD, constants.ENV_SSHPASS))
		}

		password = rsyncPassword
	} else if len(rsyncPassword) > 0 {
		password = rsyncPassword
	} else if len(sshPassword) > 0 {
		password = sshPassword
	} else {
		panic(fmt.Errorf("missing password for remote server, either environment variable %s or %s or flag --%s is required", constants.ENV_RSYNC_PASSWORD, constants.ENV_SSHPASS, constants.FLAG_PASSWORD_FILE))
	}

	if utils.HasToolSshPass() {
		if len(rsyncPassword) > 0 && len(sshPassword) == 0 {
			fmt.Println("Copied environment variable value from", constants.ENV_RSYNC_PASSWORD, "to", constants.ENV_SSHPASS)
		}
		fmt.Println("Using sshpass to passing password via environment variable", constants.ENV_SSHPASS)

		var cmdArgs []string
		if sshPassPhrase {
			//goland:noinspection SpellCheckingInspection
			cmdArgs = []string{"-P", "assphrase", "-e", toolName}
		} else {
			cmdArgs = []string{"-e", toolName}
		}
		cmdArgs = append(cmdArgs, options...)
		cmdArgs = append(cmdArgs, "-e", "ssh", src, dest)

		run("sshpass", cmdArgs, fmt.Sprintf("%s=%s", constants.ENV_SSHPASS, password))
		return
	}

	if len(sshPassword) > 0 && len(rsyncPassword) == 0 {
		fmt.Println("Copied environment variable value from", constants.ENV_SSHPASS, "to", constants.ENV_RSYNC_PASSWORD)
	}
	fmt.Println("Using environment variable", constants.ENV_RSYNC_PASSWORD, "to passing password to rsync")
	fmt.Println("**WARNING: if remote machine does not have rsync service running, password prompt still appears")
	run(toolName, append(options, "-e", "ssh", src, dest), fmt.Sprintf("%s=%s", constants.ENV_RSYNC_PASSWORD, password))
}

func run(toolName string, args []string, additionalEnvVars ...string) {
	rsyncCmd := exec.Command(toolName, args...)

	rsyncCmd.Env = append(additionalEnvVars, additionalEnvVars...)
	// stdin, _ := rsyncCmd.StdinPipe()
	stdout, _ := rsyncCmd.StdoutPipe()
	stderr, _ := rsyncCmd.StderrPipe()
	rsyncStdOutScanner := bufio.NewScanner(stdout)
	rsyncStdErrScanner := bufio.NewScanner(stderr)
	err := rsyncCmd.Start()
	if err != nil {
		fmt.Println("problem when starting", toolName, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			oScan := rsyncStdOutScanner.Scan()
			eScan := rsyncStdErrScanner.Scan()
			if oScan {
				fmt.Println("INF:", rsyncStdOutScanner.Text())
			}
			if eScan {
				fmt.Println("ERR:", rsyncStdErrScanner.Text())
			}
			if !oScan && !eScan {
				break
			}
		}
		err = rsyncCmd.Wait()
		if err != nil {
			fmt.Println("problem when waiting process", err)
		}
		defer wg.Done()
	}()

	wg.Wait()
}

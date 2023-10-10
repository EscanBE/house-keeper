package cmd

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	flagConfirmExecution = "yes"
)

var predefinedAliases map[string]predefinedAlias
var longestUseDesc int

/*
Sample content for alias file .hkd_alias:
echo "say-hello	echo \"Hello World\"" >> ~/.hkd_alias
hkd a say-hello
*/

// aliasCmd represents the 'a' command, it executes commands based on pre-defined input alias
var aliasCmd = &cobra.Command{
	Use:     "a [alias]",
	Aliases: []string{"alias"},
	Short:   "Execute commands based on pre-defined alias",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		predefinedAliases = make(map[string]predefinedAlias)

		registerStartupPredefinedAliases()
		registerPredefinedAliasesFromFile()

		if len(args) < 1 {
			lineFormat := " %-" + fmt.Sprintf("%d", longestUseDesc+1) + "s: %s\n"

			fmt.Println("Registered aliases:")
			for _, alias := range goe.NewIEnumerable[string](libutils.GetKeys(predefinedAliases)...).Order().GetOrderedEnumerable().ToArray() {
				pa := predefinedAliases[alias]
				if pa.overridden {
					fmt.Printf(" *overriden*")
				}
				fmt.Printf(lineFormat, pa.use, strings.Join(pa.command, " "))
			}
			fmt.Printf("Alias can be customized by adding into ~/%s (TSV format with each line content \"<alias><tab><command>\")\n", constants.PREDEFINED_ALIAS_FILE_NAME)
			return
		}

		selectedAlias := args[0]
		pa, found := predefinedAliases[selectedAlias]
		if !found {
			fmt.Printf("Alias '%s' has not been registered before\n", selectedAlias)
			os.Exit(1)
		}

		command := pa.command
		if len(args) > 1 && pa.alter != nil {
			command = (*pa.alter)(command, args[1:])
		}

		if len(command) < 1 {
			panic("empty command")
		}

		joinedCommand := strings.Join(command, " ")

		confirmExecution, _ := cmd.Flags().GetBool(flagConfirmExecution)

		if confirmExecution {
			const waitingTime = 10
			fmt.Println("Pending execution command:")
			fmt.Printf("> %s\n", joinedCommand)
			fmt.Printf("(actual command: [/bin/bash] [-c] [%s])\n", joinedCommand)
			fmt.Printf("Executing in %d seconds...\n", waitingTime)
			time.Sleep(waitingTime * time.Second)
		} else {
			fmt.Println("Are you sure want to execute the following command?")
			fmt.Printf("> %s\n", joinedCommand)
			fmt.Printf("(actual command: [/bin/bash] [-c] [%s])\n", joinedCommand)
			fmt.Println("Yes/No?")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(strings.ToLower(text))

			switch text {
			case "y":
				break
			case "yes":
				break
			case "n":
				fmt.Println("Aborted")
				os.Exit(1)
			case "no":
				fmt.Println("Aborted")
				os.Exit(1)
			default:
				fmt.Printf("Aborted! '%s' is not an accepted answer!\n", text)
				fmt.Println("Your answer must be Yes/No (or Y/N)")
				os.Exit(1)
			}
		}

		fmt.Println("Executing...")

		ec := utils.LaunchAppWithDirectStd("/bin/bash", []string{"-c", joinedCommand}, nil)
		if ec != 0 {
			os.Exit(ec)
		}
	},
}

//goland:noinspection SpellCheckingInspection
func registerStartupPredefinedAliases() {
	currentUser, errGetUser := user.Current()
	if errGetUser != nil {
		libutils.PrintlnStdErr("ERR: failed to get current user:", errGetUser.Error())
		os.Exit(1)
	}
	home, errGetUserHomeDir := os.UserHomeDir()
	if errGetUserHomeDir != nil {
		libutils.PrintlnStdErr("ERR: failed to get home directory:", errGetUserHomeDir.Error())
		os.Exit(1)
	}
	isUserRoot := currentUser.Username == "root"

	registerPredefinedAliasForNode := func(binaryName, prefix string) {
		hasBinary := utils.HasBinaryName(binaryName)
		if hasBinary || isExistsServiceFile(binaryName) {
			registerPredefinedAlias(fmt.Sprintf("%srs", prefix), []string{"sudo", "systemctl", "restart", binaryName}, nil)
			registerPredefinedAlias(fmt.Sprintf("%sstop", prefix), []string{"sudo", "systemctl", "stop", binaryName}, nil)
			registerPredefinedAlias(fmt.Sprintf("%sl [?since]", prefix), []string{"sudo", "journalctl", "-fu", binaryName}, &genericAlterJournalctl)
		}
		if !isUserRoot && hasBinary {
			nodeHome := path.Join(home, "."+binaryName)
			if _, err := os.Stat(nodeHome); err == nil {
				homeActualPath, err := utils.TryReadSymlink(nodeHome)
				if err != nil { // blind accept
					homeActualPath = nodeHome
				}
				nodeData := path.Join(homeActualPath, "data")
				dataActualPath, err := utils.TryReadSymlink(nodeData)
				if err != nil { // blind accept
					dataActualPath = nodeData
				}

				const notSupportCommandResetIfDataSizeGETerrabyte = 1.0
				const notSupportCommandResetIfDataSizeGEBytes = int64(notSupportCommandResetIfDataSizeGETerrabyte * 1_000_000_000_000)

				totalSize, err := utils.SumDirectorySize(dataActualPath, notSupportCommandResetIfDataSizeGEBytes)
				if err == nil || utils.IsErrorLimitSumDirectorySizeReached(err) {
					resetAlias := fmt.Sprintf("%sreset", prefix)
					if totalSize >= notSupportCommandResetIfDataSizeGEBytes {
						libutils.PrintfStdErr("WARN: %s is not supported for node with data size >= %.2f TB\n", resetAlias, notSupportCommandResetIfDataSizeGETerrabyte)
					} else {
						registerPredefinedAlias(resetAlias, []string{binaryName, "tendermint", "unsafe-reset-all", "--home", nodeHome, "--keep-addr-book"}, nil)
					}
				} else {
					libutils.PrintfStdErr("WARN: failed to calculate total data size of %s: %s\n", dataActualPath, err.Error())
				}
			}
		}
	}

	// Manage Evmos nodes
	registerPredefinedAliasForNode("evmosd", "es")

	// Manage Dymension nodes
	registerPredefinedAliasForNode("dymd", "dym")

	// Manage Ethermint dev nodes
	registerPredefinedAliasForNode("ethermintd", "eth")

	// Manage CosmosHub nodes
	registerPredefinedAliasForNode("gaid", "ga")

	// Manage indexer
	if utils.HasBinaryName("crawld") || isExistsServiceFile("crawld") {
		registerPredefinedAlias("ecrs", []string{"sudo", "systemctl", "restart", "crawld"}, nil)
		registerPredefinedAlias("ecstop", []string{"sudo", "systemctl", "stop", "crawld"}, nil)
		registerPredefinedAlias("ecl [?since]", []string{"sudo", "journalctl", "-fu", "crawld"}, &genericAlterJournalctl)
	}

	// Manage proxy
	if utils.HasBinaryName("epod") || isExistsServiceFile("epod") {
		registerPredefinedAlias("eprs", []string{"sudo", "systemctl", "restart", "epod"}, nil)
		registerPredefinedAlias("epstop", []string{"sudo", "systemctl", "stop", "epod"}, nil)
		registerPredefinedAlias("epl [?since]", []string{"sudo", "journalctl", "-fu", "epod"}, &genericAlterJournalctl)
	}

	// Read logging
	registerPredefinedAlias("log [?service] [?since]", []string{"sudo", "journalctl"}, &aliasLogHandler)

	// Git
	if _, err := os.Stat(".git"); err == nil {
		registerPredefinedAlias("pull [?branch] [?branch2] [...]", []string{"git", "fetch", "--all", "&&", "git", "checkout", "main", "&&", "git", "pull"}, &gitPullHandler)
	}
}

var aliasLogHandler commandAlter = func(_, args []string) []string {
	service := args[0]
	command := []string{"sudo", "journalctl", "-fu", service}

	if len(args) > 1 {
		command = genericAlterJournalctl(command, args[1:])
	}

	return command
}

var genericAlterJournalctl commandAlter = func(command, args []string) []string {
	return append(command, "--since", "'"+strings.Join(args, " ")+"'")
}

var gitPullHandler commandAlter = func(_, args []string) []string {
	command := []string{"git", "fetch", "--all"}

	for _, branch := range args {
		command = append(command, "&&", "git", "checkout", branch, "&&", "git", "pull")
	}

	return command
}

func registerPredefinedAliasesFromFile() {
	home, errGetUserHomeDir := os.UserHomeDir()
	if errGetUserHomeDir != nil {
		fmt.Println("ERR: failed to get home directory:", errGetUserHomeDir.Error())
		return
	}

	aliasFile := path.Join(home, constants.PREDEFINED_ALIAS_FILE_NAME)

	file, errFile := os.Stat(aliasFile)
	if errFile != nil {
		if os.IsNotExist(errFile) {
			return
		}
		fmt.Printf("ERR: unable to check alias file %s: %s\n", aliasFile, errFile.Error())
		return
	}

	if file.IsDir() {
		return
	}

	bz, err := os.ReadFile(aliasFile)
	if err != nil {
		fmt.Printf("ERR: failed to read alias file %s: %s\n", aliasFile, err.Error())
		return
	}

	tsvLines := goe.NewIEnumerable(strings.Split(string(bz), "\n")...).Select(func(line string) any {
		return strings.TrimSpace(line)
	}).CastString()

	regexReplaceContinousSpace := regexp.MustCompile("[\\s\\t]+")

	for _, line := range tsvLines.ToArray() {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if libutils.IsBlank(line) {
			continue
		}

		spl := strings.Split(
			regexReplaceContinousSpace.ReplaceAllString(strings.Replace(line, "\t", " ", -1), " "),
			" ",
		)

		if len(spl) < 2 {
			panic(fmt.Errorf("malformed %s", constants.PREDEFINED_ALIAS_FILE_NAME))
		}

		alias := spl[0]
		command := spl[1:]
		if pa, found := predefinedAliases[alias]; found {
			pa.command = command
			pa.use = alias
			pa.overridden = true
			predefinedAliases[alias] = pa
		} else {
			registerPredefinedAlias(alias, command, nil)
		}
	}
}

func registerPredefinedAlias(use string, defaultCommand []string, alter *commandAlter) {
	spl := strings.Split(use, " ")
	alias := spl[0]
	predefinedAliases[alias] = predefinedAlias{
		alias:   alias,
		use:     use,
		command: defaultCommand,
		alter:   alter,
	}
	longestUseDesc = libutils.MaxInt(longestUseDesc, len(use))
}

func isExistsServiceFile(serviceName string) bool {
	file, err := os.Stat(path.Join("/etc/systemd/system", serviceName+".service"))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false // treat as not
	}
	return !file.IsDir()
}

func init() {
	aliasCmd.PersistentFlags().Bool(
		flagConfirmExecution,
		false,
		"skip confirmation before executing the command, but wait few seconds before executing",
	)

	rootCmd.AddCommand(aliasCmd)
}

type predefinedAlias struct {
	alias      string
	use        string
	command    []string
	alter      *commandAlter
	overridden bool
}

type commandAlter func(command, args []string) []string

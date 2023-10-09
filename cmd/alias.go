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
	"path"
	"regexp"
	"strings"
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

		fmt.Println("Executing...")

		var envVars []string
		envVars = append(envVars, "SYSTEMD_PAGER=cat") // disable systemd pager

		ec := utils.LaunchApp("/bin/bash", []string{"-c", joinedCommand}, envVars)

		if ec != 0 {
			fmt.Println("Exited with status code:", ec)
		}

		os.Exit(ec)
	},
}

//goland:noinspection SpellCheckingInspection
func registerStartupPredefinedAliases() {
	home, errGetUserHomeDir := os.UserHomeDir()

	// Manage Evmos nodes
	registerPredefinedAlias("esrs", []string{"sudo", "systemctl", "restart", "evmosd"}, nil)
	registerPredefinedAlias("esstop", []string{"sudo", "systemctl", "stop", "evmosd"}, nil)
	registerPredefinedAlias("esl [since]", []string{"sudo", "journalctl", "-fu", "evmosd"}, &genericAlterJournalctl)
	if errGetUserHomeDir != nil {
		fmt.Println("ERR: Failed to register predefined alias esreset")
	} else {
		registerPredefinedAlias("esreset", []string{"evmosd", "tendermint", "unsafe-reset-all", "--home", path.Join(home, ".evmosd"), "--keep-addr-book"}, nil)
	}

	// Manage indexer
	registerPredefinedAlias("ecrs", []string{"sudo", "systemctl", "restart", "crawld"}, nil)
	registerPredefinedAlias("ecstop", []string{"sudo", "systemctl", "stop", "crawld"}, nil)
	registerPredefinedAlias("ecl [since]", []string{"sudo", "journalctl", "-fu", "crawld"}, &genericAlterJournalctl)

	// Manage proxy
	registerPredefinedAlias("eprs", []string{"sudo", "systemctl", "restart", "epod"}, nil)
	registerPredefinedAlias("epstop", []string{"sudo", "systemctl", "stop", "epod"}, nil)
	registerPredefinedAlias("epl [since]", []string{"sudo", "journalctl", "-fu", "epod"}, &genericAlterJournalctl)

	// Read logging
	registerPredefinedAlias("log [service] [since]", []string{"sudo", "journalctl"}, &aliasLogHandler)

	// Git
	registerPredefinedAlias("pull [branch] [branch2] [...]", []string{"git", "fetch", "--all", "&&", "git", "checkout", "main", "&&", "git", "pull"}, &gitPullHandler)
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

func init() {
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

package gen

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/spf13/cobra"
	"os/user"
	"strings"
	"time"
)

// GenerateVisudoCommands registers a sub-tree of commands
func GenerateVisudoCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "visudo [service] [user_name]",
		Short: "Generate config for user to start/stop read log of service",
		Args:  cobra.RangeArgs(1, 2),
		Run:   generateVisudo,
	}

	return cmd
}

func generateVisudo(_ *cobra.Command, args []string) {
	service := args[0]
	var userName string
	if len(args) < 2 {
		u, err := user.Current()
		libutils.PanicIfErr(err, "failed to get current user")
		userName = u.Username
	} else {
		userName = args[1]
	}
	fmt.Println("Generating visudo config...")
	fmt.Println("Username:", userName)
	fmt.Println("Service :", service)
	time.Sleep(5 * time.Second)

	template := `# Allow user '$USR' to manage '$SVC' service
$USR ALL= NOPASSWD: /usr/bin/systemctl start $SVC
$USR ALL= NOPASSWD: /usr/bin/systemctl stop $SVC
$USR ALL= NOPASSWD: /usr/bin/systemctl restart $SVC
$USR ALL= NOPASSWD: /usr/bin/systemctl enable $SVC
$USR ALL= NOPASSWD: /usr/bin/systemctl disable $SVC
$USR ALL= NOPASSWD: /usr/bin/systemctl status $SVC
$USR ALL= NOPASSWD: /usr/bin/journalctl`
	if userName == "statd" {
		template += `

# Allow user '$USR' to get current firewall status
$USR ALL= NOPASSWD: /usr/sbin/ufw status`
	} else if userName == "bots" {
		template += `

# Allow user '$USR' to restart/shutdown server
$USR ALL= NOPASSWD: /usr/sbin/reboot
$USR ALL= NOPASSWD: /usr/sbin/shutdown now`
	}

	content := strings.ReplaceAll(template, "$USR", userName)
	content = strings.ReplaceAll(content, "$SVC", service)

	fmt.Println(content)
}

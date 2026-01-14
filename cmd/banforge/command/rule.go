package command

import (
	"fmt"
	"os"

	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/spf13/cobra"
)

var (
	name    string
	service string
	path    string
	status  string
	method  string
)

var RuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage rules",
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "CLI interface for add new rule to file /etc/banforge/rules.toml",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			fmt.Printf("Rule name can't be empty\n")
			os.Exit(1)
		}
		if service == "" && path == "" && status == "" && method == "" {
			fmt.Printf("At least 1 rule field must be filled in.")
			os.Exit(1)
		}

		err := config.NewRule(name, service, path, status, method)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Rule added successfully!")
	},
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := config.LoadRuleConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, rule := range r {
			fmt.Printf("Name: %s\nService: %s\nPath: %s\nStatus: %s\nMethod: %s\n\n", rule.Name, rule.ServiceName, rule.Path, rule.Status, rule.Method)
		}
	},
}

func RuleRegister() {
	RuleCmd.AddCommand(AddCmd)
	RuleCmd.AddCommand(ListCmd)
	AddCmd.Flags().StringVarP(&name, "name", "n", "", "rule name (required)")
	AddCmd.Flags().StringVarP(&service, "service", "s", "", "service name")
	AddCmd.Flags().StringVarP(&path, "path", "p", "", "request path")
	AddCmd.Flags().StringVarP(&status, "status", "c", "", "HTTP status code")
	AddCmd.Flags().StringVarP(&method, "method", "m", "", "HTTP method")
}

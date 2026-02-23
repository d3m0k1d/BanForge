package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	name     string
	service  string
	path     string
	status   string
	method   string
	ttl      string
	maxRetry int
	editName string
)

var RuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage rules",
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new rule to /etc/banforge/rules.d/",
	Long:  "Creates a new rule file in /etc/banforge/rules.d/<name>.toml",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			fmt.Println("Rule name can't be empty (use -n flag)")
			os.Exit(1)
		}
		if service == "" {
			fmt.Println("Service name can't be empty (use -s flag)")
			os.Exit(1)
		}
		if path == "" && status == "" && method == "" {
			fmt.Println("At least one rule field must be filled: path, status, or method")
			os.Exit(1)
		}
		if ttl == "" {
			ttl = "1y"
		}
		err := config.NewRule(name, service, path, status, method, ttl, maxRetry)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Rule added successfully!")
	},
}

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing rule",
	Long:  "Edit rule fields by name. Only specified fields will be updated.",
	Run: func(cmd *cobra.Command, args []string) {
		if editName == "" {
			fmt.Println("Rule name is required (use -n flag)")
			os.Exit(1)
		}
		if service == "" && path == "" && status == "" && method == "" {
			fmt.Println("At least one field must be specified to edit: -s, -p, -c, or -m")
			os.Exit(1)
		}
		err := config.EditRule(editName, service, path, status, method)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Rule updated successfully!")
	},
}

var RemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a rule by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ruleName := args[0]
		fileName := config.SanitizeRuleFilename(ruleName) + ".toml"
		filePath := filepath.Join("/etc/banforge/rules.d", fileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Rule '%s' not found\n", ruleName)
			os.Exit(1)
		}

		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Failed to remove rule: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Rule '%s' removed successfully\n", ruleName)
	},
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rules",
	Run: func(cmd *cobra.Command, args []string) {
		rules, err := config.LoadRuleConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(rules) == 0 {
			fmt.Println("No rules found")
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"Name", "Service", "Path", "Status", "Method", "MaxRetry", "BanTime",
		})

		for _, rule := range rules {
			t.AppendRow(table.Row{
				rule.Name,
				rule.ServiceName,
				rule.Path,
				rule.Status,
				rule.Method,
				rule.MaxRetry,
				rule.BanTime,
			})
		}
		t.Render()
	},
}

func RuleRegister() {
	RuleCmd.AddCommand(AddCmd)
	RuleCmd.AddCommand(EditCmd)
	RuleCmd.AddCommand(RemoveCmd)
	RuleCmd.AddCommand(ListCmd)

	AddCmd.Flags().StringVarP(&name, "name", "n", "", "rule name (required)")
	AddCmd.Flags().StringVarP(&service, "service", "s", "", "service name (required)")
	AddCmd.Flags().StringVarP(&path, "path", "p", "", "request path")
	AddCmd.Flags().StringVarP(&status, "status", "c", "", "status code")
	AddCmd.Flags().StringVarP(&method, "method", "m", "", "HTTP method")
	AddCmd.Flags().StringVarP(&ttl, "ttl", "t", "", "ban time (e.g., 1h, 1d, 1y)")
	AddCmd.Flags().IntVarP(&maxRetry, "max_retry", "r", 0, "max retry before ban")

	EditCmd.Flags().StringVarP(&editName, "name", "n", "", "rule name to edit (required)")
	EditCmd.Flags().StringVarP(&service, "service", "s", "", "new service name")
	EditCmd.Flags().StringVarP(&path, "path", "p", "", "new path")
	EditCmd.Flags().StringVarP(&status, "status", "c", "", "new status code")
	EditCmd.Flags().StringVarP(&method, "method", "m", "", "new HTTP method")
}

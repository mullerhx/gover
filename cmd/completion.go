package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Bash:

  $ source <(gover completion bash)
  $ gover completion bash > /etc/bash_completion.d/gover

Zsh:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  $ gover completion zsh > "${fpath[1]}/_gover"

Fish:

  $ gover completion fish | source
  $ gover completion fish > ~/.config/fish/completions/gover.fish

PowerShell:

  PS> gover completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactValidArgs(1),
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			_ = RootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			_ = RootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = RootCmd.GenPowerShellCompletion(os.Stdout)
		default:
			_, _ = fmt.Fprintf(os.Stderr, "Unsupported shell type: %q\n", args[0])
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}

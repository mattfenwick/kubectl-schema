package apiversions

import "github.com/spf13/cobra"

//type KindArgs struct {}

func SetupKindCommand() *cobra.Command {
	//args := &KindArgs{}

	command := &cobra.Command{
		Use:   "kind",
		Short: "compare types from across swagger specs",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunKind()
		},
	}

	return command
}

func RunKind() {
	ParseKindResults()
}

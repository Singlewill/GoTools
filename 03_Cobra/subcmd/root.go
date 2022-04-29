package subcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var active bool

var RootCmd = &cobra.Command{
	Use:   "mainCMD [subcmd]",
	Short: "mainCMD short Message",
	//命令执行时，按照以下顺序依次执行
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("mainCMD PersistentPreRun!")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("mainCMD PreRun!")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(active)
		fmt.Println(args, errors.New("unrecognized subcommand"))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("mainCMD PostRun!")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("mainCMD PersistentPostRun!")
	},
}

func init() {
	//设置全局现象，适用所有subcmd
	//--active=xxx   --a=xxx
	//xxx = true, True, Flase, flase, 1, 0...
	RootCmd.PersistentFlags().BoolVarP(&active, "active", "a", false, "counter detected malicious activity (dangerous, may clobber)")
	RootCmd.AddCommand(KaloSubCmd)
	//标记必选，针对PersistentFlags标记
	RootCmd.MarkPersistentFlagRequired("active")
}

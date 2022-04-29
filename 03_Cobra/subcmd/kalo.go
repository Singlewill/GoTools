package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var KaloString string

var KaloSubCmd = &cobra.Command{
	//子command名称
	Use: "kaloSubCmd",
	//别名，
	Aliases: []string{"kalo", "KALO"},
	Short:   "kalo's sub cmd short message",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fmt.Println(KaloString)
		fmt.Println(active)
	},
}

func init() {
	// --name==hello
	// -n=hello
	KaloSubCmd.Flags().StringVarP(&KaloString, "name", "n", "hello", "no use")
	//标记必选，针对普通Flags标记
	KaloSubCmd.MarkFlagRequired("name")
	RootCmd.AddCommand(KaloSubCmd)
}

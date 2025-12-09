package godash

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-admin-full/config"
	"go-admin-full/internal/server"
)

var rootCmd = &cobra.Command{
	Use:   "godash",
	Short: "go-admin 服务",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		if err := server.Run(cfg); err != nil {
			fmt.Printf("启动失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "配置文件路径")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	cobra.OnInitialize(func() {
		cfgPath := viper.GetString("config")
		config.Init(cfgPath)
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

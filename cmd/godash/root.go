package godash

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-admin/config"
	"go-admin/internal/server"
)

// 声明当前变量。
var rootCmd = &cobra.Command{
	// 处理当前语句逻辑。
	Use: "godash",
	// 处理当前语句逻辑。
	Short: "go-admin 服务",
	// 调用func完成当前处理。
	Run: func(cmd *cobra.Command, args []string) {
		// 定义并初始化当前变量。
		cfg := config.GetConfig()
		// 判断条件并进入对应分支逻辑。
		if err := server.Run(cfg); err != nil {
			// 调用fmt.Printf完成当前处理。
			fmt.Printf("启动失败: %v\n", err)
			// 调用os.Exit完成当前处理。
			os.Exit(1)
		}
	},
}

// Execute 执行命令入口。
func Execute() {
	// 调用rootCmd.PersistentFlags完成当前处理。
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "配置文件路径")
	// 调用viper.BindPFlag完成当前处理。
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		return
	}

	// 调用cobra.OnInitialize完成当前处理。
	cobra.OnInitialize(func() {
		// 定义并初始化当前变量。
		cfgPath := viper.GetString("config")
		// 调用config.Init完成当前处理。
		config.Init(cfgPath)
	})

	// 判断条件并进入对应分支逻辑。
	if err := rootCmd.Execute(); err != nil {
		// 调用fmt.Println完成当前处理。
		fmt.Println(err)
		// 调用os.Exit完成当前处理。
		os.Exit(1)
	}
}

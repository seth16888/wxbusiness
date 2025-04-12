package cmd

import (
	"fmt"
	"os"

	"github.com/seth16888/wxbusiness/internal/bootstrap"
	"github.com/seth16888/wxbusiness/internal/di"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "c",
		"conf/conf.yaml", "--conf config file (default is conf/conf.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "wxbusiness [command] [flags] [args]",
	Short: "A WX application business server",
	Long:  `A WX application business server`,
  PreRunE: func(cmd *cobra.Command, args []string) error {
    defer func() {
      if err:= recover(); err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(1)
      }
    }()

    conf,err:=bootstrap.InitConfig(configFile)
    if err != nil {
      return err
    }
    log:=bootstrap.InitLogger(conf.Log)

    di.Get().Conf=conf
    di.Get().Log=log

    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    defer func() {
			if di.Get() != nil && di.Get().Log != nil {
				di.Get().Log.Sync() // flushes buffer, if any
			}
		}()

    return bootstrap.StartApp()
  },
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/bootstrap"
	"github.com/seth16888/wxbusiness/internal/data"
	"github.com/seth16888/wxbusiness/internal/di"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"github.com/seth16888/wxcommon/hc"
	"github.com/spf13/cobra"
)

var configFile string

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
			if err := recover(); err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
		}()

		conf, err := bootstrap.InitConfig(configFile)
		if err != nil {
			return err
		}
		log := bootstrap.InitLogger(conf.Log)

		di.Get().Conf = conf
		di.Get().Log = log

		di.Get().HttpClient = hc.NewClient(10*time.Second,
			3*time.Minute, hc.CommonCheckRedirect)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func() {
			if di.Get() != nil && di.Get().Log != nil {
				di.Get().Log.Sync() // flushes buffer, if any
			}
		}()

		err := bootstrap.InitDatabase(di.Get().Conf.DB, di.Get().Log)
		if err != nil {
			return err
		}

		tokenServerAddr := di.Get().Conf.TokenServer.Addr
		if len(tokenServerAddr) == 0 {
			return fmt.Errorf("token server addr is empty")
		}
		tokenClient, err := bootstrap.InitTokenCLient(tokenServerAddr)
		if err != nil {
			return err
		}
		di.Get().TokenClient = tokenClient
		di.Get().Validator = validator.NewValidator()

		userAppRepo := data.NewUserAppData(di.Get().DB, di.Get().Log)
		platformAppRepo := data.NewPlatformAppRepo(di.Get().DB, di.Get().Log) // 先初始化

		// gRPC timeout
		gRPCTimeout := 15 * time.Second

		tokenProxy := biz.NewAccessTokenUsecase(tokenClient, di.Get().Log, gRPCTimeout, platformAppRepo)
		portalUsecase := biz.NewPortalUsecase(platformAppRepo, tokenProxy, di.Get().Log)
		di.Get().PortalUsecase = portalUsecase

		di.Get().AppUsecase = biz.NewAppUsecase(platformAppRepo, di.Get().Log)

		apiProxyClient, err := bootstrap.InitAPIProxyClient(di.Get().Conf.ProxyServer.Addr)
		if err != nil {
			return err
		}
		apiProxy := biz.NewAPIProxyUsecase(apiProxyClient, di.Get().Log, gRPCTimeout, tokenProxy)
    menuRepo := data.NewMPMenuData(di.Get().Log, di.Get().DB)
		menuUsecase := biz.NewMPMenuUsecase(platformAppRepo, tokenProxy, apiProxy, di.Get().Log, menuRepo)
		di.Get().MenuUsecase = menuUsecase

		coAuthClient, err := bootstrap.InitAuthClient(di.Get().Conf.CoAuthServer.Addr)
		if err != nil {
			return err
		}
		di.Get().CoAuthClient = coAuthClient

		userUc := biz.NewUserUsecase(userAppRepo)
		di.Get().UserUsecase = userUc

		tagRepo := data.NewMemberTagData(di.Get().DB, di.Get().Log)
		tagUc := biz.NewMemberTagUsecase(tagRepo, di.Get().Log, apiProxy)
		di.Get().MemberTagUsecase = tagUc

		memberRepo := data.NewMPMemberData(di.Get().DB, di.Get().Log)
		blockRepo := data.NewMPBlackListData(di.Get().DB, di.Get().Log)
		memberUc := biz.NewMPMemberUsecase(di.Get().Log, memberRepo, apiProxy, blockRepo)
		di.Get().MPMemberUsecase = memberUc

		materialRepo := data.NewMPMaterialData(di.Get().DB, di.Get().Log)
		materialUc := biz.NewMaterialUsecase(di.Get().Log, materialRepo, apiProxy,
			di.Get().HttpClient)
		di.Get().MaterialUsecase = materialUc

		qrcodeUc := biz.NewMpQRCodeUsecase(di.Get().Log, platformAppRepo, tokenProxy, apiProxy)
		di.Get().MpQRCodeUsecase = qrcodeUc

		return bootstrap.StartApp(di.Get())
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

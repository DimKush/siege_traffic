package main

import (
	"fmt"
	"github.com/DimKush/siege_traffic/internal/logger"
	"github.com/DimKush/siege_traffic/internal/option"
	sniffer "github.com/DimKush/siege_traffic/internal/sniffer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

func NewApp() *cobra.Command {

	opt := &option.Options{}
	app := &cobra.Command{
		Use:   "siege_traffic",
		Short: "Catch and analyze games traffic UDP (in-game packages)",
		Run: func(cmd *cobra.Command, args []string) {
			if opt.FilterServerIP == "" {
				fmt.Println("cannot run. Server ip is not setted")
				os.Exit(1)
			}

			snf := sniffer.NewSniffer(opt)
			snf.Start()
		},
	}

	app.Flags().StringVarP(&opt.DeviceName, "device-name", "d", "\\Device\\NPF_{BC0CDC3B-9E5A-4A2C-9A69-6CC121A05865}", "device name")
	app.Flags().StringVarP(&opt.FilterClientIP, "client-ip", "c", "192.168.3.8", "client ip")
	app.Flags().StringVarP(&opt.FilterServerIP, "server-ip", "s", "", "server ip")

	app.Flags().PrintDefaults()
	return app
}

func init() {
	_, err := logger.NewLogger()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//log.Logger = zLog
		log.Info().Msg("Logger initialized.")
	}
}
func main() {
	app := NewApp()
	if err := app.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

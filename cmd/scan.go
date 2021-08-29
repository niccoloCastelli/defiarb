package cmd

import (
	"context"
	"fmt"
	"github.com/niccoloCastelli/defiarb/exchanges"
	_ "github.com/niccoloCastelli/defiarb/exchanges/bakery"
	_ "github.com/niccoloCastelli/defiarb/exchanges/pancake"
	"github.com/niccoloCastelli/defiarb/services/bsc_client"
	"github.com/niccoloCastelli/defiarb/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan exchange1 exchange2",
	Long:  `Find LP for exchange(s)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if conf == nil {
			return errors.New("config not found")
		}
		fmt.Println("scan called", walletAddr)
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()
		client, err := bsc_client.NewClient(conf.NodeUrl, ctx)
		if err != nil {
			return errors.WithStack(err)
		}

		Db, err := storage.NewDb(&conf.Db)
		if err != nil {
			return errors.WithStack(err)
		}
		Db = Db.LogMode(true)
		if err := storage.Migrate(Db); err != nil {
			return errors.WithStack(err)
		}
		for _, exchangeName := range args {
			fmt.Println("Scan exchange: ", exchangeName)
			exchange, err := exchanges.New(exchangeName, client, "")
			if err != nil {
				return err
			}
			_, addrChan, errChan := exchange.GetAllPairsAsync(ctx)
			if len(errChan) > 0 {
				return <-errChan
			}
			for {
				select {
				case err := <-errChan:
					fmt.Println("Get pairs error: ", err.Error())
					cancelFn()
					goto endloop
				case pairAddr := <-addrChan:
					if pairAddr == nil {
						goto endloop
					}
					lp, err := exchange.GetPairInfo(*pairAddr)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if err := Db.Save(lp).Error; err != nil {
						return err
					}
				}
			}
		endloop:
		}
		return nil
	},
}

/**/

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

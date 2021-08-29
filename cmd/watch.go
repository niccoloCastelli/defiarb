package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/contracts/bep20"
	"github.com/niccoloCastelli/defiarb/services/wallet"
	"github.com/niccoloCastelli/defiarb/storage"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/niccoloCastelli/defiarb/tokens"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"math"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Check wallet status",
	Long:  `Check wallet status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("wallet called", walletAddr)
		ctx := context.Background()
		Db, err := storage.NewDb(&conf.Db)
		if err != nil {
			return errors.WithStack(err)
		}
		Db = Db.LogMode(true)
		supportedLpTokens := models.LiquidityPools{}
		if err := supportedLpTokens.LoadSupportedPools(Db); err != nil {
			return err
		}
		client, err := ethclient.Dial(conf.NodeUrl)
		if err != nil {
			return err
		}
		chainId, err := client.ChainID(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Chain ID: ", chainId.String())
		address := common.HexToAddress(walletAddr)
		Wallet := wallet.NewWallet(address, client)
		balance, err := Wallet.GetBalance(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Balance: ", float64(balance)/math.Pow(10, 18))
		supportedTokens := tokens.SupportedTokens()
		for _, token := range supportedTokens {
			tokenBalance, err := token.FormattedBalance(client, address)
			if err != nil {
				return err
			}
			fmt.Printf("Balance: %s %f [%s]\n", token.Code(), tokenBalance, token.Address())
		}

		for _, lpToken := range supportedLpTokens {
			tokenBalance, err := lpToken.LpToken().FormattedBalance(client, address)
			if err != nil {
				return err
			}
			fmt.Printf("Balance: %s %f [%s]\n", lpToken.Symbol, tokenBalance, lpToken.Address)
		}

		txs, err := Wallet.GetTransactions(ctx)
		if err != nil {
			return err
		}
		for _, tx := range txs {
			instance, err := bep20.NewBep20(tx.Address, client)
			if err != nil {
				return err
			}
			t, err := instance.ParseTransfer(tx)
			if err != nil {
				return err
			}
			fmt.Println(tx.Address.String(), t.From.String(), t.To.String(), float64(t.Value.Uint64())/math.Pow(10, 18))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
}

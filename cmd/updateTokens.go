package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/storage"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"regexp"
)

type tokenInfo struct {
	Name_       string `json:"name"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Explorer    string `json:"explorer"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	Status      string `json:"status"`
	Id          string `json:"id"`
}

func (t *tokenInfo) LogoUrl() string {
	return ""
}

func (t *tokenInfo) Code() string {
	return t.Symbol
}

func (t *tokenInfo) Name() string {
	return t.Name_
}

func (t *tokenInfo) Address() string {
	return t.Id
}

func (t *tokenInfo) BalanceOf(client *ethclient.Client, account common.Address) (uint64, error) {
	return 0, errors.New("not supported")
}

func (t *tokenInfo) FormattedBalance(client *ethclient.Client, account common.Address) (float64, error) {
	return 0, errors.New("not supported")
}

// updateTokensCmd represents the updateTokens command
var updateTokensCmd = &cobra.Command{
	Use:   "updateTokens",
	Short: "Update tokens list",
	Long:  `Download token list from trustwallet github`,
	RunE: func(cmd *cobra.Command, args []string) error {

		Db, err := storage.NewDb(&conf.Db)
		if err != nil {
			return errors.WithStack(err)
		}

		matchAssetRe, err := regexp.Compile(`^blockchains/smartchain/assets/([0-9A-Za-z]+)/info\.json$`)
		if err != nil {
			return err
		}

		fmt.Println("Cloning git repository...")
		r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: "https://github.com/trustwallet/assets.git",
		})
		if err != nil {
			return err
		}
		head, err := r.Head()
		if err != nil {
			return err
		}
		commit, err := r.CommitObject(head.Hash())
		if err != nil {
			return err
		}
		fileIter, err := commit.Files()
		if err != nil {
			return err
		}
		fmt.Println("Searching files...")
		err = fileIter.ForEach(func(f *object.File) error {

			m := matchAssetRe.FindStringSubmatch(f.Name)
			if len(m) == 0 {
				return nil
			}
			content, err := f.Contents()
			if err != nil {
				return err
			}

			tokenInfo := tokenInfo{}
			if err := json.Unmarshal([]byte(content), &tokenInfo); err != nil {
				return err
			}
			if tokenInfo.Type == "BEP20" && tokenInfo.Status == "active" {
				dbModel := models.NewTokenFromERC20(&tokenInfo)
				fmt.Println(tokenInfo.Name, tokenInfo.Symbol)
				err := Db.Save(&dbModel).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
		return err
	},
}

func init() {
	rootCmd.AddCommand(updateTokensCmd)
}

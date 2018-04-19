package app

import (
	"testing"
	"path/filepath"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/cosmos/cosmos-sdk/server"
	"fmt"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"
	"github.com/spf13/viper"
)

func TestNewIrisApp(t *testing.T) {
	ctx :=server.NewDefaultContext()
	appCreator := func(rootDir string, logger log.Logger) (abci.Application, error){
		rootDir = "/Users/zhangzhiqiang/.iris-hub1"
		dbMain, err := dbm.NewGoLevelDB("iris", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbAcc, err := dbm.NewGoLevelDB("iris-acc", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbPow, err := dbm.NewGoLevelDB("iris-pow", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbIBC, err := dbm.NewGoLevelDB("iris-ibc", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbStaking, err := dbm.NewGoLevelDB("iris-staking", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbIris, err := dbm.NewGoLevelDB("iris-service", filepath.Join(rootDir, "data"))
		if err != nil {
			fmt.Print(err)
		}
		dbs := map[string]dbm.DB{
			"main":    dbMain,
			"acc":     dbAcc,
			"pow":     dbPow,
			"ibc":     dbIBC,
			"staking": dbStaking,
			"iservice": dbIris,
		}
		bapp := NewIrisApp(ctx.Logger, dbs)
		return bapp,nil
	}

	cmd := server.StartCmd(appCreator,ctx)

	viper.Set("with-tendermint", true)
	//viper.Set("address", "tcp://0.0.0.0:46658")

	cmd.Execute()


}
package mock

import (
	"math/rand"
	"os"

	bam "github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/guardian"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	chainID = ""
)

const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	bech32PrefixAccAddr = "faa"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	bech32PrefixAccPub = "fap"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	bech32PrefixValAddr = "fva"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	bech32PrefixValPub = "fvp"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	bech32PrefixConsAddr = "fca"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	bech32PrefixConsPub = "fcp"
)

var (
	IrisCt = types.NewDefaultCoinType(stakeTypes.StakeDenomName)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*bam.BaseApp
	Cdc              *codec.Codec
	TotalCoinsSupply sdk.Coins
}

// NewApp partially constructs a new app on the memstore for module and genesis
// testing.
func NewApp() *App {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)

	cdc := MakeCodec()

	bApp := bam.NewBaseApp("mock", logger, db, auth.DefaultTxDecoder(cdc), bam.SetPruning("nothing"))

	// create your application object
	var app = &App{
		BaseApp: 		  bApp,
		Cdc:              cdc,
		TotalCoinsSupply: sdk.Coins{},
	}
	engine := protocol.NewProtocolEngine()

	protocol0 := NewProtocolVersion0(cdc)
	engine.Add(protocol0)

	engine.LoadCurrentProtocol()
	app.SetProtocolEngine(engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())
	err := app.LoadLatestVersion(engine.GetKeyMain())
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	record.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

//// CompleteSetup completes the application setup after the routes have been
//// registered.
//func (app *App) CompleteSetup(newKeys ...sdk.StoreKey) error {
//	newKeys = append(newKeys, protocol.KeyMain)
//	newKeys = append(newKeys, protocol.KeyAccount)
//	newKeys = append(newKeys, protocol.KeyParams)
//	newKeys = append(newKeys, protocol.KeyStake)
//	newKeys = append(newKeys, protocol.KeyFeeCollection)
//	newKeys = append(newKeys, protocol.TkeyParams)
//	newKeys = append(newKeys, protocol.TkeyStake)
//
//	for _, key := range newKeys {
//		switch key.(type) {
//		case *sdk.KVStoreKey:
//			app.MountStore(key, sdk.StoreTypeIAVL)
//		case *sdk.TransientStoreKey:
//			app.MountStore(key, sdk.StoreTypeTransient)
//		default:
//			return fmt.Errorf("unsupported StoreKey: %+v", key)
//		}
//	}
//
//	err := app.LoadLatestVersion(protocol.KeyMain)
//	return err
//}


// CreateGenAccounts generates genesis accounts loaded with coins, and returns
// their addresses, pubkeys, and privkeys.
func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (genAccs []auth.Account, addrs []sdk.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {
	for i := 0; i < numAccs; i++ {
		privKey := ed25519.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		genAcc := &auth.BaseAccount{
			Address: addr,
			Coins:   genCoins,
		}

		genAccs = append(genAccs, genAcc)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrs = append(addrs, addr)
	}

	return
}

// SetGenesis sets the mock app genesis accounts.
func SetGenesis(app *App, accs []auth.Account) {
	// Pass the accounts in via the application (lazy) instead of through
	// RequestInitChain.
	app.Engine.GetCurrent().(*ProtocolVersion0).GenesisAccounts = accs

	app.InitChain(abci.RequestInitChain{})
	app.Commit()
}

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	// Make the transaction free
	fee := auth.StdFee{
		Amount: sdk.Coins{sdk.NewInt64Coin(stakeTypes.StakeDenom, 0)},
		Gas:    20000,
	}

	sigs := make([]auth.StdSignature, len(priv))
	memo := "testmemotestmemo"

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:        p.PubKey(),
			Signature:     sig,
			AccountNumber: accnums[i],
			Sequence:      seq[i],
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, memo)
}

// GeneratePrivKeys generates a total n Ed25519 private keys.
func GeneratePrivKeys(n int) (keys []crypto.PrivKey) {
	// TODO: Randomize this between ed25519 and secp256k1
	keys = make([]crypto.PrivKey, n, n)
	for i := 0; i < n; i++ {
		keys[i] = ed25519.GenPrivKey()
	}

	return
}

// GeneratePrivKeyAddressPairs generates a total of n private key, address
// pairs.
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	keys = make([]crypto.PrivKey, n, n)
	addrs = make([]sdk.AccAddress, n, n)
	for i := 0; i < n; i++ {
		if rand.Int63()%2 == 0 {
			keys[i] = secp256k1.GenPrivKey()
		} else {
			keys[i] = ed25519.GenPrivKey()
		}
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

// GeneratePrivKeyAddressPairsFromRand generates a total of n private key, address
// pairs using the provided randomness source.
func GeneratePrivKeyAddressPairsFromRand(rand *rand.Rand, n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	keys = make([]crypto.PrivKey, n, n)
	addrs = make([]sdk.AccAddress, n, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		if err != nil {
			panic("Could not read randomness")
		}
		if rand.Int63()%2 == 0 {
			keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		} else {
			keys[i] = ed25519.GenPrivKeyFromSecret(secret)
		}
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

// RandomSetGenesis set genesis accounts with random coin values using the
// provided addresses and coin denominations.
func RandomSetGenesis(r *rand.Rand, app *App, addrs []sdk.AccAddress, denoms []string) {
	accts := make([]auth.Account, len(addrs), len(addrs))
	//randCoinIntervals := []BigInterval{
	//	{sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 1)},
	//	{sdk.NewIntWithDecimal(1, 2), sdk.NewIntWithDecimal(1, 3)},
	//	{sdk.NewIntWithDecimal(1, 40), sdk.NewIntWithDecimal(1, 50)},
	//}

	for i := 0; i < len(accts); i++ {
		coins := make([]sdk.Coin, len(denoms), len(denoms))

		amount := sdk.NewIntWithDecimal(1, 2)
		// generate a random coin for each denomination
		for j := 0; j < len(denoms); j++ {
			coins[j] = sdk.Coin{Denom: denoms[j],
				//Amount: RandFromBigInterval(r, randCoinIntervals).Add(amount),
				Amount: amount,
			}
		}

		app.TotalCoinsSupply = app.TotalCoinsSupply.Plus(coins)
		baseAcc := auth.NewBaseAccountWithAddress(addrs[i])

		(&baseAcc).SetCoins(coins)
		accts[i] = &baseAcc
	}

	app.Engine.GetCurrent().(*ProtocolVersion0).GenesisAccounts = accts
}

// GetAllAccounts returns all accounts in the AccountKeeper.
func GetAllAccounts(mapper auth.AccountKeeper, ctx sdk.Context) []auth.Account {
	accounts := []auth.Account{}
	appendAccount := func(acc auth.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	mapper.IterateAccounts(ctx, appendAccount)
	return accounts
}

// GenSequenceOfTxs generates a set of signed transactions of messages, such
// that they differ only by having the sequence numbers incremented between
// every transaction.
func GenSequenceOfTxs(msgs []sdk.Msg, accnums []uint64, initSeqNums []uint64, numToGenerate int, priv ...crypto.PrivKey) []auth.StdTx {
	txs := make([]auth.StdTx, numToGenerate, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		txs[i] = GenTx(msgs, accnums, initSeqNums, priv...)
		incrementAllSequenceNumbers(initSeqNums)
	}

	return txs
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}

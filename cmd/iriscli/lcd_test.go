package main

import ("testing"
	"github.com/irisnet/irishub/app"
	"github.com/spf13/viper"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"bytes"
	"log"
	"os"
	"encoding/hex"
	"github.com/irisnet/irishub/tools/prometheus"
)


func TestRestServer(t *testing.T) {
	cdc := app.MakeCodec()
	comm := ServeCommand(cdc)
	viper.Set("chain-id","fuxi")
	viper.Set("node","tcp://localhost:26657")
	viper.Set("laddr","tcp://localhost:1317")

	comm.ExecuteC()
}

func TestMetricsCmd(t *testing.T){
	cdc := app.MakeCodec()
	comm := prometheus.MonitorCommand("stake",cdc)
	viper.Set("node","tcp://0.0.0.0:46657")
	viper.Set("chain-id","fuxi")
	viper.Set("commands","iris start;htop")
	viper.Set("paths","/etc ")
	comm.ExecuteC()
}

func TestAmino(t *testing.T) {
	var cdc = amino.NewCodec()

	type Person struct {
		Name string
		Age  int
	}

	p := Person{
		Name :"zhangsn",
		Age:18,
	}

	bz,_ := cdc.MarshalJSON(p)
	fmt.Println(bz)

	var p1 Person
	if err := cdc.UnmarshalJSON(bz,&p1);err != nil {
		fmt.Println(err)
		return
	}
	//cdc.MustUnmarshalBinaryBare(bz,p1)
	fmt.Printf("%+v",p1)

}

func TestMarshalJson(t *testing.T) {
	add,_ := types.GetAccAddressHex("1EC2E86065D5EF88A3ED65B8B3A43210FAD9C7B2")
	conis := types.Coins{
		types.Coin{
			Denom:"iris",
			Amount:10,
		},
	}
	//fmt.Println(types.MustBech32ifyAcc(add))
	input := bank.Input{
		Address:add,
		Coins:conis,
	}
	input1Byte := input.GetSignBytes()
	fmt.Println(input1Byte)

	input2Byte,_ := json.Marshal(struct {
		Address string    `json:"address"`
		Coins   types.Coins `json:"coins"`
	}{
		Address: types.MustBech32ifyAcc(add),
		Coins:   conis,
	})

	assert.Equal(t,input1Byte,input2Byte)

}

func TestMarshal(t *testing.T) {
	add,_ := types.GetAccAddressHex("1EC2E86065D5EF88A3ED65B8B3A43210FAD9C7B2")
	conis := types.Coins{
		types.Coin{
			Denom:"iris",
			Amount:10,
		},
	}
	//fmt.Println(types.MustBech32ifyAcc(add))
	input := bank.Input{
		Address:add,
		Coins:conis,
	}
	input1Byte := input.GetSignBytes()

	input2Byte,_ := json.Marshal(struct {
		Address string    `json:"address"`
		Coins   types.Coins `json:"coins"`
	}{
		Address: types.MustBech32ifyAcc(add),
		Coins:   conis,
	})

	assert.Equal(t,input1Byte,input2Byte)

}

func TestSignByte(t *testing.T) {
	from,_ := types.GetAccAddressHex("1EC2E86065D5EF88A3ED65B8B3A43210FAD9C7B2")

	chainID := "fuxi"
	accnums := []int64{0}
	sequences := []int64{3}

	conis := types.Coins{
		types.Coin{
			Denom:"iris",
			Amount:10,
		},
	}

	fees := types.Coins{
		types.Coin{
			Denom:"iris",
			Amount:0,
		},
	}

	type put struct {
		Address string    `json:"address"`
		Coins   types.Coins `json:"coins"`
	}

	input := bank.Input{
		Address:from,
		Coins:conis,
	}

	input1 := input.GetSignBytes()
	input2,_ := json.Marshal(put{
		Address: types.MustBech32ifyAcc(from),
		Coins:   conis,
	})

	assert.Equal(t,input1,input2)

	to,_ := types.GetAccAddressHex("3A058A8B5468AE0EA2D2517CE3BAFDD281E50C2F")
	output := bank.Output{
		Address:to,
		Coins:conis,
	}

	output1 := output.GetSignBytes()
	output2,_ := json.Marshal(put{
		Address: types.MustBech32ifyAcc(to),
		Coins:   conis,
	})

	assert.Equal(t,output1,output2)

	msg := bank.MsgSend{
		[]bank.Input{input},[]bank.Output{output},
	}

	msg1 := msg.GetSignBytes()

	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}

	msg2,_ := json.Marshal(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	//fmt.Print(msg2)
	assert.Equal(t,msg1,msg2)

	fee := auth.StdFee{
		Amount:fees,
		Gas:int64(0),
	}

	fee1 := fee.Bytes()
	fee2,_:= json.Marshal(struct {
		Amount types.Coins `json:"amount"`
		Gas    int64       `json:"gas"`
	}{
		Amount:fees,
		Gas:int64(0),
	})

	assert.Equal(t,fee1,fee2)

	signMsg := auth.StdSignMsg{
		ChainID:chainID,
		AccountNumbers:accnums,
		Sequences:sequences,
		Msg:msg,
		Fee:fee,
	}

	signByte1 := signMsg.Bytes()
	signByte2,_ := json.Marshal(auth.StdSignDoc{
		ChainID:        chainID,
		AccountNumbers: accnums,
		Sequences:      sequences,
		FeeBytes:       fee2,
		MsgBytes:       msg2,
	})

	str := string(signByte1)

	fmt.Println(str)

	//printJson(signByte2)

	assert.Equal(t,signByte1,signByte2)

	fmt.Println(signByte1)

}

func printJson(b []byte)  {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	out.WriteTo(os.Stdout)
}

func TestH(t *testing.T){
	client := []byte{123,34,99,104,97,105,110,95,105,100,34,58,34,102,117,120,105,34,44,34,97,99,99,111,117,110,116,95,110,117,109,98,101,114,115,34,58,91,48,93,44,34,115,101,113,117,101,110,99,101,115,34,58,91,51,93,44,34,102,101,101,95,98,121,116,101,115,34,58,34,101,121,74,104,98,87,57,49,98,110,81,105,79,108,116,55,73,109,70,116,98,51,86,117,100,67,73,54,77,67,119,105,90,71,86,117,98,50,48,105,79,105,74,112,99,109,108,122,73,110,49,100,76,67,74,110,89,88,77,105,79,106,66,57,34,44,34,109,115,103,95,98,121,116,101,115,34,58,34,101,121,74,112,98,110,66,49,100,72,77,105,79,108,116,55,73,109,70,107,90,72,74,108,99,51,77,105,79,105,74,106,98,51,78,116,98,51,78,104,89,50,78,104,90,71,82,121,77,88,74,116,99,72,100,122,89,51,73,53,78,109,104,111,89,122,78,110,98,71,82,50,97,51,86,48,79,71,90,119,97,110,112,121,89,87,82,117,77,50,70,113,100,110,74,116,98,71,100,107,73,105,119,105,89,50,57,112,98,110,77,105,79,108,116,55,73,109,70,116,98,51,86,117,100,67,73,54,77,84,65,115,73,109,82,108,98,109,57,116,73,106,111,105,97,88,74,112,99,121,74,57,88,88,49,100,76,67,74,118,100,88,82,119,100,88,82,122,73,106,112,98,101,121,74,104,90,71,82,121,90,88,78,122,73,106,111,105,89,50,57,122,98,87,57,122,89,87,78,106,89,87,82,107,99,106,69,52,90,51,112,106,78,72,111,50,78,87,82,54,97,72,70,104,90,50,116,113,77,106,107,51,100,122,104,51,97,71,69,50,77,110,69,51,77,110,74,119,77,72,100,108,101,87,53,48,79,67,73,115,73,109,78,118,97,87,53,122,73,106,112,98,101,121,74,104,98,87,57,49,98,110,81,105,79,106,69,119,76,67,74,107,90,87,53,118,98,83,73,54,73,109,108,121,97,88,77,105,102,86,49,57,88,88,48,61,34,44,34,97,108,116,95,98,121,116,101,115,34,58,110,117,108,108,125}
	clientS := string(client)

	//fmt.Println(str1)
	fmt.Println(clientS)

	server := []byte{123,34,99,104,97,105,110,95,105,100,34,58,34,102,117,120,105,34,44,34,97,99,99,111,117,110,116,95,110,117,109,98,101,114,115,34,58,91,48,93,44,34,115,101,113,117,101,110,99,101,115,34,58,91,51,93,44,34,102,101,101,95,98,121,116,101,115,34,58,34,101,121,74,104,98,87,57,49,98,110,81,105,79,108,116,55,73,109,82,108,98,109,57,116,73,106,111,105,97,88,74,112,99,121,73,115,73,109,70,116,98,51,86,117,100,67,73,54,77,72,49,100,76,67,74,110,89,88,77,105,79,106,66,57,34,44,34,109,115,103,95,98,121,116,101,115,34,58,34,101,121,74,112,98,110,66,49,100,72,77,105,79,108,116,55,73,109,70,107,90,72,74,108,99,51,77,105,79,105,74,106,98,51,78,116,98,51,78,104,89,50,78,104,90,71,82,121,77,88,74,116,99,72,100,122,89,51,73,53,78,109,104,111,89,122,78,110,98,71,82,50,97,51,86,48,79,71,90,119,97,110,112,121,89,87,82,117,77,50,70,113,100,110,74,116,98,71,100,107,73,105,119,105,89,50,57,112,98,110,77,105,79,108,116,55,73,109,82,108,98,109,57,116,73,106,111,105,97,88,74,112,99,121,73,115,73,109,70,116,98,51,86,117,100,67,73,54,77,84,66,57,88,88,49,100,76,67,74,118,100,88,82,119,100,88,82,122,73,106,112,98,101,121,74,104,90,71,82,121,90,88,78,122,73,106,111,105,89,50,57,122,98,87,57,122,89,87,78,106,89,87,82,107,99,106,69,52,90,51,112,106,78,72,111,50,78,87,82,54,97,72,70,104,90,50,116,113,77,106,107,51,100,122,104,51,97,71,69,50,77,110,69,51,77,110,74,119,77,72,100,108,101,87,53,48,79,67,73,115,73,109,78,118,97,87,53,122,73,106,112,98,101,121,74,107,90,87,53,118,98,83,73,54,73,109,108,121,97,88,77,105,76,67,74,104,98,87,57,49,98,110,81,105,79,106,69,119,102,86,49,57,88,88,48,61,34,44,34,97,108,116,95,98,121,116,101,115,34,58,110,117,108,108,125}
	serverS := string(server)
	fmt.Println(serverS)
}

func TestBech32(t *testing.T){
	addr,_ := types.GetAccAddressBech32("cosmosaccaddr1wp4shn6zzv0l52hfat9c2ryu0gvwa4h3dj2k92")
	fmt.Print(hex.EncodeToString(addr))
}





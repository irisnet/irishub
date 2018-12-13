package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v0"
)

func TestIrisCLISubmitRecord(t *testing.T) {
	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	// submit q first record onchain test
	srStr := fmt.Sprintf("iriscli record submit %v", flags)
	srStr += fmt.Sprintf(" --description=%s", "test")
	srStr += fmt.Sprintf(" --onchain-data=%s", "record-test")
	srStr += fmt.Sprintf(" --from=%s", "foo")
	srStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	srStr += " --json"

	recordTxHash := executeSubmitRecordAndGetTxHash(t, srStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	recordID1 := executeGetRecordID(t, fmt.Sprintf("iriscli tendermint tx %v --output json --trust-node=true %v", recordTxHash, flags))

	// Submit same record twice
	res, _ := tests.ExecuteT(t, srStr, "")
	require.Equal(t, fmt.Sprintf("Warning: Record ID %v already exists.", string(recordID1)), res)

	record1 := executeGetRecord(t, fmt.Sprintf("iriscli record query --record-id=%s --output=json %v", recordID1, flags))
	require.Equal(t, recordID1, record1.RecordID)
	require.Equal(t, fooAddr, record1.OwnerAddress)
	require.Equal(t, "record-test", record1.Data)

	downloadOK := executeDownloadRecord(t, fmt.Sprintf("iriscli record download --record-id=%s --file-name=%s --path=%s %v", recordID1, "download.txt", iriscliHome, flags), iriscliHome+"/download.txt", true)
	require.Equal(t, true, downloadOK)

	res, _ = tests.ExecuteT(t, fmt.Sprintf("iriscli record download --record-id=%s --file-name=%s --path=%s %v", recordID1, "download.txt", iriscliHome, flags), "")
	//require.Equal(t, fmt.Sprintf("Warning: %s already exists, please try another file name.", iriscliHome+"/download.txt"), res)

	// submit a second record onchain test
	srStr = fmt.Sprintf("iriscli record submit %v", flags)
	srStr += fmt.Sprintf(" --description=%s", "test")
	srStr += fmt.Sprintf(" --onchain-data=%s", "record-test2")
	srStr += fmt.Sprintf(" --from=%s", "foo")
	srStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	srStr += " --json"

	recordTxHash = executeSubmitRecordAndGetTxHash(t, srStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	recordID2 := executeGetRecordID(t, fmt.Sprintf("iriscli tendermint tx %v --output json --trust-node=true %v", recordTxHash, flags))

	record2 := executeGetRecord(t, fmt.Sprintf("iriscli record query --record-id=%s --output=json %v", recordID2, flags))
	require.Equal(t, recordID2, record2.RecordID)
	require.Equal(t, fooAddr, record2.OwnerAddress)
	require.Equal(t, "record-test2", record2.Data)
}

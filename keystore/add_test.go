package keystore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/irisnet/irishub/address"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testHome       = "./keystest"
	defaultKeyPass = "1234567890"
)

func Test_runAddCmdBasic(t *testing.T) {
	address.ConfigureBech32Prefix()
	for i := 0; i < 10000; i++ {
		_, out, err1 := executeWriteRetStdStreams(t, fmt.Sprintf("iriscli keys add --home=%s foo", testHome), defaultKeyPass)
		fmt.Println(string(out))
		_ = err1

		ss := strings.Split(out, "\n")
		s := ss[len(ss)-2]
		addr := strings.Split(ss[1], "	")[2]
		fmt.Println(addr)

		cmd := keys.AddKeyCommand()
		cmd.Flags().AddFlagSet(keys.Commands("home").PersistentFlags())

		mockIn := testutil.ApplyMockIODiscardOutErr(cmd)
		kbHome := t.TempDir()

		kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
		require.NoError(t, err)

		clientCtx := client.Context{}.WithKeyringDir(kbHome)
		ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

		t.Cleanup(func() {
			_ = kb.Delete("keyname1")
		})

		// In recovery mode
		cmd.SetArgs([]string{
			"keyname1",
			fmt.Sprintf("--%s=true", "recover"),
			fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
			fmt.Sprintf("--%s=%s", cli.OutputFlag, keys.OutputFormatText),
			fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
			fmt.Sprintf("--%s=%s", flags.FlagKeyAlgorithm, string(hd.Secp256k1Type)),
		})

		// use valid mnemonic and complete recovery key generation successfully
		mockIn.Reset(s + "\n")
		require.NoError(t, cmd.ExecuteContext(ctx))

		info, err := kb.Key("keyname1")
		require.NoError(t, err)
		require.Equal(t, addr, info.GetAddress().String())
		executeWriteCheckErr(t, fmt.Sprintf("iriscli keys delete --home=%s foo", testHome), defaultKeyPass)
	}
}

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	// broadcast transaction and return after the transaction is included by a block
	if strings.Contains(cmdStr, "--from") && strings.Contains(cmdStr, "--fee") {
		cmdStr = cmdStr + " --commit"
	}

	exitSuccess, _, _ = executeWriteRetStdStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStdStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
	proc := GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}
	// Log output.
	if len(stdout) > 0 {
		t.Log("Stdout:", string(stdout))
	}
	if len(stderr) > 0 {
		t.Log("Stderr:", string(stderr))
	}

	proc.Wait()
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

func GoExecuteT(t *testing.T, cmd string) (proc *Process) {
	t.Log("Running", cmd)

	// Split cmd to name and args.
	split := strings.Split(cmd, " ")
	require.True(t, len(split) > 0, "no command provided")
	name, args := split[0], []string(nil)
	if len(split) > 1 {
		args = split[1:]
	}

	// Start process.
	proc, err := StartProcess("", name, args)
	require.NoError(t, err)
	return proc
}

func StartProcess(dir string, name string, args []string) (*Process, error) {
	proc, err := CreateProcess(dir, name, args)
	if err != nil {
		return nil, err
	}
	// cmd start
	if err := proc.Cmd.Start(); err != nil {
		return nil, err
	}
	proc.Pid = proc.Cmd.Process.Pid

	return proc, nil
}

// Same as StartProcess but doesn't start the process
func CreateProcess(dir string, name string, args []string) (*Process, error) {
	var cmd = exec.Command(name, args...) // is not yet started.
	// cmd dir
	if dir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cmd.Dir = pwd
	} else {
		cmd.Dir = dir
	}
	// cmd stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	proc := &Process{
		ExecPath:   name,
		Args:       args,
		StartTime:  time.Now(),
		Cmd:        cmd,
		ExitState:  nil,
		StdinPipe:  stdin,
		StdoutPipe: stdout,
		StderrPipe: stderr,
	}
	return proc, nil
}

type Process struct {
	ExecPath   string
	Args       []string
	Pid        int
	StartTime  time.Time
	EndTime    time.Time
	Cmd        *exec.Cmd        `json:"-"`
	ExitState  *os.ProcessState `json:"-"`
	StdinPipe  io.WriteCloser   `json:"-"`
	StdoutPipe io.ReadCloser    `json:"-"`
	StderrPipe io.ReadCloser    `json:"-"`
}

// ReadAll calls ioutil.ReadAll on the StdoutPipe and StderrPipe.
func (proc *Process) ReadAll() (stdout []byte, stderr []byte, err error) {
	outbz, err := ioutil.ReadAll(proc.StdoutPipe)
	if err != nil {
		return nil, nil, err
	}
	errbz, err := ioutil.ReadAll(proc.StderrPipe)
	if err != nil {
		return nil, nil, err
	}
	return outbz, errbz, nil
}

// wait for the process
func (proc *Process) Wait() {
	err := proc.Cmd.Wait()
	if err != nil {
		// fmt.Printf("Process exit: %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			proc.ExitState = exitError.ProcessState
		}
	}
	proc.ExitState = proc.Cmd.ProcessState
	proc.EndTime = time.Now() // TODO make this goroutine-safe
}

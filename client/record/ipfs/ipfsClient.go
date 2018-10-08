package ipfs

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	tar "github.com/whyrusleeping/tar-utils"
)

const (
	DefaultURL = "localhost:5001"
)

type IpfsClient struct {
	Client *http.Client
	Url    string
	Para   *RequestParameter
}

type RequestParameter struct {
	Command string
	Args    []string
	Opts    map[string]string
	Headers map[string]string
	Body    io.Reader
}

func getIpfsNodeUrl() string {
	if a, err := ma.NewMultiaddr(DefaultURL); err == nil {
		_, url, err := manet.DialArgs(a)
		if err == nil {
			return url
		}
	}
	return DefaultURL
}

func NewIpfsclient() *IpfsClient {
	return &IpfsClient{
		Client: &http.Client{},
		Url:    getIpfsNodeUrl(),
		Para:   &RequestParameter{},
	}
}

func (ipfs *IpfsClient) Add(r io.Reader, pin, rawLeaves bool) (string, error) {

	ipfs.Para.Command = "add"

	var rc io.ReadCloser
	if rclose, ok := r.(io.ReadCloser); ok {
		rc = rclose
	} else {
		rc = ioutil.NopCloser(r)
	}

	// handler expects an array of files
	fr := files.NewReaderFile("", "", rc, nil)
	slf := files.NewSliceFile("", "", []files.File{fr})
	fileReader := files.NewMultiFileReader(slf, true)

	type object struct {
		Hash string
	}
	var recordHash object

	if ipfs.Para.Opts == nil {
		ipfs.Para.Opts = make(map[string]string, 1)
	}
	ipfs.Para.Opts["progress"] = strconv.FormatBool(false)
	ipfs.Para.Opts["pin"] = strconv.FormatBool(pin)
	ipfs.Para.Opts["raw-leaves"] = strconv.FormatBool(rawLeaves)
	ipfs.Para.Body = fileReader

	err := ipfs.Exec(context.Background(), &recordHash)

	return recordHash.Hash, err

}

func (ipfs *IpfsClient) Get(hash, outdir string) error {

	ipfs.Para.Command = "get"

	ipfs.Para.Args = []string{hash}
	if ipfs.Para.Opts == nil {
		ipfs.Para.Opts = make(map[string]string, 1)
	}
	ipfs.Para.Opts["create"] = "true"

	resp, err := ipfs.Send(context.Background())

	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}

	extractor := &tar.Extractor{Path: outdir}
	return extractor.Extract(resp.Output)
}

func (ipfs *IpfsClient) Send(ctx context.Context) (*Response, error) {
	req := NewRequest(ctx, ipfs.Url, ipfs.Para.Command, ipfs.Para.Args...)
	req.Opts = ipfs.Para.Opts
	req.Headers = ipfs.Para.Headers
	req.Body = ipfs.Para.Body
	return req.Send(ipfs.Client)
}

func (ipfs *IpfsClient) Exec(ctx context.Context, res interface{}) error {
	httpRes, err := ipfs.Send(ctx)
	if err != nil {
		return err
	}

	if res == nil {
		httpRes.Close()
		if httpRes.Error != nil {
			return httpRes.Error
		}
		return nil
	}

	return httpRes.Decode(res)
}

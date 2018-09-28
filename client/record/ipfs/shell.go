package ipfs

import (
	"context"
	"io"
	"io/ioutil"
	gohttp "net/http"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	tar "github.com/whyrusleeping/tar-utils"
)

type Shell struct {
	url     string
	httpcli *gohttp.Client
}

func NewShell(url string) *Shell {
	c := &gohttp.Client{
		Transport: &gohttp.Transport{
			Proxy:             gohttp.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	return NewShellWithClient(url, c)
}

func NewShellWithClient(url string, c *gohttp.Client) *Shell {
	if a, err := ma.NewMultiaddr(url); err == nil {
		_, host, err := manet.DialArgs(a)
		if err == nil {
			url = host
		}
	}

	return &Shell{
		url:     url,
		httpcli: c,
	}
}

func (s *Shell) Add(r io.Reader) (string, error) {
	return s.AddWithOpts(r, true, false)
}

type object struct {
	Hash string
}

func (s *Shell) AddWithOpts(r io.Reader, pin bool, rawLeaves bool) (string, error) {
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

	var out object
	return out.Hash, s.Request("add").
		Option("progress", false).
		Option("pin", pin).
		Option("raw-leaves", rawLeaves).
		Body(fileReader).
		Exec(context.Background(), &out)
}

func (s *Shell) Request(command string, args ...string) *RequestBuilder {
	return &RequestBuilder{
		command: command,
		args:    args,
		shell:   s,
	}
}

func (s *Shell) Get(hash, outdir string) error {
	resp, err := s.Request("get", hash).Option("create", true).Send(context.Background())
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

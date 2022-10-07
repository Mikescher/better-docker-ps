package docker

import (
	"better-docker-ps/cli"
	"better-docker-ps/consts"
	pserr "better-docker-ps/fferr"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

func ListContainer(ctx *cli.PSContext) ([]byte, error) {
	if ctx.Opt.Input != nil {
		data, err := os.ReadFile(*ctx.Opt.Input)
		if err != nil {
			return nil, pserr.DirectOutput.Wrap(err, "Failed to read --input file")
		}
		return data, nil
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", ctx.Opt.Socket)
			},
		},
	}

	uri := fmt.Sprintf("%s?1=1", consts.DockerAPIContainerList)

	if ctx.Opt.All {
		uri += "&all=true"
	}

	if ctx.Opt.WithSize {
		uri += "&size=true"
	}

	if ctx.Opt.Filter != nil {
		bin, err := json.Marshal(*ctx.Opt.Filter)
		if err != nil {
			return nil, errorx.InternalError.Wrap(err, "Failed to marshal filter")
		}

		uri += "&filter=" + url.PathEscape(string(bin))
	}

	response, err := client.Get(uri)
	if err != nil {
		return nil, pserr.DirectOutput.Wrap(err, "Call to unix socket failed")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "Failed to read unix socket response")
	}

	return body, nil
}

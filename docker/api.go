package docker

import (
	"better-docker-ps/cli"
	"better-docker-ps/consts"
	"better-docker-ps/pserr"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joomcode/errorx"
)

func ListContainer(ctx *cli.PSContext) ([]byte, error) {
	if ctx.Opt.Input != nil {
		data, err := os.ReadFile(*ctx.Opt.Input)
		if err != nil {
			return nil, pserr.DirectOutput.Wrap(err, "Failed to read --input file")
		}
		return data, nil
	}

	client := newSocketClient(ctx.Opt.GetSocket())

	uri := fmt.Sprintf("%s?1=1", consts.DockerAPIContainerList)

	if ctx.Opt.All {
		uri += "&all=true"
	}

	if ctx.Opt.WithSize {
		uri += "&size=true"
	}

	if ctx.Opt.Limit != -1 {
		uri += "&limit=" + strconv.Itoa(ctx.Opt.Limit)
	}

	if ctx.Opt.Filter != nil {
		bin, err := json.Marshal(*ctx.Opt.Filter)
		if err != nil {
			return nil, errorx.InternalError.Wrap(err, "Failed to marshal filter")
		}

		uri += "&filters=" + url.PathEscape(string(bin))
	}

	response, err := client.Get(uri)
	if err != nil {
		if strings.Contains(err.Error(), "connect: permission denied") {
			return nil, pserr.DirectOutput.Wrap(err, "Call to unix socket failed (permission denied)")
		} else {
			return nil, pserr.DirectOutput.Wrap(err, "Call to unix socket failed")
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "Failed to read unix socket response")
	}

	return body, nil
}

// InspectContainer queries the container-inspect endpoint for a single container.
// This is needed for data that the (cheaper) container-list endpoint does not return, e.g. Config.User.
func InspectContainer(ctx *cli.PSContext, containerID string) (*ContainerInspectSchema, error) {
	client := newSocketClient(ctx.Opt.GetSocket())

	uri := fmt.Sprintf(consts.DockerAPIContainerInspect, containerID)

	response, err := client.Get(uri)
	if err != nil {
		if strings.Contains(err.Error(), "connect: permission denied") {
			return nil, pserr.DirectOutput.Wrap(err, "Call to unix socket failed (permission denied)")
		} else {
			return nil, pserr.DirectOutput.Wrap(err, "Call to unix socket failed")
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "Failed to read unix socket response")
	}

	var data ContainerInspectSchema
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errorx.InternalError.Wrap(err, "Failed to decode Docker API (inspect) response")
	}

	return &data, nil
}

func newSocketClient(socket string) http.Client {
	return http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socket)
			},
		},
	}
}

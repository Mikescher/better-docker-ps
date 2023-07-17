package docker

import (
	"better-docker-ps/cli"
	"strings"
)

var registryPrefixList = []string{
	".com",
	".de",
	".net",
	".io",
	".org",
}

func SplitDockerImage(ctx *cli.PSContext, img string) (string, string, string) {

	resultRegistry := ""
	resultImage := ""
	resultTag := ""

	if v := strings.Split(img, ":"); len(v) > 1 {
		last := v[len(v)-1]
		if !strings.Contains(last, "/") {
			resultTag = last
			img = img[0 : len(img)-len(last)-1]
		}
	}

	if v := strings.Split(img, "/"); len(v) > 1 {
		first := v[0]
		if len(v) == 3 {
			resultRegistry = first
			img = img[len(resultRegistry)+1:]
		} else {
			for _, rpl := range registryPrefixList {
				if strings.HasSuffix(first, rpl) {
					resultRegistry = first
					img = img[len(resultRegistry)+1:]
				}
				break
			}
		}
	}

	resultImage = img

	if resultImage == "sha256" && len(resultTag) == 64 && ctx.Opt.Truncate {
		resultImage = "(sha256)"
		resultTag = resultTag[0:12] + "..."
	}

	return resultRegistry, resultImage, resultTag

}

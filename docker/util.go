package docker

import "strings"

var registryPrefixList = []string{
	".com",
	".de",
	".net",
	".io",
	".org",
}

func SplitDockerImage(img string) (string, string, string) {

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

	return resultRegistry, resultImage, resultTag

}

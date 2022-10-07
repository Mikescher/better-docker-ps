package term

import (
	"golang.org/x/term"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// -> partly copied from [ https://github.com/jwalton/go-supportscolor/tree/master ]

func SupportsColors() bool {

	if isatty := term.IsTerminal(int(os.Stdout.Fd())); !isatty {
		return false
	}

	termenv := os.Getenv("TERM")

	if termenv == "dumb" {
		return false
	}

	if osColorEnabled := enableColor(); !osColorEnabled {
		return false
	}

	if _, ci := os.LookupEnv("CI"); ci {
		var ciEnvNames = []string{"TRAVIS", "CIRCLECI", "APPVEYOR", "GITLAB_CI", "GITHUB_ACTIONS", "BUILDKITE", "DRONE"}
		for _, ciEnvName := range ciEnvNames {
			_, exists := os.LookupEnv(ciEnvName)
			if exists {
				return true
			}
		}

		if os.Getenv("CI_NAME") == "codeship" {
			return true
		}

		return false
	}

	if teamCityVersion, isTeamCity := os.LookupEnv("TEAMCITY_VERSION"); isTeamCity {
		versionRegex := regexp.MustCompile(`^(9\.(0*[1-9]\d*)\.|\d{2,}\.)`)
		if versionRegex.MatchString(teamCityVersion) {
			return true
		}
		return false
	}

	if os.Getenv("COLORTERM") == "truecolor" {
		return true
	}

	if termProgram, termProgramPreset := os.LookupEnv("TERM_PROGRAM"); termProgramPreset {
		switch termProgram {
		case "iTerm.app":
			termProgramVersion := strings.Split(os.Getenv("TERM_PROGRAM_VERSION"), ".")
			version, err := strconv.ParseInt(termProgramVersion[0], 10, 64)
			if err == nil && version >= 3 {
				return true
			}
			return true
		case "Apple_Terminal":
			return true

		default:
			// No default
		}
	}

	var term256Regex = regexp.MustCompile("(?i)-256(color)?$")
	if term256Regex.MatchString(termenv) {
		return true
	}

	var termBasicRegex = regexp.MustCompile("(?i)^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux")

	if termBasicRegex.MatchString(termenv) {
		return true
	}

	if _, colorTerm := os.LookupEnv("COLORTERM"); colorTerm {
		return true
	}

	return false
}

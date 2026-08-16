package tzdata

import (
	// embed is required for the //go:embed directive
	_ "embed"
	"sort"
	"strings"
	// time/tzdata embeds the IANA time zone database into the binary
	_ "time/tzdata"
)

//go:embed zone1970.tab
var zoneTab string

var cachedTimezones []string

func init() {
	lines := strings.Split(zoneTab, "\n")
	tzSet := map[string]bool{
		"UTC":     true,
		"GMT":     true,
		"Etc/UTC": true,
		"Etc/GMT": true,
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 3 {
			tzName := strings.TrimSpace(parts[2])
			if tzName != "" {
				tzSet[tzName] = true
			}
		}
	}

	for tz := range tzSet {
		cachedTimezones = append(cachedTimezones, tz)
	}

	sort.Strings(cachedTimezones)
}

// GetTimezones returns all available IANA timezones supported by time/tzdata
func GetTimezones() []string {
	res := make([]string, len(cachedTimezones))
	copy(res, cachedTimezones)
	return res
}

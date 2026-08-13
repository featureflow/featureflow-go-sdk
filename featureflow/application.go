package featureflow

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Contract slug (featureflow-client-sdk-testbed/CONTRACT.md): lowercase [a-z0-9._-],
// max 64 chars. The server sanitises defensively; the SDK validates strictly so a typo
// is a visible warning here rather than a silently mangled tag there.
var applicationPattern = regexp.MustCompile(`^[a-z0-9._-]{1,64}$`)

// sanitiseApplication validates the configured application tag. Case is forgiven
// (lowercased); anything else invalid is dropped with a warning, and no
// X-Featureflow-Application header is sent at all.
func sanitiseApplication(raw string, logger *log.Logger) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return ""
	}
	if !applicationPattern.MatchString(value) {
		logger.Println(
			LOG_WARN,
			fmt.Sprintf("ignoring application %q — must be lowercase a-z, 0-9, dot, underscore or hyphen, max 64 chars", raw),
		)
		return ""
	}
	return value
}

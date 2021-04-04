package context

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateReqID creates a UUID ID with the "-" removed
func GenerateReqID() string {
	reqID := uuid.New().String()
	return strings.ReplaceAll(reqID, "-", "")
}

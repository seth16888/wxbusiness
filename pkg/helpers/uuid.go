package helpers

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() string {
	rt := uuid.NewString()
	if rt == "" {
		panic("uuid is empty")
	}
	return strings.ReplaceAll(rt, "-", "")
}

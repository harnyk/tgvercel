package vapi

import (
	"fmt"
)

type Target string

const (
	TargetProduction  Target = "production"
	TargetPreview     Target = "preview"
	TargetDevelopment Target = "development"
)

func NewTarget(target string) (Target, error) {
	t := Target(target)
	switch t {
	case TargetProduction, TargetPreview, TargetDevelopment:
		return t, nil
	default:
		return "", fmt.Errorf("invalid target: %s", target)
	}
}

type Options struct {
	Token string
}

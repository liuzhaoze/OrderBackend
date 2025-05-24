package cqrs

import (
	"fmt"
	"strings"
)

func getName[CQ any](cq CQ) string {
	return strings.Split(fmt.Sprintf("%T", cq), ".")[1]
}

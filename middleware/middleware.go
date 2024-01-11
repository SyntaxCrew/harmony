package middleware

import "github.com/SyntaxCrew/harmony"

type (
	Skipper func(ctx harmony.Context) bool
)

func defaultSkipper(ctx harmony.Context) bool {
	return false
}

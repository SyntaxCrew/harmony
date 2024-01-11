package middleware

import "github.com/SyntaxCrew/harmony"

type (
	// Skipper defines a function to skip middleware.
	Skipper func(ctx harmony.Context) bool
)

func defaultSkipper(ctx harmony.Context) bool {
	return false
}

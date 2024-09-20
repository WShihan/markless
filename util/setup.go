package util

import "markless/injection"

var (
	Env injection.Env
)

func InitENV(env *injection.Env) {
	Env = *env
	InitLogger()
}

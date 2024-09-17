package util

import "markee/injection"

var (
	Env injection.Env
)

func InitENV(env injection.Env) {
	Env = env
}

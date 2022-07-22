package x_utils

import (
	"log"
	"runtime/debug"
)

func RecoverErr() {
	if err := recover(); err != nil {
		log.Printf("Err is [%v]\n Stack:[%v]", err, string(debug.Stack()))
	}
}

package generators

import (
	"github.com/golang-module/carbon"
)

func TimeStampGenerator() any {
	//TODO range generator
	return carbon.Now().ToStdTime()
}

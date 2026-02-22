package pdf

import (
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

func GetEngine() core.Maroto {
	engine := maroto.New()

	return engine
}

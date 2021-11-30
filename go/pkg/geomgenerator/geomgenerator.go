package geomgenerator

import (
	"math/rand"

	"github.com/muzammilar/geomrpc/protos/shape"
)

func Cuboid() *shape.Cuboid {
	return &shape.Cuboid{
		Id: &shape.Identifier{
			Id: int64(rand.Uint32()),
		},
		Length: int64(10 + rand.Uint32()%25),
		Width:  int64(1 + rand.Uint32()%10),
		Height: int64(1 + rand.Uint32()%25),
	}
}

/*
 * The package implements a basic examle of a geometry grpc server
 */

package geomserver

/*
 * HTTP Server Package
 */
import (
	"github.com/muzammilar/geometric-shapes/protos/shapecalc"
	"github.com/sirupsen/logrus"
)

/*
 * Private Structs (or Public)
 */
type GeometryServer struct {
	shapecalc.UnimplementedGeometryServer
	logger *logrus.Logger
}

/*
 * Private Functions
 */

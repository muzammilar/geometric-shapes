/*
 * The package implements a basic examle of a geometry grpc server
 */

package geomserver

/*
 * Geometry Server Package
 */
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/muzammilar/geomrpc/pkg/grpcserver"
	"github.com/muzammilar/geomrpc/protos/serviceinfo"
	"github.com/muzammilar/geomrpc/protos/shape"
	"github.com/muzammilar/geomrpc/protos/shapecalc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
 * Private Structs (or Public)
 */

// GeometryServer implements two different services in the same server. This is allowed since the service names are unique
// Generally it's best to have two different structs implement them, however, this is a POC for experimentation, so it's fine
type GeometryServer struct {
	// embed the Geometry Server
	shapecalc.UnimplementedGeometryServer
	// embed the Info Server
	shapecalc.UnimplementedInfoServer
	// Other internal use variables
	logger *logrus.Logger  // a shared logger (can be a bottleneck)
	wg     *sync.WaitGroup // a wait group - since the server's serve is blocking (and shutdown closes all workers), a wait group is not needed
	// server information
	name    string // name of the server
	version string // version of the server
}

/*
 * Functions - Shared by both Service Server
 */

// If the following function is not implemented, there would be ambiguity
// GeometryServer.Version is ambiguous
// cannot use geometryHandler (type *GeometryServer) as type shapecalc.GeometryServer in argument to shapecalc.RegisterGeometryServer:
//        *GeometryServer does not implement shapecalc.GeometryServer (missing Version method)
func (g *GeometryServer) Version(ctx context.Context, e *emptypb.Empty) (*serviceinfo.Info, error) {
	return &serviceinfo.Info{
		Server: &serviceinfo.Server{
			Name: fmt.Sprintf("%s (%T)", g.name, g),
		},
		Version: &serviceinfo.Version{
			Name: g.version,
		},
	}, nil
}

/*
 * Functions - Info Service Server
 */

func (g *GeometryServer) ComputeRectangleArea(ctx context.Context, r *shape.Rectangle) (*shape.ShapeInfo_Mesurement, error) {
	// get client information
	clientAddr := grpcserver.GetRemoteHostFromContext(ctx)
	errTemplate := "Error occured while computing GeometryServer/ComputeRectangleArea to '%s': %s"

	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		g.logger.Debugf(errTemplate, clientAddr, err.Error())
		return nil, err
	}
	// compute area
	var m = new(shape.ShapeInfo_Mesurement)
	m.Name = shape.ShapeInfo_AREA
	m.Value = float64(r.Length * r.Width)
	return m, nil
}

func (g *GeometryServer) ListRectangleCoordinates(r *shape.Rectangle, stream shapecalc.Geometry_ListRectangleCoordinatesServer) error {
	// get client information
	clientAddr := grpcserver.GetRemoteHostFromContext(stream.Context())
	errTemplate := "Error occured while streaming GeometryServer/ListRectangleCoordinates to '%s': %s"

	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		g.logger.Debugf(errTemplate, clientAddr, err.Error())
		return err
	}
	// list coordinates (all points on the rectangle)
	// note: since length of 2 means, that both 0 and 2 will be included [0,2] in the coordinates
	var x, y int64
	for x = 0; x <= r.Length; x++ {
		for y = 0; y <= r.Width; y++ {
			// get coordinates
			pc := &shape.PlanarCoordinates{
				Id:    r.Id,
				Shape: shape.ShapeInfo_RECTANGLE,
				X:     x,
				Y:     y,
			}
			// send and check for errors
			if err := stream.Send(pc); err != nil {
				g.logger.Infof(errTemplate, clientAddr, err.Error())
				return err
			}
		}
	}
	// no error
	return nil
}

/*
 * Functions - Geometry Service Server
 */

func (g *GeometryServer) RectangleInfo(ctx context.Context, r *shape.Rectangle) (*shape.ShapeInfo, error) {
	// get client information
	clientAddr := grpcserver.GetRemoteHostFromContext(ctx)
	errTemplate := "Error occured while computing GeometryServer/RectangleInfo to '%s': %s"

	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		g.logger.Debugf(errTemplate, clientAddr, err.Error())
		return nil, err
	}
	// compute info
	area, err := g.ComputeRectangleArea(ctx, r)
	if err != nil {
		g.logger.Debugf(errTemplate, clientAddr, err.Error())
		return nil, err
	}
	perimeter, err := g.ComputeRectanglePerimeter(ctx, r)
	if err != nil {
		g.logger.Debugf(errTemplate, clientAddr, err.Error())
		return nil, err
	}
	measurements := make([]*shape.ShapeInfo_Mesurement, 3)
	measurements = append(measurements, area)
	measurements = append(measurements, perimeter)

	si := &shape.ShapeInfo{
		Id:          r.Id,
		Shape:       shape.ShapeInfo_RECTANGLE,
		Mesurements: measurements,
		Timestamp:   timestamppb.New(time.Now()),
	}
	g.logger.Tracef("Computing GeometryServer/RectangleInfo and sent to '%s': %#v", clientAddr, si)
	return si, nil
}

func (g *GeometryServer) ComputeRectanglePerimeter(ctx context.Context, r *shape.Rectangle) (*shape.ShapeInfo_Mesurement, error) {
	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		return nil, err
	}
	// compute area
	var m = new(shape.ShapeInfo_Mesurement)
	m.Name = shape.ShapeInfo_PERIMETER
	m.Value = 2 * float64(r.Length+r.Width)
	return m, nil
}

/*
 * Private Functions
 */

// validateRectangleDimensions makes sure that the dimensions of a rectangle are positive only
// note: We do not discuss negative areas and vector coordinates for now. We only discuss scalars.
func validateRectangleDimensions(r *shape.Rectangle) error {
	if r.Length < 0 {
		return fmt.Errorf("The length field can not be negative: %#v", r.Length)
	}
	if r.Width < 0 {
		return fmt.Errorf("The width field can not be negative: %#v", r.Width)
	}
	return nil
}

// validateCuboidDimensions makes sure that the dimensions of a cuboid are positive only
func validateCuboidDimensions(c *shape.Cuboid) error {
	if c.Length < 0 {
		return fmt.Errorf("The length field can not be negative: %#v", c.Length)
	}
	if c.Width < 0 {
		return fmt.Errorf("The width field can not be negative: %#v", c.Width)
	}
	if c.Height < 0 {
		return fmt.Errorf("The height field can not be negative: %#v", c.Width)
	}
	return nil
}

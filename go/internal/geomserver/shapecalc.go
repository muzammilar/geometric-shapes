/*
 * The package implements a basic examle of a geometry grpc server
 */

package geomserver

/*
 * HTTP Server Package
 */
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/muzammilar/geometric-shapes/protos/shape"
	"github.com/muzammilar/geometric-shapes/protos/shapecalc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
 * Private Structs (or Public)
 */

// GeometryServer implements two different services in the same server. This is allowed since the service names are unique
type GeometryServer struct {
	// embed the Geometry Server
	shapecalc.UnimplementedGeometryServer
	// embed the Info Server
	shapecalc.UnimplementedInfoServer
	// Other internal use variables
	logger *logrus.Logger  // a shared logger (can be a bottleneck)
	wg     *sync.WaitGroup // a wait group to track all the request
}

/*
 * Functions - Info Service Server
 */

func (g *GeometryServer) ComputeRectangleArea(ctx context.Context, r *shape.Rectangle) (*shape.ShapeInfo_Mesurement, error) {
	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		return nil, err
	}
	// compute area
	var m = new(shape.ShapeInfo_Mesurement)
	m.Name = shape.ShapeInfo_AREA
	m.Value = float64(r.Length * r.Width)
	return m, nil
}

func (g *GeometryServer) ListRectangleCoordinates(r *shape.Rectangle, stream shapecalc.Geometry_ListRectangleCoordinatesServer) error {
	// validate data
	if err := validateRectangleDimensions(r); err != nil {
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
	// validate data
	if err := validateRectangleDimensions(r); err != nil {
		return nil, err
	}
	// compute info
	area, err := g.ComputeRectangleArea(ctx, r)
	if err != nil {
		return nil, err
	}
	perimeter, err := g.ComputeRectanglePerimeter(ctx, r)
	if err != nil {
		return nil, err
	}
	measurements := make([]*shape.ShapeInfo_Mesurement, 3)
	measurements = append(measurements, area)
	measurements = append(measurements, perimeter)

	return &shape.ShapeInfo{
		Id:          r.Id,
		Shape:       shape.ShapeInfo_RECTANGLE,
		Mesurements: measurements,
		Timestamp:   timestamppb.New(time.Now()),
	}, nil

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

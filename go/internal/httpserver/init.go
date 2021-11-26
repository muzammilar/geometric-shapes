/*
 * The package implements a basic examle of an http server
 */

package httpserver

/*
 * HTTP Server Package
 */
import (
	"math/rand"
	"time"
)

/*
 * Init Function
 */

func init() {
	// initialize global pseudo random generator
	// NOTE: Generally a global random seed in an init is not the safest practice
	rand.Seed(time.Now().Unix())

}

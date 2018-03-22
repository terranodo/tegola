package slippy

import (
	"math"
	"github.com/go-spatial/tegola/maths"
	"github.com/go-spatial/tegola"
	"github.com/go-spatial/tegola/internal/log"
)

func NewTile(z, x, y uint, buffer float64, srid uint64) *Tile {
	return &Tile{
		z:      z,
		x:      x,
		y:      y,
		Buffer: buffer,
		SRID:   srid,
	}
}

func DegToNum(zoom uint, lat, lon float64) (x, y uint) {
	lat_rad := maths.RadToDeg(lat)
	n := float64(maths.Exp2(uint64(zoom)))

	x = uint(n * (lon + 180.0) / 360.0)
	y = uint(n *
		(1.0 - math.Log(
			math.Tan(lat_rad)+(1/math.Cos(lat_rad))) / 2.0))

	return
}

func NewTileLatLon(z uint, lat, lon, buffer float64, srid uint64) *Tile {
	x, y := DegToNum(z, lat, lon)
	return &Tile{
		z:      z,
		x:      x,
		y:      y,
		Buffer: buffer,
		SRID:   srid,
	}
}

type Tile struct {
	// zoom
	z uint
	// column
	x uint
	// row
	y uint
	// buffer will add a buffer to the tile bounds. this buffer is expected to use the same units as the SRID
	// of the projected tile (i.e. WebMercator = pixels, 3395 = meters)
	Buffer float64
	// spatial reference id
	SRID uint64
}

func (t *Tile) ZXY() (uint, uint, uint) {
	return t.z, t.x, t.y
}

// TODO(arolek): return geom.Extent once it has been refactored
// TODO(arolek): support alternative SRIDs. Currently this assumes 3857
// Extent will return the tile extent excluding the tile's buffer and the Extent's SRID
func (t *Tile) Extent() (extent [2][2]float64, srid uint64) {
	max := 20037508.34

	//	resolution
	res := (max * 2) / math.Exp2(float64(t.z))

	//	unbuffered extent
	return [2][2]float64{
		{
			-max + (float64(t.x) * res), // MinX
			max - (float64(t.y) * res),  // Miny
		},
		{
			-max + (float64(t.x) * res) + res, // MaxX
			max - (float64(t.y) * res) - res,  // MaxY
		},
	}, t.SRID
}

// BufferedExtent will return the tile extent including the tile's buffer and the Extent's SRID
func (t *Tile) BufferedExtent() (bufferedExtent [2][2]float64, srid uint64) {
	extent, _ := t.Extent()

	// TODO(arolek): the following value is hard coded for MVT, but this concept needs to be abstracted to support different projections
	mvtTileWidthHeight := 4096.0
	// the bounds / extent
	mvtTileExtent := [2][2]float64{{0 - t.Buffer, 0 - t.Buffer}, {mvtTileWidthHeight + t.Buffer, mvtTileWidthHeight + t.Buffer}}

	xspan := extent[1][0] - extent[0][0]
	yspan := extent[1][1] - extent[0][1]

	bufferedExtent[0][0] = (mvtTileExtent[0][0] * xspan / mvtTileWidthHeight) + extent[0][0]
	bufferedExtent[0][1] = (mvtTileExtent[0][1] * yspan / mvtTileWidthHeight) + extent[0][1]
	bufferedExtent[1][0] = (mvtTileExtent[1][0] * xspan / mvtTileWidthHeight) + extent[0][0]
	bufferedExtent[1][1] = (mvtTileExtent[1][1] * yspan / mvtTileWidthHeight) + extent[0][1]

	return bufferedExtent, t.SRID
}

// calls f on every vertically related to t at the specified zoom
// the method name can be misleading because it also handles parents
func (t *Tile) RangeChildren(zoom uint, f func(*Tile) error) error {
	if zoom <= t.z {
		mag := t.z - zoom
		arg := NewTile(zoom, t.x>>mag, t.y>>mag, t.Buffer, t.SRID)
		return f(arg)
	}

	mag := zoom - t.z
	delta := uint(maths.Exp2(uint64(mag)))

	leastX := t.x << mag
	leastY := t.y << mag

	log.Info("info: ", mag, delta, leastY, leastY)

	for x := leastX; x < leastX+delta; x++ {
		for y := leastY; y < leastY+delta; y++ {
			err := f(NewTile(zoom, x, y, 0, tegola.WebMercator))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
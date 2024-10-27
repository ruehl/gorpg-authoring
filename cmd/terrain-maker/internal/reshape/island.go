package reshape

import (
	"github.com/ruehl/gorpg-authoring/cmd/terrain-maker/internal/noise"
)

func ReshapeIsland(elevationData noise.NoiseData, mix float64) noise.NoiseData {
	result := noise.NoiseData{
		Data:   make([]float64, elevationData.Width*elevationData.Height),
		Width:  elevationData.Width,
		Height: elevationData.Height,
	}

	max := 0.0
	min := 1.0

	for y := 0; y < elevationData.Height; y++ {
		for x := 0; x < elevationData.Width; x++ {
			offset := y*elevationData.Width + x
			elevation := elevationData.Data[offset]
			nx := 2*float64(x)/float64(elevationData.Width) - 1
			ny := 2*float64(y)/float64(elevationData.Height) - 1
			d := distanceFunctionSquareBump(nx, ny)
			result.Data[offset] = lerp(elevation, 1-d, mix)

			if result.Data[offset] > max {
				max = result.Data[offset]
			}

			if result.Data[offset] < min {
				min = result.Data[offset]
			}
		}
	}

	result.MinValue = min
	result.MaxValue = max

	return result
}

func distanceFunctionSquareBump(nx float64, ny float64) float64 {
	return 1 - (1-nx*nx)*(1-ny*ny)
}

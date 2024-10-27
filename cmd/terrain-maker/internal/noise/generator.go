package noise

import (
	"math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type NoiseGeneratorParams struct {
	Frequency float64   // Frequency of the noise
	Octaves   []float64 // Octaves to apply on the noise
	Power     float64   // Coefficient of the redistribution power function
	Seed      int64     // Seed for the random number generator
}

func DefaultNoiseGenertoParams(frequency float64) NoiseGeneratorParams {
	return NoiseGeneratorParams{
		Frequency: frequency,
		Octaves:   []float64{1.0, 0.5, 0.25, 0.125},
		Power:     1.2,
		Seed:      rand.Int63(),
	}
}

type NoiseData struct {
	Data     []float64
	Width    int
	Height   int
	MinValue float64
	MaxValue float64
}

func Generate(width, height int, params NoiseGeneratorParams) NoiseData {
	result := NoiseData{
		Data:   make([]float64, width*height),
		Width:  width,
		Height: height,
	}

	// Setup noise generators for each octave
	generators := make([]opensimplex.Noise, 0, len(params.Octaves))
	seed := params.Seed
	for i := 0; i < len(params.Octaves); i++ {
		generators = append(generators, opensimplex.NewNormalized(seed+int64(i)))
	}

	// Calculate sum of octaves
	sumOctaves := 0.0
	for _, octave := range params.Octaves {
		sumOctaves += octave
	}

	max := 0.0
	min := 1.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := float64(x)/float64(width) - 0.5
			ny := float64(y)/float64(height) - 0.5
			offset := y*width + x
			e := 0.0

			for index, octave := range params.Octaves {
				e += octave * float64(generators[index].Eval2(
					params.Frequency*nx*1.0/octave,
					params.Frequency*ny*1.0/octave,
				))
			}

			result.Data[offset] = math.Pow(e/sumOctaves, params.Power)

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

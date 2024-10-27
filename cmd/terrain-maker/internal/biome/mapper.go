package biome

import (
	"image/color"
	"sort"

	"github.com/ruehl/gorpg-authoring/cmd/terrain-maker/internal/noise"
)

type BiomeDefinition struct {
	Id           uint8
	Name         string
	MaxElevation float64
	MaxMoisture  float64
	Color        color.NRGBA
}

type biomeDefinitions []BiomeDefinition

func (a biomeDefinitions) Len() int { return len(a) }
func (a biomeDefinitions) Less(i, j int) bool {
	if a[i].MaxElevation < a[j].MaxElevation {
		return true
	} else if a[i].MaxElevation == a[j].MaxElevation {
		return a[i].MaxMoisture < a[j].MaxMoisture
	} else {
		return false
	}
}
func (a biomeDefinitions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func Map(
	elevationData noise.NoiseData,
	moistureData noise.NoiseData,
	biomes []BiomeDefinition,
) (biomeMap []uint8, colorMap []color.NRGBA) {
	biomeMap = make([]uint8, elevationData.Width*elevationData.Height)
	colorMap = make([]color.NRGBA, elevationData.Width*elevationData.Height)

	// We need to remap elevations into the range [0, 1]
	scale := 1.0 / (elevationData.MaxValue - elevationData.MinValue)

	// Sort the biome definitions by elevation and moisture
	// This is necessary to ensure that the biomes are applied in the correct order
	// The order of the biome definitions is important because the first matching biome will be used.
	sortedBiomes := make(biomeDefinitions, len(biomes))
	copy(sortedBiomes, biomes)
	sort.Sort(sortedBiomes)

	for y := 0; y < elevationData.Height; y++ {
		for x := 0; x < elevationData.Width; x++ {
			offset := y*elevationData.Width + x
			elevation := elevationData.Data[offset]*scale - elevationData.MinValue
			moisture := moistureData.Data[offset]
			biome := pickBiomeByElevationAndMoisture(elevation, moisture, sortedBiomes)
			biomeMap[offset] = biome.Id
			colorMap[offset] = biome.Color
		}
	}

	return biomeMap, colorMap
}

func pickBiomeByElevationAndMoisture(elevation, moisture float64, defs []BiomeDefinition) BiomeDefinition {
	for _, biome := range defs {
		if elevation <= biome.MaxElevation && moisture <= biome.MaxMoisture {
			return biome
		}
	}

	// No biome found, return the last one
	return defs[len(defs)-1]
}

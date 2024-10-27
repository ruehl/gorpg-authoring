package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/ruehl/gorpg-authoring/cmd/terrain-maker/internal/biome"
	"github.com/ruehl/gorpg-authoring/cmd/terrain-maker/internal/noise"
	"github.com/ruehl/gorpg-authoring/cmd/terrain-maker/internal/reshape"
)

const width, height = 1024, 1024
const freq = 8.0

func main() {
	// Generate elevation data
	elevationParams := noise.DefaultNoiseGenertoParams(freq)
	elevationData := noise.Generate(width, height, elevationParams)

	// Rehape elevation data into an island shape
	reshapedNoise := reshape.ReshapeIsland(elevationData, 0.4)

	// Generate moisture data
	moistureParams := noise.DefaultNoiseGenertoParams(freq)
	moistureData := noise.Generate(width, height, moistureParams)

	// Define biomes
	biomeDefinitions := []biome.BiomeDefinition{
		{
			Id:           0,
			Name:         "Deep Sea",
			MaxElevation: 0.40,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0x02, G: 0x2F, B: 0x8e, A: 255},
		},
		{
			Id:           1,
			Name:         "Coastal Water",
			MaxElevation: 0.5,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0x1c, G: 0x70, B: 0xc8, A: 255},
		},
		{
			Id:           2,
			Name:         "Beach",
			MaxElevation: 0.52,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0xc2, G: 0xb2, B: 0x80, A: 255},
		},
		{
			Id:           3,
			Name:         "Desert",
			MaxElevation: 0.75,
			MaxMoisture:  0.15,
			Color:        color.NRGBA{R: 0xfa, G: 0xd5, B: 0xa5, A: 255},
		},
		{
			Id:           4,
			Name:         "Grassland",
			MaxElevation: 0.75,
			MaxMoisture:  0.4,
			Color:        color.NRGBA{R: 0x3f, G: 0x9b, B: 0x0b, A: 255},
		},
		{
			Id:           5,
			Name:         "Forest",
			MaxElevation: 0.75,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0x2e, G: 0x6f, B: 0x40, A: 255},
		},
		{
			Id:           6,
			Name:         "Tundra",
			MaxElevation: 0.85,
			MaxMoisture:  0.5,
			Color:        color.NRGBA{R: 0x8c, G: 0x9c, B: 0x5c, A: 255},
		},
		{
			Id:           7,
			Name:         "Highlands",
			MaxElevation: 0.85,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0x7a, G: 0x94, B: 0x61, A: 255},
		},
		{
			Id:           8,
			Name:         "Mountains",
			MaxElevation: 1.0,
			MaxMoisture:  0.6,
			Color:        color.NRGBA{R: 64, G: 64, B: 64, A: 255},
		},
		{
			Id:           8,
			Name:         "Snow",
			MaxElevation: 1.0,
			MaxMoisture:  1.0,
			Color:        color.NRGBA{R: 0xff, G: 0xfa, B: 0xfa, A: 255},
		},
	}

	// Map elevations and moisture to biomes
	_, colorMap := biome.Map(reshapedNoise, moistureData, biomeDefinitions)

	if err := exportImage(colorMap, width, height); err != nil {
		log.Fatalf("could not export elevation map: %v", err)
	}
}

func exportImage(data []color.NRGBA, w, h int) error {
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, data[y*w+x])
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		return err
	}

	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}

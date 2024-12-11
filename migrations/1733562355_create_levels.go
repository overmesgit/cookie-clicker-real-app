package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"log"
)

type Level struct {
	name        string
	baseRate    int
	baseCost    int
	description string
}

var levels = []Level{
	{
		name:        "Quantum Fluctuations",
		baseRate:    1,
		baseCost:    10,
		description: "The smallest possible energy variations. Starting point of universe",
	},
	{
		name:        "Elementary Particles",
		baseRate:    10,
		baseCost:    100,
		description: "Quarks, electrons, neutrinos. Basic building blocks",
	},
	{
		name:        "Atomic Nuclei",
		baseRate:    100,
		baseCost:    1000,
		description: "Protons and neutrons forming. First atomic cores",
	},
	{
		name:        "Atoms",
		baseRate:    1000,
		baseCost:    10000,
		description: "Hydrogen and Helium formation. First complete atoms",
	},
	{
		name:        "Molecular Clouds",
		baseRate:    10000,
		baseCost:    100000,
		description: "Gas clouds in space. Building blocks of stars",
	},
	{
		name:        "Stars",
		baseRate:    100000,
		baseCost:    1000000,
		description: "Nuclear fusion begins. Energy generation",
	},
	{
		name:        "Supernovas",
		baseRate:    1000000,
		baseCost:    10000000,
		description: "Heavy element creation. Spreading materials across space",
	},
	{
		name:        "Solar Systems",
		baseRate:    10000000,
		baseCost:    100000000,
		description: "Stars with planets. Matter organization",
	},
	{
		name:        "Galaxies",
		baseRate:    100000000,
		baseCost:    1000000000,
		description: "Star clusters. Massive matter structures",
	},
	{
		name:        "Galaxy Clusters",
		baseRate:    1000000000,
		baseCost:    10000000000,
		description: "Multiple galaxies. Large scale structures",
	},
	{
		name:        "Observable Universe",
		baseRate:    10000000000,
		baseCost:    100000000000,
		description: "Everything visible. Maximum known scale",
	},
}

func init() {
	m.Register(
		func(app core.App) error {
			collection, err := app.FindCollectionByNameOrId("levels")
			if err != nil {
				log.Printf(err.Error())
				return err
			}

			for _, level := range levels {
				record := core.NewRecord(collection)
				record.Set("rate", level.baseRate)
				record.Set("cost", level.baseCost)
				record.Set("name", level.name)
				record.Set("description", level.description)
				err = app.Save(record)
				if err != nil {
					log.Printf(err.Error())
					return err
				}
			}
			return nil

		}, func(app core.App) error {
			// add down queries...

			return nil
		},
	)
}

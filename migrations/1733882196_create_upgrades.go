package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"log"
)

type Upgrade struct {
	name        string
	baseRate    int
	baseCost    int
	description string
	tag         string
}

var upgrades = []Upgrade{
	{
		name:        "Quantum Tunneling",
		baseRate:    1,
		description: "Harness quantum uncertainty for more efficient clicking",
		tag:         "quantum-tunneling",
	},
	{
		name:        "Hadron Accelerator",
		baseRate:    10,
		description: "Channel particle formation energy into clicks",
		tag:         "hadron-accelerator",
	},

	{
		name:        "Nuclear Fusion Enhancement",
		baseRate:    25,
		description: "Convert nuclear binding energy to click power",
		tag:         "nuclear-fusion",
	},

	{
		name:        "Atomic Resonator",
		baseRate:    50,
		description: "Amplify clicks through electron shell transitions",
		tag:         "atomic-resonator",
	},

	{
		name:        "Stellar Forge",
		baseRate:    100,
		description: "Harness star formation energy for clicking",
		tag:         "stellar-forge",
	},

	{
		name:        "Galactic Core Tap",
		baseRate:    250,
		description: "Channel supermassive black hole energy",
		tag:         "galactic-core",
	},

	{
		name:        "Solar Wind Collector",
		baseRate:    500,
		description: "Convert stellar radiation into click power",
		tag:         "solar-wind",
	},

	{
		name:        "Dark Matter Manipulator",
		baseRate:    1000,
		description: "Utilize invisible mass for enhanced clicking",
		tag:         "dark-matter",
	},

	{
		name:        "Red Giant Amplifier",
		baseRate:    2500,
		description: "Harvest dying star energy",
		tag:         "red-giant",
	},

	{
		name:        "Hawking Radiation Collector",
		baseRate:    5000,
		description: "Convert black hole evaporation to clicks",
		tag:         "hawking-radiation",
	},

	{
		name:        "Vacuum Energy Extractor",
		baseRate:    10000,
		description: "Harvest the final remnants of universal energy",
		tag:         "vacuum-energy",
	},
}

func init() {
	m.Register(
		func(app core.App) error {
			collection, err := app.FindCollectionByNameOrId("upgrades")
			if err != nil {
				log.Printf(err.Error())
				return err
			}

			for _, upgrade := range upgrades {
				record := core.NewRecord(collection)
				record.Set("rate", upgrade.baseRate)
				record.Set("name", upgrade.name)
				record.Set("tag", upgrade.tag)
				record.Set("description", upgrade.description)
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

package dto

const PriceMultiplier = 10000

// Property Type - Price
type PTPrice struct {
	Amount   int64  `json:"amount" validate:"gte=0"`
	Currency string `json:"currency" validate:"iso4217"`
}

// MeasureMultiplier scales the amount of a measured property type (mass, volume) the same way
// PriceMultiplier scales a price: the amount is stored as an integer in the unit the user picked,
// multiplied by this, so four decimal places survive without floating point drift.
const MeasureMultiplier = 10000

// Property Type - Mass
type PTMass struct {
	Amount int64  `json:"amount" validate:"gte=0"`
	Unit   string `json:"unit" validate:"oneof=mg g kg"`
}

// Property Type - Volume
type PTVolume struct {
	Amount int64  `json:"amount" validate:"gte=0"`
	Unit   string `json:"unit" validate:"oneof=mL L"`
}

// MassUnitFactors and VolumeUnitFactors convert an amount into the canonical base unit for its
// dimension - milligram for mass, millilitre for volume - so amounts recorded in different units
// can be summed: baseAmount = amount * factor, still scaled by MeasureMultiplier. Keep the keys in
// sync with the oneof tags on PTMass.Unit and PTVolume.Unit.
var (
	MassUnitFactors   = map[string]int64{"mg": 1, "g": 1_000, "kg": 1_000_000}
	VolumeUnitFactors = map[string]int64{"mL": 1, "L": 1_000}
)

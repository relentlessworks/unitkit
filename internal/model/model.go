package model

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Unit represents a single unit of measurement.
type Unit struct {
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	Factor   float64 `json:"factor"`
	Offset   float64 `json:"offset"`
	Category string  `json:"category"`
}

// Category holds metadata about a measurement category.
type Category struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BaseUnit    string   `json:"base_unit"`
	Units       []string `json:"units"`
}

// ConversionResult holds the result of a unit conversion.
type ConversionResult struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Value    float64 `json:"value"`
	Result   float64 `json:"result"`
	Category string  `json:"category"`
}

// unitRegistry maps exact unit symbols to Unit definitions.
var unitRegistry map[string]Unit

// unitLowerMap maps lowercase unit symbols to all matching units (for case-insensitive fallback).
var unitLowerMap map[string][]Unit

// categoryRegistry maps category names (lowercase) to Category definitions.
var categoryRegistry map[string]Category

// unitToCategory maps exact unit symbols to their category name.
var unitToCategory map[string]string

func init() {
	unitRegistry = make(map[string]Unit)
	unitLowerMap = make(map[string][]Unit)
	categoryRegistry = make(map[string]Category)
	unitToCategory = make(map[string]string)
	registerAll()
}

// registerCategory adds a category and all its units to the registries.
func registerCategory(name, description, baseUnit string, units []Unit) {
	symbols := make([]string, len(units))
	for i, u := range units {
		u.Category = name
		key := u.Symbol
		unitRegistry[key] = u
		unitToCategory[key] = name
		lowerKey := strings.ToLower(u.Symbol)
		unitLowerMap[lowerKey] = append(unitLowerMap[lowerKey], u)
		symbols[i] = u.Symbol
	}
	categoryRegistry[strings.ToLower(name)] = Category{
		Name:        name,
		Description: description,
		BaseUnit:    baseUnit,
		Units:       symbols,
	}
}

func registerAll() {
	// Length — base: meter
	registerCategory("length", "Distance and length measurements", "m", []Unit{
		{Symbol: "m", Name: "meter", Factor: 1},
		{Symbol: "km", Name: "kilometer", Factor: 1000},
		{Symbol: "cm", Name: "centimeter", Factor: 0.01},
		{Symbol: "mm", Name: "millimeter", Factor: 0.001},
		{Symbol: "um", Name: "micrometer", Factor: 1e-6},
		{Symbol: "nm", Name: "nanometer", Factor: 1e-9},
		{Symbol: "mi", Name: "mile", Factor: 1609.344},
		{Symbol: "yd", Name: "yard", Factor: 0.9144},
		{Symbol: "ft", Name: "foot", Factor: 0.3048},
		{Symbol: "in", Name: "inch", Factor: 0.0254},
		{Symbol: "nmi", Name: "nautical mile", Factor: 1852},
		{Symbol: "ly", Name: "light year", Factor: 9.4607304725808e15},
		{Symbol: "au", Name: "astronomical unit", Factor: 1.495978707e11},
		{Symbol: "pc", Name: "parsec", Factor: 3.0856775814913673e16},
	})

	// Mass — base: kilogram
	registerCategory("mass", "Mass and weight measurements", "kg", []Unit{
		{Symbol: "kg", Name: "kilogram", Factor: 1},
		{Symbol: "g", Name: "gram", Factor: 0.001},
		{Symbol: "mg", Name: "milligram", Factor: 1e-6},
		{Symbol: "ug", Name: "microgram", Factor: 1e-9},
		{Symbol: "t", Name: "metric ton", Factor: 1000},
		{Symbol: "lb", Name: "pound", Factor: 0.45359237},
		{Symbol: "oz", Name: "ounce", Factor: 0.028349523125},
		{Symbol: "st", Name: "stone", Factor: 6.35029318},
		{Symbol: "ton_us", Name: "US ton (short ton)", Factor: 907.18474},
		{Symbol: "ton_uk", Name: "UK ton (long ton)", Factor: 1016.0469088},
		{Symbol: "ct", Name: "carat", Factor: 0.0002},
		{Symbol: "gr", Name: "grain", Factor: 6.479891e-5},
	})

	// Temperature — base: Celsius (uses special conversion)
	registerCategory("temperature", "Temperature measurements", "C", []Unit{
		{Symbol: "C", Name: "Celsius", Factor: 1, Offset: 0},
		{Symbol: "F", Name: "Fahrenheit", Factor: 1, Offset: 0},
		{Symbol: "K", Name: "Kelvin", Factor: 1, Offset: 0},
		{Symbol: "R", Name: "Rankine", Factor: 1, Offset: 0},
	})

	// Volume — base: liter
	registerCategory("volume", "Volume and capacity measurements", "L", []Unit{
		{Symbol: "L", Name: "liter", Factor: 1},
		{Symbol: "mL", Name: "milliliter", Factor: 0.001},
		{Symbol: "m3", Name: "cubic meter", Factor: 1000},
		{Symbol: "cm3", Name: "cubic centimeter", Factor: 0.001},
		{Symbol: "mm3", Name: "cubic millimeter", Factor: 1e-6},
		{Symbol: "ft3", Name: "cubic foot", Factor: 28.316846592},
		{Symbol: "in3", Name: "cubic inch", Factor: 0.016387064},
		{Symbol: "gal_us", Name: "US gallon", Factor: 3.785411784},
		{Symbol: "gal_uk", Name: "UK gallon", Factor: 4.54609},
		{Symbol: "qt_us", Name: "US quart", Factor: 0.946352946},
		{Symbol: "pt_us", Name: "US pint", Factor: 0.473176473},
		{Symbol: "cup_us", Name: "US cup", Factor: 0.2365882365},
		{Symbol: "floz_us", Name: "US fluid ounce", Factor: 0.0295735295625},
		{Symbol: "tbsp_us", Name: "US tablespoon", Factor: 0.01478676478125},
		{Symbol: "tsp_us", Name: "US teaspoon", Factor: 0.00492892159375},
		{Symbol: "bbl_oil", Name: "oil barrel", Factor: 158.987294928},
	})

	// Area — base: square meter
	registerCategory("area", "Area measurements", "m2", []Unit{
		{Symbol: "m2", Name: "square meter", Factor: 1},
		{Symbol: "km2", Name: "square kilometer", Factor: 1e6},
		{Symbol: "cm2", Name: "square centimeter", Factor: 0.0001},
		{Symbol: "mm2", Name: "square millimeter", Factor: 1e-6},
		{Symbol: "ha", Name: "hectare", Factor: 10000},
		{Symbol: "ac", Name: "acre", Factor: 4046.8564224},
		{Symbol: "ft2", Name: "square foot", Factor: 0.09290304},
		{Symbol: "in2", Name: "square inch", Factor: 0.00064516},
		{Symbol: "yd2", Name: "square yard", Factor: 0.83612736},
		{Symbol: "mi2", Name: "square mile", Factor: 2589988.110336},
		{Symbol: "nmi2", Name: "square nautical mile", Factor: 3429904},
	})

	// Speed — base: meter per second
	registerCategory("speed", "Speed and velocity measurements", "m/s", []Unit{
		{Symbol: "m/s", Name: "meter per second", Factor: 1},
		{Symbol: "km/h", Name: "kilometer per hour", Factor: 0.2777777777777778},
		{Symbol: "mph", Name: "mile per hour", Factor: 0.44704},
		{Symbol: "ft/s", Name: "foot per second", Factor: 0.3048},
		{Symbol: "kn", Name: "knot", Factor: 0.5144444444444444},
		{Symbol: "mach", Name: "Mach (at sea level)", Factor: 343},
		{Symbol: "c", Name: "speed of light", Factor: 299792458},
	})

	// Time — base: second
	registerCategory("time", "Time duration measurements", "s", []Unit{
		{Symbol: "ns", Name: "nanosecond", Factor: 1e-9},
		{Symbol: "us", Name: "microsecond", Factor: 1e-6},
		{Symbol: "ms", Name: "millisecond", Factor: 0.001},
		{Symbol: "s", Name: "second", Factor: 1},
		{Symbol: "min", Name: "minute", Factor: 60},
		{Symbol: "h", Name: "hour", Factor: 3600},
		{Symbol: "d", Name: "day", Factor: 86400},
		{Symbol: "wk", Name: "week", Factor: 604800},
		{Symbol: "mo", Name: "month (30 days)", Factor: 2592000},
		{Symbol: "yr", Name: "year (365 days)", Factor: 31536000},
		{Symbol: "dec", Name: "decade", Factor: 315360000},
		{Symbol: "cent", Name: "century", Factor: 3153600000},
	})

	// Data — base: byte
	registerCategory("data", "Digital data and storage measurements", "B", []Unit{
		{Symbol: "bit", Name: "bit", Factor: 0.125},
		{Symbol: "B", Name: "byte", Factor: 1},
		{Symbol: "KB", Name: "kilobyte (decimal)", Factor: 1000},
		{Symbol: "MB", Name: "megabyte (decimal)", Factor: 1e6},
		{Symbol: "GB", Name: "gigabyte (decimal)", Factor: 1e9},
		{Symbol: "TB", Name: "terabyte (decimal)", Factor: 1e12},
		{Symbol: "PB", Name: "petabyte (decimal)", Factor: 1e15},
		{Symbol: "KiB", Name: "kibibyte (binary)", Factor: 1024},
		{Symbol: "MiB", Name: "mebibyte (binary)", Factor: 1048576},
		{Symbol: "GiB", Name: "gibibyte (binary)", Factor: 1073741824},
		{Symbol: "TiB", Name: "tebibyte (binary)", Factor: 1099511627776},
		{Symbol: "PiB", Name: "pebibyte (binary)", Factor: 1125899906842624},
		{Symbol: "Kb", Name: "kilobit (decimal)", Factor: 125},
		{Symbol: "Mb", Name: "megabit (decimal)", Factor: 125000},
		{Symbol: "Gb", Name: "gigabit (decimal)", Factor: 125000000},
	})

	// Pressure — base: pascal
	registerCategory("pressure", "Pressure measurements", "Pa", []Unit{
		{Symbol: "Pa", Name: "pascal", Factor: 1},
		{Symbol: "kPa", Name: "kilopascal", Factor: 1000},
		{Symbol: "MPa", Name: "megapascal", Factor: 1e6},
		{Symbol: "GPa", Name: "gigapascal", Factor: 1e9},
		{Symbol: "bar", Name: "bar", Factor: 100000},
		{Symbol: "mbar", Name: "millibar", Factor: 100},
		{Symbol: "atm", Name: "atmosphere", Factor: 101325},
		{Symbol: "psi", Name: "pound per square inch", Factor: 6894.757293168},
		{Symbol: "mmHg", Name: "millimeter of mercury", Factor: 133.322387415},
		{Symbol: "inHg", Name: "inch of mercury", Factor: 3386.389},
		{Symbol: "torr", Name: "torr", Factor: 133.32236842105263},
		{Symbol: "hPa", Name: "hectopascal", Factor: 100},
	})

	// Energy — base: joule
	registerCategory("energy", "Energy measurements", "J", []Unit{
		{Symbol: "J", Name: "joule", Factor: 1},
		{Symbol: "kJ", Name: "kilojoule", Factor: 1000},
		{Symbol: "MJ", Name: "megajoule", Factor: 1e6},
		{Symbol: "cal", Name: "calorie", Factor: 4.184},
		{Symbol: "kcal", Name: "kilocalorie", Factor: 4184},
		{Symbol: "Wh", Name: "watt-hour", Factor: 3600},
		{Symbol: "kWh", Name: "kilowatt-hour", Factor: 3600000},
		{Symbol: "MWh", Name: "megawatt-hour", Factor: 3600000000},
		{Symbol: "BTU", Name: "British thermal unit", Factor: 1055.05585262},
		{Symbol: "ftlb", Name: "foot-pound", Factor: 1.3558179483314004},
		{Symbol: "eV", Name: "electronvolt", Factor: 1.602176634e-19},
		{Symbol: "erg", Name: "erg", Factor: 1e-7},
		{Symbol: "therm", Name: "therm (US)", Factor: 105480400},
	})

	// Power — base: watt
	registerCategory("power", "Power measurements", "W", []Unit{
		{Symbol: "W", Name: "watt", Factor: 1},
		{Symbol: "kW", Name: "kilowatt", Factor: 1000},
		{Symbol: "MW", Name: "megawatt", Factor: 1e6},
		{Symbol: "GW", Name: "gigawatt", Factor: 1e9},
		{Symbol: "mW", Name: "milliwatt", Factor: 0.001},
		{Symbol: "hp", Name: "horsepower (mechanical)", Factor: 745.6998715822702},
		{Symbol: "hp_m", Name: "horsepower (metric)", Factor: 735.49875},
		{Symbol: "BTU/h", Name: "BTU per hour", Factor: 0.29307107},
		{Symbol: "ftlb/s", Name: "foot-pound per second", Factor: 1.3558179483314004},
		{Symbol: "cal/s", Name: "calorie per second", Factor: 4.184},
	})

	// Angle — base: degree
	registerCategory("angle", "Angle measurements", "deg", []Unit{
		{Symbol: "deg", Name: "degree", Factor: 1},
		{Symbol: "rad", Name: "radian", Factor: 57.29577951308232},
		{Symbol: "grad", Name: "gradian", Factor: 0.9},
		{Symbol: "arcmin", Name: "arcminute", Factor: 0.016666666666666666},
		{Symbol: "arcsec", Name: "arcsecond", Factor: 0.0002777777777777778},
		{Symbol: "turn", Name: "turn (full circle)", Factor: 360},
	})

	// Frequency — base: hertz
	registerCategory("frequency", "Frequency measurements", "Hz", []Unit{
		{Symbol: "Hz", Name: "hertz", Factor: 1},
		{Symbol: "kHz", Name: "kilohertz", Factor: 1000},
		{Symbol: "MHz", Name: "megahertz", Factor: 1e6},
		{Symbol: "GHz", Name: "gigahertz", Factor: 1e9},
		{Symbol: "THz", Name: "terahertz", Factor: 1e12},
		{Symbol: "rpm", Name: "revolutions per minute", Factor: 0.016666666666666666},
		{Symbol: "rad/s", Name: "radian per second", Factor: 0.15915494309189535},
	})
}

// LookupUnit finds a unit by symbol. Tries exact match first, then case-insensitive.
func LookupUnit(symbol string) (Unit, bool) {
	if u, ok := unitRegistry[symbol]; ok {
		return u, true
	}
	lower := strings.ToLower(symbol)
	matches, ok := unitLowerMap[lower]
	if !ok || len(matches) == 0 {
		return Unit{}, false
	}
	if len(matches) == 1 {
		return matches[0], true
	}
	return Unit{}, false
}

// LookupUnitWithHint finds a unit by symbol and returns a helpful hint on failure.
func LookupUnitWithHint(symbol string) (Unit, string, bool) {
	if u, ok := unitRegistry[symbol]; ok {
		return u, "", true
	}
	lower := strings.ToLower(symbol)
	matches, ok := unitLowerMap[lower]
	if !ok || len(matches) == 0 {
		return Unit{}, fmt.Sprintf("unknown unit: %s | hint: call GET /categories to list all categories, then GET /units/{category} to see available units", symbol), false
	}
	if len(matches) == 1 {
		return matches[0], "", true
	}
	symbols := make([]string, len(matches))
	for i, m := range matches {
		symbols[i] = m.Symbol
	}
	return Unit{}, fmt.Sprintf("ambiguous unit: %s matches %s | hint: use the exact symbol (case-sensitive): %s", symbol, strings.Join(symbols, ", "), strings.Join(symbols, " or ")), false
}

// LookupCategory finds a category by name (case-insensitive).
func LookupCategory(name string) (Category, bool) {
	c, ok := categoryRegistry[strings.ToLower(name)]
	return c, ok
}

// ListCategories returns all category names sorted alphabetically.
func ListCategories() []string {
	names := make([]string, 0, len(categoryRegistry))
	for name := range categoryRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListUnits returns all unit symbols for a category.
func ListUnits(category string) ([]string, bool) {
	c, ok := categoryRegistry[strings.ToLower(category)]
	if !ok {
		return nil, false
	}
	return c.Units, true
}

// AllCategories returns all category definitions sorted by name.
func AllCategories() []Category {
	cats := make([]Category, 0, len(categoryRegistry))
	for _, c := range categoryRegistry {
		cats = append(cats, c)
	}
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Name < cats[j].Name
	})
	return cats
}

// GetCategoryForUnit returns the category name for a given unit symbol.
func GetCategoryForUnit(symbol string) (string, bool) {
	cat, ok := unitToCategory[symbol]
	return cat, ok
}

// Convert performs a unit conversion from one unit to another.
func Convert(from, to string, value float64) (*ConversionResult, error) {
	fromUnit, hint, ok := LookupUnitWithHint(from)
	if !ok {
		return nil, fmt.Errorf("%s", hint)
	}

	toUnit, hint, ok := LookupUnitWithHint(to)
	if !ok {
		return nil, fmt.Errorf("%s", hint)
	}

	fromCat, _ := GetCategoryForUnit(fromUnit.Symbol)
	toCat, _ := GetCategoryForUnit(toUnit.Symbol)
	if fromCat != toCat {
		return nil, fmt.Errorf("cannot convert %s (%s) to %s (%s) — different categories | hint: both units must be in the same category. Call GET /units/%s to see valid units", fromUnit.Symbol, fromCat, toUnit.Symbol, toCat, fromCat)
	}

	var result float64

	if fromCat == "temperature" {
		result = convertTemperature(fromUnit.Symbol, toUnit.Symbol, value)
	} else {
		result = value * fromUnit.Factor / toUnit.Factor
	}

	result = roundToPrecision(result, 12)

	return &ConversionResult{
		From:     fromUnit.Symbol,
		To:       toUnit.Symbol,
		Value:    value,
		Result:   result,
		Category: fromCat,
	}, nil
}

// convertTemperature handles the non-linear temperature conversions.
func convertTemperature(from, to string, value float64) float64 {
	var celsius float64
	switch from {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	case "R":
		celsius = (value - 491.67) * 5 / 9
	default:
		celsius = value
	}

	switch to {
	case "C":
		return celsius
	case "F":
		return celsius*9/5 + 32
	case "K":
		return celsius + 273.15
	case "R":
		return (celsius + 273.15) * 9 / 5
	default:
		return celsius
	}
}

// roundToPrecision rounds to n significant decimal places.
func roundToPrecision(val float64, precision int) float64 {
	if val == 0 {
		return 0
	}
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}

// FormatValue formats a float64 nicely, removing trailing zeros.
func FormatValue(val float64) string {
	if val == math.Trunc(val) && math.Abs(val) < 1e15 {
		return fmt.Sprintf("%.0f", val)
	}
	s := fmt.Sprintf("%.10g", val)
	return s
}

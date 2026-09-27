package model

import (
	"math"
	"testing"
)

func TestLookupUnit(t *testing.T) {
	tests := []struct {
		symbol string
		want   bool
		name   string
	}{
		{"m", true, "meter"},
		{"M", true, "meter (case insensitive)"},
		{"km", true, "kilometer"},
		{"kg", true, "kilogram"},
		{"C", true, "Celsius"},
		{"c", true, "speed of light (exact match)"},
		{"Hz", true, "hertz"},
		{"hz", true, "hertz (case insensitive)"},
		{"foo", false, "unknown unit"},
		{"", false, "empty string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, ok := LookupUnit(tt.symbol)
			if ok != tt.want {
				t.Errorf("LookupUnit(%q) ok = %v, want %v", tt.symbol, ok, tt.want)
			}
			if ok && u.Symbol == "" {
				t.Errorf("LookupUnit(%q) returned empty symbol", tt.symbol)
			}
		})
	}
}

func TestLookupCategory(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"length", true},
		{"Length", true},
		{"mass", true},
		{"temperature", true},
		{"volume", true},
		{"area", true},
		{"speed", true},
		{"time", true},
		{"data", true},
		{"pressure", true},
		{"energy", true},
		{"power", true},
		{"angle", true},
		{"frequency", true},
		{"foo", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := LookupCategory(tt.name)
			if ok != tt.want {
				t.Errorf("LookupCategory(%q) ok = %v, want %v", tt.name, ok, tt.want)
			}
		})
	}
}

func TestListCategories(t *testing.T) {
	cats := ListCategories()
	if len(cats) != 13 {
		t.Errorf("ListCategories() returned %d categories, want 13", len(cats))
	}
	for i := 1; i < len(cats); i++ {
		if cats[i-1] > cats[i] {
			t.Errorf("ListCategories() not sorted: %q > %q", cats[i-1], cats[i])
		}
	}
}

func TestListUnits(t *testing.T) {
	units, ok := ListUnits("length")
	if !ok {
		t.Fatal("ListUnits(length) returned false")
	}
	if len(units) < 10 {
		t.Errorf("ListUnits(length) returned %d units, want at least 10", len(units))
	}

	_, ok = ListUnits("nonexistent")
	if ok {
		t.Error("ListUnits(nonexistent) should return false")
	}
}

func TestConvertLength(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"m", "km", 1000, 1, 1e-9},
		{"km", "m", 1, 1000, 1e-9},
		{"m", "ft", 1, 3.280839895013123, 1e-9},
		{"ft", "m", 1, 0.3048, 1e-9},
		{"mi", "km", 1, 1.609344, 1e-9},
		{"in", "cm", 1, 2.54, 1e-9},
		{"nmi", "km", 1, 1.852, 1e-9},
		{"m", "m", 5, 5, 1e-9},
		{"yd", "m", 1, 0.9144, 1e-9},
		{"ly", "km", 1, 9.4607304725808e12, 1e3},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v (diff: %v)", tt.from, tt.to, tt.value, result.Result, tt.want, math.Abs(result.Result-tt.want))
			}
		})
	}
}

func TestConvertMass(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"kg", "g", 1, 1000, 1e-9},
		{"g", "kg", 1000, 1, 1e-9},
		{"lb", "kg", 1, 0.45359237, 1e-9},
		{"kg", "lb", 1, 2.2046226218487757, 1e-9},
		{"oz", "g", 1, 28.349523125, 1e-9},
		{"t", "kg", 1, 1000, 1e-9},
		{"st", "lb", 1, 14, 1e-9},
		{"ct", "g", 1, 0.2, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertTemperature(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"C", "F", 0, 32, 1e-9},
		{"C", "F", 100, 212, 1e-9},
		{"F", "C", 32, 0, 1e-9},
		{"F", "C", 212, 100, 1e-9},
		{"C", "K", 0, 273.15, 1e-9},
		{"K", "C", 273.15, 0, 1e-9},
		{"C", "C", 25, 25, 1e-9},
		{"F", "K", 32, 273.15, 1e-9},
		{"K", "F", 300, 80.33, 0.01},
		{"C", "R", 0, 491.67, 1e-9},
		{"R", "C", 491.67, 0, 1e-9},
		{"F", "R", 212, 671.67, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertVolume(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"L", "mL", 1, 1000, 1e-9},
		{"mL", "L", 1000, 1, 1e-9},
		{"gal_us", "L", 1, 3.785411784, 1e-9},
		{"L", "gal_us", 3.785411784, 1, 1e-9},
		{"gal_uk", "L", 1, 4.54609, 1e-9},
		{"m3", "L", 1, 1000, 1e-9},
		{"cup_us", "mL", 1, 236.5882365, 1e-6},
		{"tbsp_us", "tsp_us", 1, 3, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertSpeed(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"km/h", "m/s", 36, 10, 1e-9},
		{"m/s", "km/h", 10, 36, 1e-9},
		{"mph", "km/h", 1, 1.609344, 1e-9},
		{"kn", "km/h", 1, 1.852, 1e-9},
		{"mph", "m/s", 1, 0.44704, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertData(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"KB", "B", 1, 1000, 1e-9},
		{"KiB", "B", 1, 1024, 1e-9},
		{"MB", "KB", 1, 1000, 1e-9},
		{"MiB", "KiB", 1, 1024, 1e-9},
		{"GB", "MB", 1, 1000, 1e-9},
		{"bit", "B", 8, 1, 1e-9},
		{"B", "bit", 1, 8, 1e-9},
		{"KB", "KiB", 1.024, 1, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertTime(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"min", "s", 1, 60, 1e-9},
		{"h", "min", 1, 60, 1e-9},
		{"d", "h", 1, 24, 1e-9},
		{"wk", "d", 1, 7, 1e-9},
		{"yr", "d", 1, 365, 1e-9},
		{"ms", "s", 1000, 1, 1e-9},
		{"s", "ms", 1, 1000, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertPressure(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"bar", "Pa", 1, 100000, 1e-3},
		{"atm", "Pa", 1, 101325, 1e-3},
		{"psi", "bar", 14.503773773, 1, 1e-3},
		{"kPa", "Pa", 1, 1000, 1e-9},
		{"mmHg", "Pa", 1, 133.322387415, 1e-6},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertEnergy(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"kJ", "J", 1, 1000, 1e-9},
		{"kcal", "cal", 1, 1000, 1e-9},
		{"kWh", "J", 1, 3600000, 1e-3},
		{"cal", "J", 1, 4.184, 1e-9},
		{"BTU", "J", 1, 1055.05585262, 1e-6},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertPower(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"kW", "W", 1, 1000, 1e-9},
		{"hp", "W", 1, 745.6998715822702, 1e-6},
		{"MW", "kW", 1, 1000, 1e-9},
		{"BTU/h", "W", 1, 0.29307107, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertAngle(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"rad", "deg", 1, 57.29577951308232, 1e-9},
		{"deg", "rad", 180, 3.141592653589793, 1e-9},
		{"turn", "deg", 1, 360, 1e-9},
		{"grad", "deg", 100, 90, 1e-9},
		{"arcmin", "deg", 60, 1, 1e-9},
		{"arcsec", "deg", 3600, 1, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertFrequency(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"kHz", "Hz", 1, 1000, 1e-9},
		{"MHz", "kHz", 1, 1000, 1e-9},
		{"GHz", "MHz", 1, 1000, 1e-9},
		{"rpm", "Hz", 60, 1, 1e-9},
		{"Hz", "rpm", 1, 60, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertArea(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		value   float64
		want    float64
		epsilon float64
	}{
		{"km2", "m2", 1, 1e6, 1e-3},
		{"ha", "m2", 1, 10000, 1e-9},
		{"ac", "m2", 1, 4046.8564224, 1e-6},
		{"ft2", "m2", 1, 0.09290304, 1e-9},
		{"mi2", "km2", 1, 2.589988110336, 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("Convert(%q, %q, %v) error: %v", tt.from, tt.to, tt.value, err)
			}
			if math.Abs(result.Result-tt.want) > tt.epsilon {
				t.Errorf("Convert(%q, %q, %v) = %v, want %v", tt.from, tt.to, tt.value, result.Result, tt.want)
			}
		})
	}
}

func TestConvertErrors(t *testing.T) {
	_, err := Convert("foo", "bar", 1)
	if err == nil {
		t.Error("Convert with unknown units should return error")
	}

	_, err = Convert("m", "kg", 1)
	if err == nil {
		t.Error("Convert across categories should return error")
	}
	if err != nil && !contains(err.Error(), "different categories") {
		t.Errorf("Cross-category error should mention 'different categories', got: %v", err)
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{1000, "1000"},
		{3.14, "3.14"},
		{-5, "-5"},
		{0.5, "0.5"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatValue(tt.val)
			if got != tt.want {
				t.Errorf("FormatValue(%v) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestGetCategoryForUnit(t *testing.T) {
	cat, ok := GetCategoryForUnit("m")
	if !ok || cat != "length" {
		t.Errorf("GetCategoryForUnit(m) = %q, %v, want length, true", cat, ok)
	}

	_, ok = GetCategoryForUnit("nonexistent")
	if ok {
		t.Error("GetCategoryForUnit(nonexistent) should return false")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package units

type Value struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

func ConvertBytes(unit Value, toUnit string) Value {
	var result = Value{Value: 0, Unit: "B"}

	switch unit.Unit {
	case "B":
		result.Value = unit.Value
	case "KB":
		result.Value = unit.Value * 1000
	case "MB":
		result.Value = unit.Value * (1000 * 1000)
	case "GB":
		result.Value = unit.Value * (1000 * 1000 * 1000)
	case "TB":
		result.Value = unit.Value * (1000 * 1000 * 1000 * 1000)
	case "PB":
		result.Value = unit.Value * (1000 * 1000 * 1000 * 1000 * 1000)
	}

	switch toUnit {
	case "B":
		result.Value = result.Value
	case "KB":
		result.Value = result.Value / 1000
	case "MB":
		result.Value = result.Value / (1000 * 1000)
	case "GB":
		result.Value = result.Value / (1000 * 1000 * 1000)
	case "TB":
		result.Value = result.Value / (1000 * 1000 * 1000 * 1000)
	case "PB":
		result.Value = result.Value / (1000 * 1000 * 1000 * 1000 * 1000)
	}

	return result
}

/*
	func ToB(unit Value) Value {
		return ConvertBytes(unit, "B")
	}

	func ToKB(unit Value) Value {
		return ConvertBytes(unit, "KB")
	}

	func ToMB(unit Value) Value {
		return ConvertBytes(unit, "MB")
	}
*/
func ToGB(unit Value) Value {
	return ConvertBytes(unit, "GB")
}

/*
func ToPB(unit Value) Value {
	return ConvertBytes(unit, "PB")
}
*/

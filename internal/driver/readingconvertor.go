package driver

// convertReadingValueType converts  int, uint, float to correct value type according to profileValueType
func convertReadingValueType(profileValueType string, readingValueType string, reading interface{}) interface{} {
	switch readingValueType {
	case "int":
		reading = convertInt(profileValueType, reading)
	case "uint":
		reading = convertUnit(profileValueType, reading)
	case "float32":
		reading = convertFloat32(profileValueType, reading)
	case "float64":
		reading = convertFloat64(profileValueType, reading)
	}
	return reading
}

func convertInt(profileValueType string, reading interface{}) interface{} {
	switch profileValueType {
	case "int8":
		reading = int8(reading.(int64))
	case "int16":
		reading = int16(reading.(int64))
	case "int32":
		reading = int32(reading.(int64))
	case "int64":
		reading = int64(reading.(int64))
	case "uint8":
		reading = uint8(reading.(int64))
	case "uint16":
		reading = uint16(reading.(int64))
	case "uint32":
		reading = uint32(reading.(int64))
	case "uint64":
		reading = uint64(reading.(int64))
	}
	return reading
}

func convertUnit(profileValueType string, reading interface{}) interface{} {
	switch profileValueType {
	case "int8":
		reading = int8(reading.(uint64))
	case "int16":
		reading = int16(reading.(uint64))
	case "int32":
		reading = int32(reading.(uint64))
	case "int64":
		reading = int64(reading.(uint64))
	case "uint8":
		reading = uint8(reading.(uint64))
	case "uint16":
		reading = uint16(reading.(uint64))
	case "uint32":
		reading = uint32(reading.(uint64))
	case "uint64":
		reading = uint64(reading.(uint64))
	}
	return reading
}

func convertFloat32(profileValueType string, reading interface{}) interface{} {
	switch profileValueType {
	case "int8":
		reading = int8(reading.(float32))
	case "int16":
		reading = int16(reading.(float32))
	case "int32":
		reading = int32(reading.(float32))
	case "int64":
		reading = int64(reading.(float32))
	case "uint8":
		reading = uint8(reading.(float32))
	case "uint16":
		reading = uint16(reading.(float32))
	case "uint32":
		reading = uint32(reading.(float32))
	case "uint64":
		reading = uint64(reading.(float32))
	case "float64":
		reading = float64(reading.(float32))
	}
	return reading
}

func convertFloat64(profileValueType string, reading interface{}) interface{} {
	switch profileValueType {
	case "int8":
		reading = int8(reading.(float64))
	case "int16":
		reading = int16(reading.(float64))
	case "int32":
		reading = int32(reading.(float64))
	case "int64":
		reading = int64(reading.(float64))
	case "uint8":
		reading = uint8(reading.(float64))
	case "uint16":
		reading = uint16(reading.(float64))
	case "uint32":
		reading = uint32(reading.(float64))
	case "uint64":
		reading = uint64(reading.(float64))
	case "float32":
		reading = float32(reading.(float64))
	}
	return reading
}

package sunspec

func toInt16(v interface{}) int16 {
	switch v := v.(type) {
	case int:
		return int16(v)
	case float64:
		return int16(v)
	case int16:
		return v
	}
	return 0
}

func toInt32(v interface{}) int32 {
	switch v := v.(type) {
	case int:
		return int32(v)
	case float64:
		return int32(v)
	case int32:
		return v
	}
	return 0
}

func toInt64(v interface{}) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case int64:
		return v
	}
	return 0
}

func toUint16(v interface{}) uint16 {
	switch v := v.(type) {
	case int:
		return uint16(v)
	case float64:
		return uint16(v)
	case uint16:
		return v
	}
	return 0
}

func toUint32(v interface{}) uint32 {
	switch v := v.(type) {
	case int:
		return uint32(v)
	case float64:
		return uint32(v)
	case uint32:
		return v
	}
	return 0
}

func toUint64(v interface{}) uint64 {
	switch v := v.(type) {
	case int:
		return uint64(v)
	case float64:
		return uint64(v)
	case uint64:
		return v
	}
	return 0
}

func toFloat32(v interface{}) float32 {
	switch v := v.(type) {
	case int:
		return float32(v)
	case float64:
		return float32(v)
	case float32:
		return v
	}
	return 0
}

func toFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	}
	return 0
}

func toByteS(v interface{}) []byte {
	switch v := v.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	}
	return nil
}

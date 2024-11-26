package foundation

func String(v string) *string {
	return &v
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func Uint(v uint) *uint {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

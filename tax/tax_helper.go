package tax

func RoundDay(time int64) int64 {
	rounding := (time + (7 * 86400)) % 86400
	return (time - rounding)
}

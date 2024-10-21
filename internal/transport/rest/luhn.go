package rest

const (
	asciiZero = 48
	asciiTen  = 57
)

func luhnValid(num string) bool {
	if len(num) == 0 {
		return false
	}
	var sum int64
	p := len(num) % 2
	for i, d := range num {
		if d < asciiZero || d > asciiTen {
			return false
		}
		d = d - asciiZero
		if i%2 == p {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += int64(d)
	}
	return sum%10 == 0
}

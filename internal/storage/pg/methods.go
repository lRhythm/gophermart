package pg

// Fox example: 1.2345 * mod → 12345; 12345 / mod = 1.2345.
const mod = 10000

// s2r - storage points (uint) to real points (float32).
func s2r(u uint) float32 {
	return float32(u) / mod
}

// r2s - real points (float32) to storage points (uint).
func r2s(f float32) uint {
	return uint(f * mod)
}

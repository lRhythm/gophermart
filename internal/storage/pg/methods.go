package pg

// s2r - storage points (uint) to real points (float32).
func s2r(u uint) float32 {
	return float32(u) / 100
}

// r2s - real points (float32) to storage points (uint).
func r2s(f float32) uint {
	return uint(f * 100)
}

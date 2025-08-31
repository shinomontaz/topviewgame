package main

func computeMask(x, y int, tiles []*MapTile, gd GameData) uint8 {
	// Use this order: N, S, E, W, NW, NE, SW, SE
	neighbors := [8][2]int{
		{0, -1},  // N
		{0, 1},   // S
		{1, 0},   // E
		{-1, 0},  // W
		{-1, -1}, // NW
		{1, -1},  // NE
		{-1, 1},  // SW
		{1, 1},   // SE
	}

	var mask uint8 = 0
	for bit, delta := range neighbors {
		nx, ny := x+delta[0], y+delta[1]
		wall := false
		if nx < 0 || nx >= gd.MapWidth || ny < 0 || ny >= gd.MapHeight {
			wall = true
		} else {
			nIdx := gd.GetIndexFromXY(nx, ny)
			wall = tiles[nIdx].Blocked
		}
		if wall {
			mask |= 1 << bit
		}
	}
	return mask
}

// blobMaskToTile returns a tile name for a given 8-bit neighbor mask.
// Mask convention: 0 = wall, 1 = free
func blobMaskToTile(mask uint8) string {
	N := (mask >> 0) & 1
	S := (mask >> 1) & 1
	E := (mask >> 2) & 1
	W := (mask >> 3) & 1
	NW := (mask >> 4) & 1
	NE := (mask >> 5) & 1
	SW := (mask >> 6) & 1
	SE := (mask >> 7) & 1

	// 1. Full/none
	if N == 1 && S == 1 && E == 1 && W == 1 && NW == 1 && NE == 1 && SW == 1 && SE == 1 {
		return "full"
	}
	if N == 0 && S == 0 && E == 0 && W == 0 && NW == 0 && NE == 0 && SW == 0 && SE == 0 {
		return "none"
	}

	// 2. Cross (all four directions closed)
	// Cross: All four cardinal directions are occupied (1), at least one diagonal is free (0)
	if N == 1 && S == 1 && E == 1 && W == 1 {
		if NW == 0 && NE == 0 && SW == 0 && SE == 0 {
			return "cross"
		}
		// 1 diagonal free
		if NW == 0 && NE == 1 && SW == 1 && SE == 1 {
			return "cross_NE_SW_SE"
		}
		if NW == 1 && NE == 0 && SW == 1 && SE == 1 {
			return "cross_NW_SW_SE"
		}
		if NW == 1 && NE == 1 && SW == 0 && SE == 1 {
			return "cross_NW_NE_SE"
		}
		if NW == 1 && NE == 1 && SW == 1 && SE == 0 {
			return "cross_NW_NE_SW"
		}
		// 2 diagonals free
		if NW == 0 && NE == 0 && SW == 1 && SE == 1 {
			return "cross_SW_SE"
		}
		if NW == 0 && NE == 1 && SW == 0 && SE == 1 {
			return "cross_NE_SE"
		}
		if NW == 0 && NE == 1 && SW == 1 && SE == 0 {
			return "cross_NE_SW"
		}
		if NW == 1 && NE == 0 && SW == 0 && SE == 1 {
			return "cross_NE_SW"
		}
		if NW == 1 && NE == 0 && SW == 1 && SE == 0 {
			return "cross_NW_SW"
		}
		if NW == 1 && NE == 1 && SW == 0 && SE == 0 {
			return "cross_NW_NE"
		}
		// 3 diagonals free
		// TODO

	}

	// 3. T-junctions
	if N == 0 && S == 1 && E == 1 && W == 1 {
		if SW == 1 && SE == 1 {
			return "T_up_SW_SE"
		}
		if SW == 1 && SE == 0 {
			return "T_up_SE"
		}
		if SW == 0 && SE == 1 {
			return "T_up_SW"
		}
		if SW == 0 && SE == 0 {
			return "T_up"
		}
	}

	if S == 0 && N == 1 && E == 1 && W == 1 {
		if NW == 1 && NE == 1 {
			return "T_down_NW_NE"
		}
		if NW == 1 && NE == 0 {
			return "T_down_NE"
		}
		if NW == 0 && NE == 1 {
			return "T_down_NW"
		}
		if NW == 0 && NE == 0 {
			return "T_down"
		}
	}
	if E == 0 && N == 1 && S == 1 && W == 1 {
		if NW == 1 && SW == 1 {
			return "T_left_NW_SW"
		}
		if NW == 1 && SW == 0 {
			return "T_left_SW"
		}
		if NW == 0 && SW == 1 {
			return "T_left_NW"
		}
		if NW == 0 && SW == 0 {
			return "T_left"
		}

	}
	if W == 0 && N == 1 && S == 1 && E == 1 {
		if NE == 1 && SE == 1 {
			return "T_right_NE_SE"
		}
		if NE == 0 && SE == 1 {
			return "T_right_SE"
		}

		if NE == 1 && SE == 0 {
			return "T_right_NE"
		}

		if NE == 0 && SE == 0 {
			return "T_right"
		}
	}

	// 4. Straights (corridors)
	if N == 0 && S == 0 && E == 1 && W == 1 {
		return "straight_EW"
	}
	if E == 0 && W == 0 && N == 1 && S == 1 {
		return "straight_NS"
	}

	// 5. corners (classic border corners and inner corders)
	if N == 0 && W == 0 && S == 1 && E == 1 {
		if SE == 1 {
			return "corner_SE"
		}
		return "inner_corner_SE"
	}
	if N == 0 && E == 0 && S == 1 && W == 1 {
		if SW == 1 {
			return "corner_SW"
		}
		return "inner_corner_SW"
	}
	if S == 0 && W == 0 && N == 1 && E == 1 {
		if NE == 1 {
			return "corner_NE"
		}
		return "inner_corner_NE"
	}
	if S == 0 && E == 0 && N == 1 && W == 1 {
		if NW == 1 {
			return "corner_NW"
		}
		return "inner_corner_NW"
	}

	// 6. Edges (classic border edges)
	if N == 0 && S == 1 && E == 1 && W == 1 {
		return "edge_S"
	}
	if S == 0 && N == 1 && E == 1 && W == 1 {
		return "edge_N"
	}
	if W == 0 && N == 1 && S == 1 && E == 1 {
		return "edge_E"
	}
	if E == 0 && N == 1 && S == 1 && W == 1 {
		return "edge_W"
	}

	// 10. Fallback
	return "none"
}

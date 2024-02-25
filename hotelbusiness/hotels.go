//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	maxOut := guests[0].CheckOutDate
	for _, g := range guests {
		maxOut = max(maxOut, g.CheckOutDate)
	}

	table := make([]int, maxOut+1)
	for _, g := range guests {
		table[g.CheckInDate] += 1
		table[g.CheckOutDate] -= 1
	}

	result := make([]Load, 0)
	count := 0

	for date, delta := range table {
		if delta == 0 {
			continue
		}

		count += delta
		result = append(result, Load{StartDate: date, GuestCount: count})
	}

	return result
}

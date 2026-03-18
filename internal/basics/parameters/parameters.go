package parameters

type Currency struct {
	Code  string
	Money float64
	Fx    int64
}

func valueParameter(c Currency, n int) []Currency {
	var arr []Currency
	for i := 0; i < n; i++ {
		arr = append(arr, c)
	}
	return arr
}

func pointerParameter(c Currency, n int) []*Currency {
	var arr []*Currency
	for i := 0; i < n; i++ {
		arr = append(arr, &c)
	}
	return arr
}

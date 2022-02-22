package shared

const PerPage = 2

type Pagination struct {
	Page int
}

// Page returns page start and end index
// given a length and page index
// it returns 0,0 [] if length or page is 0
func Page(length, page uint) (start, end uint) {
	if length == 0 {
		return
	}

	start = page * PerPage
	if start > length {
		start = length
	}
	end = start + PerPage
	if end > length {
		end = length
	}
	return
}

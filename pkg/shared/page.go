package shared

const PerPage = 10

type Pagination struct {
	Page int
}

// Page returns page start and end index
// given a length and page index
// it returns 0,0 [] if length or page is 0
func Page(length, page, limit uint) (start, end uint) {
	if length == 0 {
		return
	}

	if limit == 0 {
		limit = PerPage
	} else if limit > length {
		limit = length
	}

	start = page * limit
	if start > length {
		start = length
	}
	end = start + limit
	if end > length {
		end = length
	}
	return
}

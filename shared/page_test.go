package shared

import (
	"testing"
)

func TestPage(t *testing.T) {

	testCases := []struct {
		len   uint
		page  uint
		start uint
		end   uint
	}{
		{0, 1, 0, 0},
		{15, 0, 0, 0},
		{5, 1, 0, 5},
		{5, 2, 5, 5},
		{15, 1, 0, 10},
		{15, 2, 10, 15},
		{35, 3, 20, 30},
		{7, 1, 0, 7},
		{10, 10, 10, 10},
	}

	for _, testCase := range testCases {
		expectedStart, expectedEnd := testCase.start, testCase.end
		gotStart, gotEnd := Page(testCase.len, testCase.page)

		if expectedStart != gotStart {
			t.Errorf("expected start to be %d, got %d", expectedStart, gotStart)
		}
		if expectedEnd != gotEnd {
			t.Errorf("expected end to be %d, got %d", expectedEnd, gotEnd)
		}
	}
}

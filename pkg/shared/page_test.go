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
		{15, 0, 0, 10},
		{5, 1, 5, 5},
		{5, 2, 5, 5},
		{15, 1, 10, 15},
		{15, 2, 15, 15},
		{35, 0, 0, 10},
		{35, 1, 10, 20},
		{35, 2, 20, 30},
		{35, 3, 30, 35},
		{35, 4, 35, 35},
		{7, 0, 0, 7},
		{10, 10, 10, 10},
	}

	for _, testCase := range testCases {
		expectedStart, expectedEnd := testCase.start, testCase.end
		gotStart, gotEnd := Page(testCase.len, testCase.page, 0)

		if expectedStart != gotStart {
			t.Errorf("expected start to be %d, got %d", expectedStart, gotStart)
		}
		if expectedEnd != gotEnd {
			t.Errorf("expected end to be %d, got %d", expectedEnd, gotEnd)
		}
	}
}

package postgres

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(` LIMIT %d OFFSET %d `, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(` LIMIT %d `, limit)
	} else if offset > 0 {
		return fmt.Sprintf(` OFFSET %d `, offset)
	}
	return ""
}

func formatOrderBy(sortBy string, desc bool, defaultSort string, safeList []string) string {

	if !inSlice(sortBy, safeList) {
		sortBy = defaultSort
	}

	order := " ORDER BY " + sortBy
	if desc {
		order += " DESC "
	} else {
		order += " ASC "
	}
	return order
}

func inSlice(val string, s []string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

func equalMapStringString(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}

func areSlicesEqual(slice1, slice2 []string) bool {
	// Check if slices have the same length
	if len(slice1) != len(slice2) {
		return false
	}

	// Sort the slices
	sort.Strings(slice1)
	sort.Strings(slice2)

	// Compare each element
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func notSameDay(t1 time.Time, t2 time.Time) bool {
	return (t1.Year() != t2.Year()) || (t1.YearDay() != t2.YearDay())
}

func rowIndex(prefix byte, v string) int {
	v = strings.TrimSpace(v)
	if len(v) < 2 || v[0] != prefix {
		return 0
	}
	if idx, err := strconv.Atoi(v[1:]); err == nil {
		return idx
	}
	return 0
}

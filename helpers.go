package main

import "strconv"

func portLabel(i int) string {
	return strconv.Itoa(i + 1)
}

// bitGet64/bitGet32 read bit i (0-based, LSB first) out of a bit-packed
// port mask. ok is false when mask is nil (field not supported on this
// board).
func bitGet64(mask *int64, i int) (val bool, ok bool) {
	if mask == nil {
		return false, false
	}
	return (*mask>>uint(i))&1 != 0, true
}

func bitGet32(mask *int32, i int) (val bool, ok bool) {
	if mask == nil {
		return false, false
	}
	return (*mask>>uint(i))&1 != 0, true
}

// arrayLenInt32/arrayLenString return the length of the longest of the
// given slice pointers, treating nil as length 0.
func arrayLenInt32(ps ...*[]int32) int {
	n := 0
	for _, p := range ps {
		if p != nil && len(*p) > n {
			n = len(*p)
		}
	}
	return n
}

func arrayLenString(ps ...*[]string) int {
	n := 0
	for _, p := range ps {
		if p != nil && len(*p) > n {
			n = len(*p)
		}
	}
	return n
}

func atInt32(p *[]int32, i int) *int32 {
	if p == nil || i < 0 || i >= len(*p) {
		return nil
	}
	v := (*p)[i]
	return &v
}

func atInt64(p *[]int64, i int) *int64 {
	if p == nil || i < 0 || i >= len(*p) {
		return nil
	}
	v := (*p)[i]
	return &v
}

func atString(p *[]string, i int) *string {
	if p == nil || i < 0 || i >= len(*p) {
		return nil
	}
	v := (*p)[i]
	return &v
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

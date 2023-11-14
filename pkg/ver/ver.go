package ver

import (
	"fmt"
	"regexp"
	"strconv"
)

var rex = regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)$`)

type Ver struct {
	Major int64
	Minor int64
	Patch int64
}

func (v Ver) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func Parse(v string) (*Ver, error) {
	data := rex.FindStringSubmatch(v)
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid format: %s", v)
	}
	result := &Ver{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}
	var err error
	if result.Major, err = strconv.ParseInt(data[1], 10, 64); err != nil {
		return nil, err
	}
	if result.Minor, err = strconv.ParseInt(data[2], 10, 64); err != nil {
		return nil, err
	}
	if result.Patch, err = strconv.ParseInt(data[3], 10, 64); err != nil {
		return nil, err
	}
	return result, nil
}

func Compare(v1, v2 string) int {
	a, e1 := Parse(v1)
	b, e2 := Parse(v2)
	switch true {
	case e1 != nil && e2 != nil:
		return 0
	case e1 != nil:
		return -1
	case e2 != nil:
		return +1
	case a.Major < b.Major:
		return -1
	case a.Major > b.Major:
		return +1
	case a.Minor < b.Minor:
		return -1
	case a.Minor > b.Minor:
		return +1
	case a.Patch < b.Patch:
		return -1
	case a.Patch > b.Patch:
		return +1
	default:
		return 0
	}
}

func Max(vers ...string) *Ver {
	result := "v0.0.0"
	for _, ver := range vers {
		if Compare(result, ver) < 0 {
			result = ver
		}
	}
	v, err := Parse(result)
	if err != nil {
		return &Ver{
			Major: 0,
			Minor: 0,
			Patch: 0,
		}
	}
	return v
}

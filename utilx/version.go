package utilx

import (
	"fmt"
)

// func main() {

// 	version1 := "2.0.12.1.aa"
// 	version2 := "2.0.20.2,bb"

// 	fmt.Println(VersionCompare(version1, version2))
// }

func VersionCompare(version1, version2 string) int {

	vs1, vs2 := versionCompareParse(version1), versionCompareParse(version2)

	fmt.Println(version1, vs1)
	fmt.Println(version2, vs2)

	length := len(vs1)
	if len(vs1) > len(vs2) {
		length = len(vs2)
	}

	for i := 0; i < length; i++ {
		if lg := vs1[i] - vs2[i]; lg > 0 {
			return 1
		} else if lg < 0 {
			return -1
		}
	}

	if lg := len(vs1) - len(vs2); lg > 0 {
		return 1
	} else if lg < 0 {
		return -1
	}

	return 0
}

func versionCompareParse(version string) []int32 {

	ves := []int32{}

	var num int32 = -1
	for _, char := range version {

		if char >= '0' && char <= '9' {

			if num > -1 {
				num *= 10
			} else {
				num = 0
			}

			num += char - '0'

		} else {

			if num > -1 {
				ves = append(ves, num)
			}

			if char >= 'a' && char <= 'z' {
				ves = append(ves, char-'a'+10)
			}

			num = -1
		}
	}

	if num > -1 {
		ves = append(ves, num)
	}

	return ves
}

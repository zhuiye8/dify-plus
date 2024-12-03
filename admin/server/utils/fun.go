package utils

import "math"

// InArray @author: [Fantasia](https://www.npc0.com)
// @function: InArray
// @description: 判断是否在数组中
// @return: err error, conf config.Server
func InArray(value interface{}, array []interface{}) (isIn bool) {
	// 判断array是否数组
	for _, item := range array {
		if value == item {
			isIn = true
			return
		}
	}
	return false
}

// InUintArray @author: [Fantasia](https://www.npc0.com)
// @function: InUintArray
// @description: 判断是否在uint数组中
// @return: err error, conf config.Server
func InUintArray(value uint, array []uint) (isIn bool) {
	// 判断array是否数组
	for _, item := range array {
		if value == item {
			isIn = true
			return
		}
	}
	return false
}

// InStringArray @author: [Fantasia](https://www.npc0.com)
// @function: InStringArray
// @description: 判断是否在字符串数组中
// @return: err error, conf config.Server
func InStringArray(value string, array []string) (isIn bool) {
	// 判断array是否数组
	for _, item := range array {
		if value == item {
			isIn = true
			return
		}
	}
	return false
}

// AddAsteriskToString @author: [Fantasia](https://www.npc0.com)
// @function: AddAsteriskToString
// @description: 字符串加星号
// @return: err error, conf config.Server
func AddAsteriskToString(s string) string {
	// 计算要插入的星号数量
	numStars := int(math.Ceil(float64(len(s)) / 5))
	stars := ""
	for i := 0; i < numStars; i++ {
		stars += "*"
	}

	// 计算插入位置
	insertPos := len(s) / 2

	// 插入星号
	result := s[:insertPos] + stars + s[insertPos:]
	return result
}

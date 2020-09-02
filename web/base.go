package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InitServer() (*gin.Engine, error) {
	return initRouters()
}

func getStringArrayFromQuery(name, sep string, r *http.Request) (arr []string) {
	val := strings.TrimSpace(r.FormValue(name))
	if len(val) == 0 {
		return
	}

	return strings.Split(val, sep)
}

func getPage(page string) int {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		p = 1
	}

	return p
}

func getPageSize(ps string) int {
	p, err := strconv.Atoi(ps)
	if err != nil || p < 1 {
		p = 50
	} else if p > 200 {
		p = 200
	}
	return p
}

func getTime(t string) time.Time {
	t = strings.TrimSpace(t)
	time, _ := time.ParseInLocation("2006-01-02", t, time.Local)
	return time
}

// K字符串是否在SS数组中
func InStringArray(k string, ss []string) bool {
	for i := range ss {
		if ss[i] == k {
			return true
		}
	}

	return false
}

// 去除重复的元素
func UniqueStringArray(a []string) []string {
	al := len(a)
	if al == 0 {
		return a
	}

	ret := make([]string, al)
	index := 0

loopa:
	for i := 0; i < al; i++ {
		for j := 0; j < index; j++ {
			if a[i] == ret[j] {
				continue loopa
			}
		}
		ret[index] = a[i]
		index++
	}

	return ret[:index]
}

// 返回存在于 a 且不存在于 b 中的元素集合
func SubtractStringArray(a, b []string) (c []string) {
	c = []string{}

	for _, _a := range a {
		if !InStringArray(_a, b) {
			c = append(c, _a)
		}
	}

	return
}

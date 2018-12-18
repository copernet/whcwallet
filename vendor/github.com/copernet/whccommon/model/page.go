package model

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE_SIZE = 50
)

func DefaultInt(key string, defaultInt int, c *gin.Context) int {
	value, ok := c.GetQuery(key)
	if !ok {
		if value, ok = c.GetPostForm(key); !ok {
			value = c.Param(key)
		}
	}
	ret, err := strconv.Atoi(value)
	if err != nil {
		return defaultInt
	}
	return ret
}

func DefaultInt64(key string, defaultInt int64, c *gin.Context) int64 {
	value, ok := c.GetQuery(key)
	if !ok {
		if value, ok = c.GetPostForm(key); !ok {
			value = c.Param(key)
		}
	}
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultInt
	}
	return ret
}

func GetPageInfo(c *gin.Context) *PageInfo {
	var info PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		return &PageInfo{
			Page:     1,
			PageSize: DEFAULT_PAGE_SIZE,
			Total:    0,
		}
	} else {
		if info.Page == 0 {
			info.Page = 1
		}
		if info.PageSize == 0 {
			info.PageSize = DEFAULT_PAGE_SIZE
		}
		return &info
	}
}

func PageLimit(pageInfo *PageInfo) (int, int) {
	if pageInfo == nil {
		return 0, DEFAULT_PAGE_SIZE
	}
	return (pageInfo.Page - 1) * pageInfo.PageSize, pageInfo.PageSize
}

func EmptyPageInfo(info *PageInfo) *PageInfo {
	info.Total = 0
	info.List = []struct{}{}

	return info
}

type PageInfo struct {
	Page     int         `form:"pageNo" json:"pageNo"`
	PageSize int         `form:"pageSize"  json:"pageSize"`
	Total    int         `form:"total"  json:"total"`
	List     interface{} `json:"list"`
}

// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package dbutil

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/gorm"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/stringutil"
)

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}
type RequestWithSortKey interface {
	Request
	GetSortKey() string
}
type RequestWithReverse interface {
	RequestWithSortKey
	GetReverse() bool
}

const (
	TagName              = "json"
	SearchWordColumnName = "search_word"
)

func getReqValue(param interface{}) interface{} {
	switch value := param.(type) {
	case string:
		if value == "" {
			return nil
		}
		return []string{value}
	case *wrappers.StringValue:
		if value == nil {
			return nil
		}
		return []string{value.GetValue()}
	case *wrappers.Int32Value:
		if value == nil {
			return nil
		}
		return []int32{value.GetValue()}
	case []string:
		var values []string
		for _, v := range value {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) == 0 {
			return nil
		}
		return values
	}
	return nil
}

func GetDisplayColumns(displayColumns []string, wholeColumns []string) []string {
	if displayColumns == nil {
		return wholeColumns
	} else if len(displayColumns) == 0 {
		return nil
	} else {
		var newDisplayColumns []string
		for _, column := range displayColumns {
			if stringutil.Contains(wholeColumns, column) {
				newDisplayColumns = append(newDisplayColumns, column)
			}
		}
		return newDisplayColumns
	}
}

func GetFieldName(field *structs.Field) string {
	tag := field.Tag(TagName)
	t := strings.Split(tag, ",")
	if len(t) == 0 {
		return "-"
	}
	return t[0]
}

type Chain struct {
	*gorm.DB
}

func GetChain(tx *gorm.DB) *Chain {
	return &Chain{
		tx,
	}
}

func (c *Chain) BuildFilterConditions(req Request, tableName string, exclude ...string) *Chain {
	return c.buildFilterConditions(req, tableName, exclude...)
}

func (c *Chain) getSearchFilter(tableName string, value interface{}, exclude ...string) {
	var conditions []string
	if vs, ok := value.([]string); ok {
		for _, v := range vs {
			for _, column := range models.SearchColumns[tableName] {
				if stringutil.Contains(exclude, column) {
					continue
				}
				// if column suffix is _id, must exact match
				if strings.HasSuffix(column, "_id") {
					conditions = append(conditions, column+" = '"+v+"'")
				} else {
					likeV := "%" + stringutil.SimplifyString(v) + "%"
					conditions = append(conditions, column+" LIKE '"+likeV+"'")
				}
			}
		}
	} else if value != nil {
		logger.Warnf(nil, "search_word [%+v] is not string", value)
	}
	condition := strings.Join(conditions, " OR ")
	c.DB = c.DB.Where(condition)
}

func (c *Chain) buildFilterConditions(req Request, tableName string, exclude ...string) *Chain {
	for _, field := range structs.Fields(req) {
		column := GetFieldName(field)
		param := field.Value()
		indexedColumns, ok := models.IndexedColumns[tableName]
		if ok && stringutil.Contains(indexedColumns, column) {
			value := getReqValue(param)
			if value != nil {
				key := column
				c.DB = c.Where(key+" in (?)", value)
			}
		}
		if column == SearchWordColumnName && stringutil.Contains(models.SearchWordColumnTable, tableName) {
			value := getReqValue(param)
			c.getSearchFilter(tableName, value, exclude...)
		}
	}
	return c
}

func (c *Chain) AddQueryOrderDir(req Request, defaultColumn string) *Chain {
	order := "DESC"
	if r, ok := req.(RequestWithReverse); ok {
		if r.GetReverse() {
			order = "ASC"
		}
	}
	if r, ok := req.(RequestWithSortKey); ok {
		s := r.GetSortKey()
		if s != "" {
			defaultColumn = s
		}
	}
	c.DB = c.Order(defaultColumn + " " + order)
	return c
}

package gaia

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

var SyncDatabaseLock bool

// GetDatabaseTableColumns
// @Tags Database
// @Summary 获取数据库表结构
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) GetDatabaseTableColumns(table string) (nameList []request.DatabaseTableColumn, err error) {
	var rows *sql.Rows
	if rows, err = global.GVA_DB.Raw("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE "+
		"table_schema=? AND table_name=?", request.PostgreSQLDefaultSchema, table).Rows(); err != nil {
		return nameList, errors.New("table columns error:" + err.Error())
	}
	// 读取column_name
	for rows.Next() {
		var columnName, dataType, isNullable sql.NullString
		if err = rows.Scan(&columnName, &dataType, &isNullable); err == nil {
			var nullable = false
			if isNullable.String == "YES" {
				nullable = true
			}
			nameList = append(nameList, request.DatabaseTableColumn{
				ColumnName: columnName.String,
				DataType:   dataType.String,
				IsNullable: nullable,
			})
		}
	}
	// return
	return nameList, nil
}

// ForeachInstall
// @Tags Database
// @Summary 循环同步数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) ForeachInstall(
	logTable, newTable, keyName, orderName string, logColumnList,
	queryColumnList []string, addColumn map[string]interface{}, groupName, groupValue string) bool {
	// init
	var err error
	var where string
	// 查询新表最旧的数据到了哪
	var rows *sql.Rows
	var orderNumber sql.NullString
	orderText := fmt.Sprintf("\"%s\",\"%s\"", groupName, orderName)
	if err = global.GVA_DB.Raw(fmt.Sprintf("SELECT \"%s\" FROM \"%s\" WHERE \"%s\"='%s' ORDER BY %s ASC LIMIT 1",
		orderName, newTable, groupName, groupValue, orderText)).Row().Scan(&orderNumber); err != nil {
		// 查询旧表数据
		if len(groupName) > 0 {
			where = fmt.Sprintf("WHERE \"%s\"='%s'", groupName, groupValue)
		}
	} else {
		if len(groupName) > 0 {
			where = fmt.Sprintf("WHERE \"%s\"='%s' AND \"%s\"<='%s'", groupName, groupValue, orderName, orderNumber.String)
		} else {
			where = fmt.Sprintf("WHERE \"%s\"<='%s'", orderName, orderNumber.String)
		}
	}
	// 查询旧表数据
	global.GVA_LOG.Info("ForeachInstall:" + groupName + " " + groupValue + " " + orderNumber.String)
	if rows, err = global.GVA_DB.Raw(fmt.Sprintf("SELECT %s FROM \"%s\" %s ORDER BY %s DESC LIMIT %d",
		strings.Join(queryColumnList, ","), logTable, where, orderText, request.PostgreSQLDataLimit)).Rows(); err != nil {
		return true
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		// 创建一个 map 来保存每一条记录
		record := make(map[string]interface{})
		vals := make([]interface{}, len(logColumnList))
		for i := range vals {
			vals[i] = new(interface{})
		}
		if err = rows.Scan(vals...); err != nil {
			global.GVA_LOG.Fatal("failed to scan row: %v" + err.Error())
		}
		for i, col := range logColumnList {
			record[col] = *(vals[i].(*interface{}))
		}
		// 不跳过执行
		results = append(results, record)
	}

	// 3. 构建并执行 INSERT INTO 语句
	var values []interface{}
	var sqlStatement []string
	for _, result := range results {
		placeholders := make([]string, 0, len(logColumnList))
		for _, col := range logColumnList {
			placeholders = append(placeholders, "?")
			values = append(values, result[col])
		}
		// 新增增加的key
		for _, value := range addColumn {
			placeholders = append(placeholders, "?")
			values = append(values, value)
		}
		// 整合成 (a, b, add)
		sqlStatement = append(sqlStatement, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
	}
	// logColumnList 追加新增键
	for key, _ := range addColumn {
		logColumnList = append(logColumnList, key)
	}
	// sqlStatement length
	if len(sqlStatement) > 0 {
		// 执行插入
		if result := global.GVA_DB.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON CONFLICT (%s) DO NOTHING;",
			newTable, strings.Join(logColumnList, ", "), strings.Join(sqlStatement, ", "), keyName), values...); result.Error != nil {
			global.GVA_LOG.Debug("insert error: %v" + result.Error.Error())
			return true
		}
	}
	// return
	if len(results) == request.PostgreSQLDataLimit {
		return true
	}
	return false
}

// SyncDatabaseTableData
// @Tags Test
// @Summary 同步数据库表数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) SyncDatabaseTableData(logTable, newTable, keyName, orderName, groupName string) {
	// 查询旧表
	var err error
	if SyncDatabaseLock {
		return
	}
	SyncDatabaseLock = true
	var logColumns, newColumns []request.DatabaseTableColumn
	if logColumns, err = e.GetDatabaseTableColumns(logTable); err != nil {
		global.GVA_LOG.Debug("log " + err.Error())
		SyncDatabaseLock = false
		return
	}
	// new columns
	if newColumns, err = e.GetDatabaseTableColumns(newTable); err != nil {
		global.GVA_LOG.Debug("log " + err.Error())
		SyncDatabaseLock = false
		return
	}
	// 储存旧列表
	var logColumnList, queryColumnList []string
	var addColumn = make(map[string]interface{})
	var newTime = time.Now().Format("2006-01-02 15:04:05")
	for _, v := range logColumns {
		logColumnList = append(logColumnList, v.ColumnName)
		queryColumnList = append(queryColumnList, fmt.Sprintf("\"%s\"", v.ColumnName))
	}
	// 判断新增了什么字段
	for _, v := range newColumns {
		if utils.InStringArray(v.ColumnName, logColumnList) {
			// 非新增数据
			continue
		}
		// 新增的字段
		if v.IsNullable {
			// 默认为空
			addColumn[v.ColumnName] = nil
		} else {
			// 新增的字段
			switch v.DataType {
			case request.PostgreSQLDataTypeUUID:
				addColumn[v.ColumnName] = uuid.New().String()
				break
			case request.PostgreSQLDataTypeCharacterVarying, request.PostgreSQLDataTypeText:
				addColumn[v.ColumnName] = ""
				break
			case request.PostgreSQLDataTypeJSON:
				addColumn[v.ColumnName] = "{}"
				break
			case request.PostgreSQLDataTypeInteger, request.PostgreSQLDataTypeDoublePrecision:
				addColumn[v.ColumnName] = 0
				break
			case request.PostgreSQLDataTypeNumeric:
				addColumn[v.ColumnName] = 0.01
				break
			case request.PostgreSQLDataTypeTimestampWithoutTZ:
				addColumn[v.ColumnName] = newTime
				break
			case request.PostgreSQLDataTypeBoolean:
				addColumn[v.ColumnName] = false
				break
			}
		}
	}
	// 判断是否有groupName
	if len(groupName) > 0 {
		var rows *sql.Rows
		if rows, err = global.GVA_DB.Raw(fmt.Sprintf("SELECT %s FROM %s GROUP BY %s", groupName, logTable, groupName)).Rows(); err != nil {
			global.GVA_LOG.Debug("log SELECT" + err.Error())
			SyncDatabaseLock = false
			return
		}
		// 循环app列表
		var i = 0
		for rows.Next() {
			// 提取所有关联 group 列
			var groupValue string
			var group interface{}
			if err = rows.Scan(&group); err != nil {
				global.GVA_LOG.Debug("log Scan" + err.Error())
				SyncDatabaseLock = false
				return
			}
			// 区分group列内容
			switch value := group.(type) {
			case string:
				groupValue = value
			case int:
				groupValue = strconv.Itoa(value)
			case int32:
				groupValue = strconv.Itoa(int(value))
			case int64:
				groupValue = strconv.Itoa(int(value))
			case uint:
				groupValue = strconv.Itoa(int(value))
			}
			// 是否有内容
			if len(groupValue) == 0 {
				global.GVA_LOG.Debug("log groupValue is null")
				SyncDatabaseLock = false
				return
			}
			// 同步完后歇半秒
			global.GVA_LOG.Info(fmt.Sprintf("SyncDatabaseTableData run app: %d", i))
			for e.ForeachInstall(logTable, newTable, keyName, orderName, logColumnList, queryColumnList, addColumn, groupName, groupValue) {
				// 延迟半秒
				time.Sleep(time.Millisecond * 500)
			}
			// i++
			i += 1
		}
	} else {
		// 同步完后歇半秒
		for e.ForeachInstall(logTable, newTable, keyName, orderName, logColumnList, queryColumnList, addColumn, groupName, "") {
			// 延迟半秒
			time.Sleep(time.Millisecond * 500)
		}
	}
	// stop
	SyncDatabaseLock = false
	fmt.Println("SyncDatabaseTableData stop")
}

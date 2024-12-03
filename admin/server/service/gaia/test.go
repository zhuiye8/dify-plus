package gaia

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var sysRegexp = regexp.MustCompile("^sys\\.(.*?)$")
var urlDick = make(map[string]string)

type TestService struct{}

var RunAppList []string
var LOCK bool

// GetAppUrl
// @Tags Test
// @Summary 获取 app 关联的url
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) GetAppUrl(appId string) (string, error) {
	var err error
	if url, ok := urlDick[appId]; ok {
		return url, err
	}
	// 查询数据库
	var app gaia.Apps
	if err = global.GVA_DB.Where("id=?", appId).First(&app).Error; err != nil {
		return "", errors.New("找不到对应appid: " + appId)
	}
	// save
	switch app.Mode {
	case "completion":
		urlDick[appId] = "/v1/completion-messages"
	case "agent-chat", "advanced-chat", "chat":
		urlDick[appId] = "/v1/chat-messages"
	case "workflow":
		urlDick[appId] = "/v1/workflows/run"
	default:
		return "", errors.New("url not found")
	}
	return urlDick[appId], err
}

// GetAppToken
// @Tags Test
// @Summary 获取 app token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) GetAppToken(appid string) (token string, err error) {
	// 获取对应的token
	var tokens gaia.ApiTokens
	if err = global.GVA_DB.Where("app_id=?", appid).First(&tokens).Error; err != nil {
		return "", errors.New(fmt.Sprintf("AppRequestTest Token Error: %s %s", appid, err.Error()))
	}
	//
	return tokens.Token, nil
}

// RunRequest
// @Tags Test
// @Summary 执行gaia请求
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) RunRequest(url, token, jsonData, query string) (id string, err error) {
	// 将请求体编码为 JSON
	client := &http.Client{}
	// 解析JSON字符串为map
	var cacheDick map[string]interface{}
	var data = make(map[string]interface{})
	if err = json.Unmarshal([]byte(jsonData), &cacheDick); err != nil {
		return id, fmt.Errorf("error unmarshalling JSON: %s %s", url, err)
	}
	// 强制阻塞模式
	for key, value := range cacheDick {
		var keyList [][]string
		if keyList = sysRegexp.FindAllStringSubmatch(key, 1); len(keyList) > 0 {
			cacheDick[keyList[0][1]] = value
			delete(cacheDick, key)
		}
	}
	// 强制替换
	if len(query) > 0 {
		data["query"] = query
	}
	data["inputs"] = cacheDick
	data["response_mode"] = "blocking"
	data["user"] = gaia.UsernameUsingApiRequest
	// 将修改后的map重新编码为JSON字符串
	var modifiedJsonStr []byte
	modifiedJsonStr, err = json.Marshal(data)
	if err != nil {
		return id, fmt.Errorf("error marshalling JSON: %s %s", url, err)
	}
	// 创建新的 POST 请求
	req, err := http.NewRequest("POST", fmt.Sprintf(
		"%s%s", global.GVA_CONFIG.Gaia.Url, url), bytes.NewBuffer(modifiedJsonStr))
	if err != nil {
		return id, fmt.Errorf("error creating request: %s %s", url, err)
	}
	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return id, fmt.Errorf("error sending request: %s %v", url, err)
	}
	defer resp.Body.Close()
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return id, fmt.Errorf("request failed with status: %s %s", url, resp.Status)
	}
	var bodyByte []byte
	if bodyByte, err = io.ReadAll(resp.Body); err != nil {
		return id, fmt.Errorf("request io.ReadAll: %s %s", url, resp.Status)
	}
	// 解析获取task_id或者id
	var result map[string]interface{}
	if err = json.Unmarshal(bodyByte, &result); err != nil {
		return id, fmt.Errorf("request json.Unmarsha: %s %s", url, resp.Status)
	}
	// 获取id
	var ok bool
	var cacheID interface{}
	if cacheID, ok = result["id"]; ok {
		// 使用类型断言来判断和转换
		if id, ok = cacheID.(string); ok {
			return id, nil
		} else {
			return id, fmt.Errorf("switch id error: %s %s", url, resp.Status)
		}
	} else if cacheID, ok = result["workflow_run_id"]; ok {
		// 使用类型断言来判断和转换
		if id, ok = cacheID.(string); ok {
			return id, nil
		} else {
			return id, fmt.Errorf("switch id error: %s %s", url, resp.Status)
		}
	} else if cacheID, ok = result["task_id"]; ok {
		// 使用类型断言来判断和转换
		if id, ok = cacheID.(string); ok {
			return id, nil
		} else {
			return id, fmt.Errorf("switch id error: %s %s", url, resp.Status)
		}
	}
	// return
	return id, fmt.Errorf("get id error: %s", url)
}

// SaveTestLog
// @Tags Test
// @Summary 储存测试日志
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @P appID, inputs, outputs, comparison, status, err string, logTime, elapsed float64, batchID uint
func (e *TestService) SaveTestLog(
	appID, inputs, outputs, comparison, status, err string, logTime, elapsed float64, batchID uint) {
	// 查询workflow列表uu
	var appNumber = 0
	var id = uint(1)
	var log gaia.AppRequestTest
	if global.GVA_DB.Select("id").Order("id desc").First(&log).Error == nil {
		id = log.ID + 1
	}
	// 是否新的appid
	if !utils.InStringArray(appID, RunAppList) {
		RunAppList = append(RunAppList, appID)
		appNumber = 1
	}
	// inputs原始 Unicode 转义字符串
	if cacheStr, iErr := strconv.Unquote(inputs); iErr == nil {
		inputs = cacheStr
	}
	// outputs原始 Unicode 转义字符串
	if cacheStr, iErr := strconv.Unquote(outputs); iErr == nil {
		outputs = cacheStr
	}
	// comparison原始 Unicode 转义字符串
	if cacheStr, iErr := strconv.Unquote(comparison); iErr == nil {
		comparison = cacheStr
	}
	// 修改创建
	global.GVA_DB.Create(&gaia.AppRequestTest{
		ID:          id,
		Error:       err,
		AppID:       appID,
		Status:      status,
		Inputs:      inputs,
		Outputs:     outputs,
		BatchId:     batchID,
		Comparison:  comparison,
		LogTime:     math.Round(logTime*100) / 100,
		ElapsedTime: math.Round(elapsed*100) / 100,
	})
	// 是否正常状态
	var failureNumber = 1
	var successNumber = 0
	if status == gaia.MessagesSucceeded || status == gaia.WorkflowSucceeded {
		successNumber = 1
		failureNumber = 0
	}
	// 修改批次状态
	global.GVA_DB.Model(&gaia.AppRequestTestBatch{}).
		Where("id = ?", batchID).
		Updates(&map[string]interface{}{
			"sum":           gorm.Expr("sum + ?", 1),
			"app":           gorm.Expr("app + ?", appNumber),
			"success_count": gorm.Expr("success_count + ?", successNumber),
			"failure_count": gorm.Expr("failure_count + ?", failureNumber),
		})
}

// TestRunWorkflow
// @Tags Test
// @Summary 运行工作流测试
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param appList, endList []string, batchID uint
func (e *TestService) TestRunWorkflow(appList, endList []string, batchID uint) {
	// workflow_runs all
	var err error
	var workflowRun []gaia.WorkflowRun
	if err = global.GVA_DB.Select("app_id").Where(
		"app_id IN (?)", appList).Group("app_id").Find(&workflowRun).Error; err != nil {
		global.GVA_LOG.Debug("AppRequestTest TestRunWorkflow Error: " + err.Error())
		return
	}
	// 提取关联app_id
	for _, v := range workflowRun {
		// 获取token
		var token string
		var appId = v.AppID
		var workflow []gaia.WorkflowRun
		if token, err = e.GetAppToken(appId); err != nil {
			fmt.Println(err.Error())
		}
		// 获取最近10个 end_user的聊天信息
		if err = global.GVA_DB.Select("inputs", "outputs", "elapsed_time").Where(
			"app_id=? AND status=? AND created_by_role=? AND created_by IN (?) AND inputs IS NOT NULL AND NOT (inputs::text = '{}' OR inputs::text = 'null')",
			appId, gaia.WorkflowSucceeded, gaia.IndirectAccessUser, endList).Order("id desc").Limit(
			gaia.TestDefaultNumber).Find(&workflow).Error; err != nil {
			global.GVA_LOG.Debug(fmt.Sprintf("AppRequestTest TestRunWorkflow Error: %s %s", appId, err.Error()))
			continue
		}
		// 执行
		for _, item := range workflow {
			// 提取id
			var id string
			if id, err = e.RunRequest("/v1/workflows/run", token, item.Inputs, ""); err != nil {
				errStr := "workflows/run error" + err.Error()
				e.SaveTestLog(appId, item.Inputs, item.Outputs, errStr, gaia.UserClosed,
					"", item.ElapsedTime, 0, batchID)
				global.GVA_LOG.Debug(errStr)
				continue
			}
			// 查询对应请求详情
			var newWorkflow gaia.WorkflowRun
			if err = global.GVA_DB.Where("id=?", id).First(&newWorkflow).Error; err != nil {
				errStr := "WorkflowRun get new error" + err.Error()
				e.SaveTestLog(appId, item.Inputs, item.Outputs, errStr, gaia.UserClosed,
					newWorkflow.Error, item.ElapsedTime, newWorkflow.ElapsedTime, batchID)
				global.GVA_LOG.Debug(errStr)
				continue
			}
			// create
			e.SaveTestLog(appId, item.Inputs, item.Outputs, newWorkflow.Outputs,
				newWorkflow.Status, newWorkflow.Error, item.ElapsedTime, newWorkflow.ElapsedTime, batchID)
		}
	}
}

// TestRunMessages
// @Tags Test
// @Summary 运行聊天测试
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) TestRunMessages(appList []string, batchID uint) {
	// workflow_runs all
	var err error
	var message []gaia.Messages
	if err = global.GVA_DB.Select("app_id").Where(
		"app_id IN (?)", appList).Group("app_id").Find(&message).Error; err != nil {
		global.GVA_LOG.Debug("AppRequestTest TestRunMessages Error: " + err.Error())
		return
	}
	// 循环聊天列表
	for _, v := range message {
		// 获取token
		var token string
		var appId = v.AppID.String()
		var messages []gaia.Messages
		if token, err = e.GetAppToken(appId); err != nil {
			global.GVA_LOG.Debug(fmt.Sprintf("AppRequestTest GetAppToken Error: %s", err.Error()))
			continue
		}
		// 获取最近10个 end_user的聊天信息
		if err = global.GVA_DB.Select("query", "app_id", "inputs", "answer", "provider_response_latency").Where(
			"app_id=? AND status=? AND from_source=? AND inputs IS NOT NULL AND NOT (inputs::text = '{}' OR inputs::text = 'null')",
			appId, gaia.MessagesSucceeded, gaia.ChatRequestTypeApi).Order("created_at desc").Limit(gaia.TestDefaultNumber).Find(
			&messages).Error; err != nil {
			global.GVA_LOG.Debug(fmt.Sprintf("AppRequestTest TestRunMessages Error: %s %s", appId, err.Error()))
			continue
		}
		// 执行
		asterisk := utils.AddAsteriskToString(token)
		for _, item := range messages {
			// 提取id
			var id, url string
			if url, err = e.GetAppUrl(appId); err != nil {
				errStr := "AppRequestTest TestRunMessages app url error" + err.Error()
				e.SaveTestLog(appId, item.Inputs, item.Answer, errStr, gaia.UserClosed,
					"", item.ProviderResponseLatency, 0, batchID)
				global.GVA_LOG.Debug(errStr)
				continue
			}
			// 请求
			if id, err = e.RunRequest(url, token, item.Inputs, item.Query); err != nil {
				errStr := fmt.Sprintf("AppRequestTest RunRequest error\n%s\ntoken:%s", err.Error(), asterisk)
				e.SaveTestLog(appId, item.Inputs, item.Answer, errStr, gaia.UserClosed,
					"", item.ProviderResponseLatency, 0, batchID)
				global.GVA_LOG.Debug(errStr)
				continue
			}
			// 查询对应请求详情
			var newWorkflow gaia.Messages
			if err = global.GVA_DB.Where("id=?", id).First(&newWorkflow).Error; err != nil {
				errStr := "AppRequestTest TestRunMessages get new error" + err.Error()
				e.SaveTestLog(appId, item.Inputs, item.Answer, errStr, gaia.UserClosed,
					newWorkflow.Error, item.ProviderResponseLatency, newWorkflow.ProviderResponseLatency, batchID)
				global.GVA_LOG.Debug(errStr)
				continue
			}
			// create
			e.SaveTestLog(appId, item.Inputs, item.Answer, newWorkflow.Answer, newWorkflow.Status,
				newWorkflow.Error, item.ProviderResponseLatency, newWorkflow.ProviderResponseLatency, batchID)
		}
	}
}

// AppRequestTest
// @Tags Test
// @Summary 应用请求测试
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *TestService) AppRequestTest() (err error) {
	// 获取APP列表
	var tenantList []string
	var endUser []gaia.EndUser
	var endList, appList []string
	var batch gaia.AppRequestTestBatch
	var tenant []gaia.TenantAccountJoin
	if LOCK {
		return errors.New("AppRequestTest is running")
	}
	if err = global.GVA_DB.Where("account_id=? AND role=?",
		global.GVA_CONFIG.Gaia.SuperAdminAccountId, "owner").Find(&tenant).Error; err != nil {
		return errors.New("AppRequestTest TenantAccountJoin Error: " + err.Error())
	}
	if len(tenant) == -0 {
		return errors.New("AppRequestTest Tenant is null ")
	}
	for _, v := range tenant {
		tenantList = append(tenantList, v.TenantID)
	}
	// 循环获取ADMIN关联空间表
	if err = global.GVA_DB.Where("tenant_id IN (?) AND \"type\"=? AND session_id != ?",
		tenantList, gaia.UserUsingApiRequest, gaia.UsernameUsingApiRequest).Find(&endUser).Error; err != nil {
		return errors.New("AppRequestTest EndUser Error: " + err.Error())
	}
	//
	if len(endUser) == 0 {
		return errors.New("AppRequestTest No EndUser")
	}
	// 循环获取用户列表和app_id列表
	for _, v := range endUser {
		appList = append(appList, v.AppID)
		endList = append(endList, v.ID)
	}
	// 获取最新的batch_id
	LOCK = true
	if err = global.GVA_DB.Order("id desc").First(&batch).Error; err == nil {
		batch.ID += 1
	} else {
		batch.ID = 1
	}
	// 创建批次
	if err = global.GVA_DB.Create(&gaia.AppRequestTestBatch{
		App:          0,
		Sum:          0,
		EndTime:      0,
		SuccessCount: 0,
		FailureCount: 0,
		ID:           batch.ID,
		CreateTime:   time.Now().Unix(),
		Status:       gaia.BatchStatusInProgress,
	}).Error; err != nil {
		LOCK = false
		return errors.New("批次创建失败")
	}
	// 异步请求
	RunAppList = []string{}
	go func(app, end []string, id uint) {
		e.TestRunWorkflow(app, end, id)
		e.TestRunMessages(app, id)
		// 标记结束
		global.GVA_DB.Model(&gaia.AppRequestTestBatch{}).
			Where("id = ?", id).
			Updates(&map[string]interface{}{
				"end_time": time.Now().Unix(),
				"status":   gaia.BatchStatusCompleted,
			})
		LOCK = false
	}(appList, endList, batch.ID)
	return err
}

// AppRequestTestList
// @Tags Test
// @Summary gaia应用请求测试结果列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param info request.PageInfo
// @Return list []response.GetQuotaManagementDataResponse, total int64, err error
func (e *TestService) AppRequestTestList(info request.GetAppRequestTestRequest) (
	_ bool, list []response.GetAppRequestTestDataResponse, total int64, err error) {
	if info.PageSize == 0 {
		info.PageSize = 10
	}
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&gaia.AppRequestTest{}).Where("batch_id = ?", info.BatchId)
	// 筛选app列表
	if len(info.Apps) > 0 {
		db.Where("app_id IN (?)", info.Apps)
	}
	// 是否筛选状态
	switch info.Status {
	case request.GetAppRequestFilterSuccess:
		db.Where("status IN (?)", []string{gaia.MessagesSucceeded, gaia.WorkflowSucceeded})
	case request.GetAppRequestFilterFailure:
		db.Where("status NOT IN (?)", []string{gaia.MessagesSucceeded, gaia.WorkflowSucceeded})
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Order("id desc").Limit(limit).Offset(offset)
	var requestTest []gaia.AppRequestTest
	err = db.Find(&requestTest).Error
	if err != nil {
		err = fmt.Errorf("查询测试信息失败：%s", err.Error())
		return
	}

	// 账号ID集合，方便后面一次性查出
	var appList []uuid.UUID
	for _, money := range requestTest {
		var uid uuid.UUID
		if uid, err = uuid.FromString(money.AppID); err == nil {
			appList = append(appList, uid)
		}
	}

	// 查询用户信息
	var apps []gaia.Apps
	var appInfos = make(map[string]string)
	if len(appList) > 0 {
		err = global.GVA_DB.Model(&gaia.Apps{}).Where("id in (?)", appList).Find(&apps).Error
		if err != nil {
			err = fmt.Errorf("查询应用信息失败：%s", err.Error())
			return
		}
	}
	// 获取appInfos
	for _, app := range apps {
		appInfos[app.ID.String()] = app.Name
	}

	// 拼接结果
	for _, item := range requestTest {
		var status bool
		var name string
		var isExist bool
		if name, isExist = appInfos[item.AppID]; !isExist {
			global.GVA_LOG.Error("AppRequestTestList app信息不存在!", zap.String("app", item.AppID))
			continue
		}
		// 区分状态
		if item.Status == gaia.MessagesSucceeded || item.Status == gaia.WorkflowSucceeded {
			status = true
		}
		// push
		list = append(list, response.GetAppRequestTestDataResponse{
			Name:        name,
			Status:      status,
			Error:       item.Error,
			Inputs:      item.Inputs,
			Outputs:     item.Outputs,
			LogTime:     item.LogTime,
			Comparison:  item.Comparison,
			ElapsedTime: item.ElapsedTime,
		})
	}
	return LOCK, list, total, err
}

// AppRequestTestBatch
// @Tags Test
// @Summary gaia应用请求测试批次列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param info request.PageInfo
// @Return list []response.GetQuotaManagementDataResponse, total int64, err error
func (e *TestService) AppRequestTestBatch(info request.GetAppRequestTestRequest) (
	_ bool, list []gaia.AppRequestTestBatch, total int64, err error) {
	// init
	if info.PageSize == 0 {
		info.PageSize = 10
	}
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&gaia.AppRequestTestBatch{})

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Order("id desc").Limit(limit).Offset(offset)
	// 查询用户信息
	err = db.Find(&list).Error
	if err != nil {
		err = fmt.Errorf("查询测试信息失败：%s", err.Error())
		return
	}
	return LOCK, list, total, err
}

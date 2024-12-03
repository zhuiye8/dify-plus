package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Extend Start: sync user

// SyncUser
// @Tags     Base
// @Summary  用户同步
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /user/sync [post]
func (b *BaseApi) SyncUser(c *gin.Context) {
	userExtendService.SyncUser()
	response.OkWithMessage("同步中", c)
}

// Extend Stop: sync user

// OaLogin
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/oaLogin [post]
func (b *BaseApi) OaLogin(c *gin.Context) {
	var l systemReq.OaLoginReq
	err := c.ShouldBindJSON(&l)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.OaLoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	/**
	请求获取accessToken
	*/
	clientGetToken := resty.New().R()
	var accessTokenResponse *resty.Response
	getTokenUrl := fmt.Sprintf("%s%s", global.GVA_CONFIG.OaLogin.Url, global.GVA_CONFIG.OaLogin.GetTokenByCodeApiPath)
	var postParams = map[string]string{
		"client_id":     global.GVA_CONFIG.OaLogin.Oauth2ClientId,
		"client_secret": global.GVA_CONFIG.OaLogin.Oauth2ClientSecret,
		"code":          l.AuthorizeCode,
		"grant_type":    "authorization_code",
		"redirect_uri":  "",
	}
	accessTokenResponse, err = clientGetToken.
		SetFormData(postParams).
		Post(getTokenUrl)
	if err != nil {
		global.GVA_LOG.Error("请求OA用户信息失败,响应数据为：", zap.Error(errors.New(accessTokenResponse.String())))
		response.FailWithMessage("请求OA用户信息失败："+err.Error(), c)
		return
	}
	tokenRes := accessTokenResponse.String()
	var oaAccessToken systemRes.OaAccessTokenRes
	err = json.Unmarshal([]byte(tokenRes), &oaAccessToken)
	if err != nil {
		global.GVA_LOG.Error("解析OA AccessToken接口返回数据失败,响应数据为：", zap.Error(errors.New(accessTokenResponse.String())))
		response.FailWithMessage("解析OA AccessToken接口返回数据失败："+err.Error(), c)
		return
	}

	/**
	请求OA，返回用户信息
	*/
	getUserInfoUrl := fmt.Sprintf("%s%s", global.GVA_CONFIG.OaLogin.Url, global.GVA_CONFIG.OaLogin.GetUserApiPath)
	clientGetUser := resty.New().R()
	var userInfoResponse *resty.Response
	userInfoResponse, err = clientGetUser.SetHeader("Authorization", oaAccessToken.AccessToken).Post(getUserInfoUrl)
	if err != nil {
		global.GVA_LOG.Error("请求OA用户信息失败,响应数据为：", zap.Error(errors.New(userInfoResponse.String())))
		response.FailWithMessage("请求OA用户信息失败："+err.Error(), c)
		return
	}
	userInfoRes := userInfoResponse.String()
	var oaUserInfo systemRes.OaUserInfoRes
	err = json.Unmarshal([]byte(userInfoRes), &oaUserInfo)
	if err != nil {
		global.GVA_LOG.Error("解析OA用户信息接口返回数据失败,响应数据为：", zap.Error(errors.New(userInfoRes)))
		global.GVA_LOG.Error("解析OA用户信息接口返回数据失败,请求Token为：", zap.String("", oaAccessToken.AccessToken))
		global.GVA_LOG.Error("解析OA用户信息接口返回数据失败,请求Token为：", zap.String("", accessTokenResponse.String()))
		response.FailWithMessage("解析OA用户信息接口返回数据失败："+err.Error(), c)
		return
	}

	// 查询数据库数据
	sysUser := &system.SysUser{}
	err = global.GVA_DB.Where("email", oaUserInfo.Data.Email).First(&sysUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		response.FailWithMessage("查询数据库信息失败："+err.Error(), c)
		return
	}
	// 判断是否需要注册
	if sysUser.ID == 0 {
		// TODO 从ldap中获取用户详细数据

		// 注册用户
		sysUser = &system.SysUser{
			Username:    oaUserInfo.Data.Username,
			NickName:    oaUserInfo.Data.Username,
			HeaderImg:   "https://hn1.oss-cn-shenzhen.aliyuncs.com/w.jpg",
			AuthorityId: system.NormalAuthorityId,
			Authorities: []system.SysAuthority{{AuthorityId: system.NormalAuthorityId}},
			Enable:      1,
			//Phone:       r.Phone, // TODO 手机需要从ldap中获取用户详细数据
			Email:    oaUserInfo.Data.Email,
			Password: utils.RandomString(16),
		}
		var userReturn system.SysUser
		userReturn, err = userService.Register(*sysUser, "")
		if err != nil {
			global.GVA_LOG.Error("注册失败!", zap.Error(err))
			response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败", c)
			return
		} else {
			global.GVA_LOG.Info("注册成功！", zap.Any("username", userReturn.Username))
		}
		sysUser = &userReturn
	}

	var user *system.SysUser
	user, err = userExtendService.OaLogin(sysUser) // 注意这个方法不检查密码
	if err != nil {
		global.GVA_LOG.Error("登陆失败! 用户名不存在!", zap.Error(err))
		response.FailWithMessage("用户名不存在", c)
		return
	}
	if sysUser.Enable != 1 {
		global.GVA_LOG.Error("登陆失败! 用户被禁止登录!")
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	b.TokenNext(c, *user)
	return

}

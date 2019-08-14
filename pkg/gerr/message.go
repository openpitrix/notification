// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import "fmt"

type ErrorMessage struct {
	Name string
	en   string
	zhCN string
}

func (em ErrorMessage) Message(locale string, err error, a ...interface{}) string {
	format := ""
	switch locale {
	case En:
		format = em.en
	case ZhCN:
		if len(em.zhCN) > 0 {
			format = em.zhCN
		} else {
			format = em.en
		}
	}
	if err != nil {
		return fmt.Sprintf("%s: %s", fmt.Sprintf(format, a...), err.Error())
	} else {
		return fmt.Sprintf(format, a...)
	}
}

var (
	ErrorMissingParameter = ErrorMessage{
		Name: "missing_parameter",
		en:   "missing parameter [%s]",
		zhCN: "缺少参数[%s]",
	}
	ErrorGetServiceConfigFailed = ErrorMessage{
		Name: "get_service_config_failed",
		en:   "get_service_config_failed",
		zhCN: "获取通知服务参数配置失败",
	}

	ErrorIllegalEmailFormat = ErrorMessage{
		Name: "illegal_email_format",
		en:   "illegal Email format [%s]",
		zhCN: "非法的Email格式[%s]",
	}
	ErrorIllegalPort = ErrorMessage{
		Name: "illegal_Port",
		en:   "illegal Port [%s]",
		zhCN: "错误的端口号[%s]",
	}
	ErrorDescribeResourcesFailed = ErrorMessage{
		Name: "describe_resources_failed",
		en:   "describe resources failed",
		zhCN: "获取资源失败",
	}
	ErrorCreateResourcesFailed = ErrorMessage{
		Name: "create_resources_failed",
		en:   "create resources failed",
		zhCN: "创建资源失败",
	}
	ErrorUpdateResourceFailed = ErrorMessage{
		Name: "update_resource_failed",
		en:   "update resource [%s] failed",
		zhCN: "更新资源[%s]失败",
	}
	ErrorUnsupportedParameterValue = ErrorMessage{
		Name: "unsupported_parameter_value",
		en:   "unsupported parameter [%s] value [%s]",
		zhCN: "参数[%s]不支持值[%s]",
	}
	ErrorInternalError = ErrorMessage{
		Name: "internal_error",
		en:   "internal error",
		zhCN: "内部错误",
	}
	ErrorValidateFailed = ErrorMessage{
		Name: "validate_failed",
		en:   "validate failed",
		zhCN: "校验失败",
	}
	ErrorDeleteResourceFailed = ErrorMessage{
		Name: "delete_resource_failed",
		en:   "delete resource [%s] failed",
		zhCN: "删除资源[%s]失败",
	}
	ErrorRetryTaskFailed = ErrorMessage{
		Name: "retry_task_failed",
		en:   "retry task [%s] failed",
		zhCN: "重试任务[%s]失败",
	}
	ErrorTaskNotExist = ErrorMessage{
		Name: "task_not_exist",
		en:   "task [%s] not exist",
		zhCN: "任务[%s]不存在",
	}
	ErrorRetryNotificationsFailed = ErrorMessage{
		Name: "retry_notification_failed",
		en:   "retry notification[%s] failed",
		zhCN: "重试通知[%s]失败",
	}
	ErrorNotificationNotExist = ErrorMessage{
		Name: "notification_not_exist",
		en:   "notification [%s] not exist",
		zhCN: "通知[%s]不存在",
	}
	ErrorAddressNotExist = ErrorMessage{
		Name: "address_not_exist",
		en:   "address [%s] not exist",
		zhCN: "地址[%s]不存在",
	}
	ErrorAddressListNotExist = ErrorMessage{
		Name: "address_list_not_exist",
		en:   "address_list [%s] not exist",
		zhCN: "地址列表[%s]不存在",
	}
	ErrorIllegalTimeFormat = ErrorMessage{
		Name: "illegal_time_format",
		en:   "illegal time format [%s]",
		zhCN: "错误的时间格式[%s]",
	}
	ErrorNotAvailableTimeRange = ErrorMessage{
		Name: "not_available_time",
		en:   "not available time [%s]-[%s]",
		zhCN: "不在有效时间范围内[%s]-[%s]",
	}
	ErrorValidateEmailService = ErrorMessage{
		Name: "error_validate_email_service",
		en:   "validate email service failed",
		zhCN: "验证邮件服务配置失败",
	}
	ErrorIllegalNotificationAddressInfo = ErrorMessage{
		Name: "illegal_notification_address_format",
		en:   "illegal notification address format, fmt should be {\"email\": [\"xxx@abc.com\",\"xxx@xxx.com\"]},或者[\"adl-xxx1\", \"adl-xxx2\"].",
		zhCN: "错误的通知地址格式，通知地址格式应为{\"email\": [\"xxx@abc.com\",\"xxx@xxx.com\"]},或者[\"adl-xxx1\", \"adl-xxx2\"].",
	}
	ErrorIllegalNotificationType = ErrorMessage{
		Name: "illegal_notification_type",
		en:   "illegal notification type [%s]",
		zhCN: "错误的通知类型[%s]",
	}
	ErrorIllegalNotificationAddressList = ErrorMessage{
		Name: "illegal_notification_address_list",
		en:   "illegal notification address list [%s]",
		zhCN: "错误的通知列表[%s]",
	}
	ErrorDecodeContentFailed = ErrorMessage{
		Name: "error_decode_content_failed",
		en:   "error decode content failed, fmt should be [{\"content_type\": \"content\"}].",
		zhCN: "解码内容失败，内容字段格式应为[{\"content_type\": \"content\"}].",
	}
	ErrorExistAddress = ErrorMessage{
		Name: "error_exist_address",
		en:   "error exist address [%s]",
		zhCN: "已存在该地址信息[%s]",
	}
	ErrorNotifyType = ErrorMessage{
		Name: "error_notify_type",
		en:   "error notify type [%s]",
		zhCN: "错误的通知类型[%s]",
	}
	ErrorNotExistItemInList = ErrorMessage{
		Name: "not_exist_item_in_list",
		en:   "not exist item in list [%s]",
		zhCN: "列表[%s]中某些元素不存在",
	}
	ErrorAlreadyDeletedItemInList = ErrorMessage{
		Name: "already_deleted_item_in_list",
		en:   "already deleted item in list [%s]",
		zhCN: "列表[%s]中某些元素已经被删除",
	}
	ErrorIllegalItemInList = ErrorMessage{
		Name: "illegal_item_in_list",
		en:   "illegal item in list [%s]",
		zhCN: "列表[%s]中某些元素不合法",
	}
	ErrorExtraIsNotJsonFmt = ErrorMessage{
		Name: "is_not_json_fmt",
		en:   "is not json fmt [%s]",
		zhCN: "不是JSON格式[%s]",
	}
	ErrorResourceNotExist = ErrorMessage{
		Name: "resource_not_exist",
		en:   "resource [%s] not exist",
		zhCN: "资源[%s]不存在",
	}
	ErrorAlreadyExistResource = ErrorMessage{
		Name: "error_already_exist_resource",
		en:   "error already exist resource [%s]",
		zhCN: "已存在该资源[%s]",
	}
	ErrorIllegalWebsocketDisabled = ErrorMessage{
		Name: "illegal_websocket_disabled",
		en:   "illegal websocket disabled [%s]",
		zhCN: "websocket配置尚未启用[%s]",
	}
	ErrorIllegalNotificationExtraBlank = ErrorMessage{
		Name: "illegal_notification_extra_blank",
		en:   "illegal notification extra, extra should be blank.",
		zhCN: "未设置websocket用户信息，extra应为空.",
	}
	ErrorIllegalNotificationExtra = ErrorMessage{
		Name: "illegal_notification_extra",
		en:   "illegal notification extra [%s]",
		zhCN: "错误的通知附加信息[%s]",
	}
	ErrorDecodeExtraFailed = ErrorMessage{
		Name: "error_decode_extra_failed",
		en:   "error decode extra failed, fmt should be [{\"ws_service\": \"xxx\",\"ws_message_type\": \"xxx\"}].",
		zhCN: "解码extra字段失败，extra字段格式应为[{\"ws_service\": \"xxx\",\"ws_message_type\": \"xxx\"}].",
	}
)

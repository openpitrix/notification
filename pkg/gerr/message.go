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
	ErrorRetryTaskNotExist = ErrorMessage{
		Name: "retry_task_not_exist",
		en:   "retry task [%s] not exist",
		zhCN: "重试任务[%s]不存在",
	}
	ErrorRetryNotificationsFailed = ErrorMessage{
		Name: "retry_notification_failed",
		en:   "retry notification[%s] failed",
		zhCN: "重发通知[%s]失败",
	}
	ErrorRetryNotificationtNotExist = ErrorMessage{
		Name: "retry_notification_not_exist",
		en:   "retry notification [%s] not exist",
		zhCN: "重试通知[%s]不存在",
	}
)

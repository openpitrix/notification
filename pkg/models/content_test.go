package models

import (
	"testing"

	"openpitrix.io/logger"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
)

func TestDecodeContent(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()

	testContentStr := "{\"html\":\"register_content_html\",  \"normal\":\"register_content_normal\"}"

	contentStruct, err := DecodeContent(testContentStr)
	if err != nil {
		logger.Errorf(nil, "error=[%+v]", err)
	}

	content2SendHtml := ""
	content, ok := (*contentStruct)[constants.ContentFmtHtml]
	if ok {
		content2SendHtml = content
	}
	logger.Debugf(nil, "content2SendHtml=[%s]", content2SendHtml)

	content2SendNormal := ""
	content, ok = (*contentStruct)[constants.ContentFmtNormal]
	if ok {
		content2SendNormal = content
	}
	logger.Debugf(nil, "content2SendNormal=[%s]", content2SendNormal)

}

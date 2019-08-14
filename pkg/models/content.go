package models

import (
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

type Content map[string]string

func DecodeContent(data string) (*Content, error) {
	content := new(Content)
	err := jsonutil.Decode([]byte(data), content)
	if err != nil {
		logger.Warnf(nil, "Try to decode as format[{\"content_type\": \"content\"}], decode [%s] into content failed: %+v", data, err)
		return nil, gerr.NewWithDetail(nil, gerr.InvalidArgument, err, gerr.ErrorDecodeContentFailed)
	}
	return content, err
}

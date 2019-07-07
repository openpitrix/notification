// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
)

func checkEmailCfgInDBisDefaultOrNot(emailCfg models.EmailConfig) bool {
	if emailCfg.EmailHost == "mail.app-center.com.cn" && emailCfg.Email == "admin@app-center.com.cn" && emailCfg.DisplaySender == "notification" {
		return true
	} else {
		return false
	}
}
func modifyEmailConfigFromCfg(cfg *config.Config) error {
	attributes := make(map[string]interface{})
	attributes[models.EmailCfgColProtocol] = cfg.Email.Protocol
	attributes[models.EmailCfgColEmailHost] = cfg.Email.EmailHost
	attributes[models.EmailCfgColPort] = (uint32(cfg.Email.Port))
	attributes[models.EmailCfgColDisplaySender] = cfg.Email.DisplaySender
	attributes[models.EmailCfgColPassword] = cfg.Email.Password
	attributes[models.EmailCfgColSSLEnable] = cfg.Email.SSLEnable
	attributes[models.EmailCfgColEmail] = cfg.Email.Email
	attributes[models.EmailCfgColStatusTime] = time.Now()

	tx := global.GetInstance().GetDB().Begin()
	_, err := GetEmailConfig4Modify(nil, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table(models.TableEmailConfig).Updates(attributes).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	resetConfig4EmailCfgFromConfig(cfg)
	logger.Debugf(nil, "Modify email config from config successfully, %+v.", config.GetInstance().Email)
	return nil
}
func ResetEmailCfg(cfg *config.Config) error {
	emailCfg, err := GetEmailConfigFromDB(nil)
	if err != nil {
		logger.Errorf(nil, "Failed to reset email config, %+v.", err)
		return err
	} else {
		isDefault := checkEmailCfgInDBisDefaultOrNot(*emailCfg)
		if isDefault {
			err = modifyEmailConfigFromCfg(cfg)
			if err != nil {
				logger.Errorf(nil, "Failed to reset email config, %+v.", err)
				return err
			} else {
				logger.Infof(nil, "Reset email config successfully, %+v.", config.GetInstance().Email)
				return nil
			}
		} else {
			resetConfig4EmailCfg(emailCfg)
			logger.Infof(nil, "Reset email config successfully, %+v.", config.GetInstance().Email)
			return nil
		}
	}
}

func getModifyAttributes(req *pb.ServiceConfig) map[string]interface{} {
	attributes := make(map[string]interface{})
	if req.EmailServiceConfig.GetProtocol().GetValue() != "" {
		attributes[models.EmailCfgColProtocol] = req.EmailServiceConfig.GetProtocol().GetValue()
	}

	if req.EmailServiceConfig.GetEmailHost().GetValue() != "" {
		attributes[models.EmailCfgColEmailHost] = req.EmailServiceConfig.GetEmailHost().GetValue()
	}

	if req.EmailServiceConfig.GetPort().GetValue() != 0 {
		attributes[models.EmailCfgColPort] = req.EmailServiceConfig.GetPort().GetValue()
	}

	if req.EmailServiceConfig.GetDisplaySender().GetValue() != "" {
		attributes[models.EmailCfgColDisplaySender] = req.EmailServiceConfig.GetDisplaySender().GetValue()
	}

	if req.EmailServiceConfig.GetEmail().GetValue() != "" {
		attributes[models.EmailCfgColEmail] = req.EmailServiceConfig.GetEmail().GetValue()
	}

	if req.EmailServiceConfig.GetPassword().GetValue() != "" {
		attributes[models.EmailCfgColPassword] = req.EmailServiceConfig.GetPassword().GetValue()
	}

	if req.EmailServiceConfig.GetSslEnable().GetValue() != false {
		attributes[models.EmailCfgColSSLEnable] = req.EmailServiceConfig.GetSslEnable().GetValue()
	}
	return attributes
}
func ModifyEmailConfig(ctx context.Context, req *pb.ServiceConfig) error {
	attributes := getModifyAttributes(req)
	attributes[models.EmailCfgColStatusTime] = time.Now()

	tx := global.GetInstance().GetDB().Begin()
	_, err := GetEmailConfig4Modify(ctx, tx)
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to modify email config, %+v.", err)
		return err
	}

	err = tx.Table(models.TableEmailConfig).Updates(attributes).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to modify email config, %+v.", err)
		return err
	}
	tx.Commit()

	resetConfig4EmailCfgFromReq(req)
	logger.Debugf(nil, "Modify email config successfully, %+v.", config.GetInstance().Email)
	return nil
}

func GetEmailConfig4Modify(ctx context.Context, tx *gorm.DB) (*models.EmailConfig, error) {
	var emailConfig models.EmailConfig
	err := tx.Table(models.TableEmailConfig).Last(&emailConfig).Set("gorm:query_option", "FOR UPDATE").Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to get email config for modify, %+v.", err)
		return nil, err
	}
	return &emailConfig, nil

}

func GetEmailConfigFromDB(ctx context.Context) (*models.EmailConfig, error) {
	var emailConfigs models.EmailConfig
	db := global.GetInstance().GetDB()
	err := db.Last(&emailConfigs).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get email config from DB, %+v.", err)
		return nil, err
	}
	return &emailConfigs, nil
}

func resetConfig4EmailCfg(emailCfg *models.EmailConfig) {
	os.Setenv("NOTIFICATION_EMAIL_PROTOCOL", emailCfg.Protocol)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL_HOST", emailCfg.EmailHost)
	p := strconv.Itoa(int(emailCfg.Port))
	os.Setenv("NOTIFICATION_EMAIL_PORT", p)
	os.Setenv("NOTIFICATION_EMAIL_DISPLAY_SENDER", emailCfg.DisplaySender)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL", emailCfg.Email)
	os.Setenv("NOTIFICATION_EMAIL_PASSWORD", emailCfg.Password)
	os.Setenv("NOTIFICATION_EMAIL_SSL_ENABLE", strconv.FormatBool(emailCfg.SSLEnable))

	config.GetInstance().LoadConf()
}
func resetConfig4EmailCfgFromReq(req *pb.ServiceConfig) {
	os.Setenv("NOTIFICATION_EMAIL_PROTOCOL", req.EmailServiceConfig.GetProtocol().GetValue())
	os.Setenv("NOTIFICATION_EMAIL_EMAIL_HOST", req.EmailServiceConfig.GetEmailHost().GetValue())
	p := strconv.Itoa(int(req.EmailServiceConfig.GetPort().GetValue()))
	os.Setenv("NOTIFICATION_EMAIL_PORT", p)
	os.Setenv("NOTIFICATION_EMAIL_DISPLAY_SENDER", req.EmailServiceConfig.GetDisplaySender().GetValue())
	os.Setenv("NOTIFICATION_EMAIL_EMAIL", req.EmailServiceConfig.GetEmail().GetValue())
	os.Setenv("NOTIFICATION_EMAIL_PASSWORD", req.EmailServiceConfig.GetPassword().GetValue())
	s := strconv.FormatBool(req.EmailServiceConfig.GetSslEnable().GetValue())
	os.Setenv("NOTIFICATION_EMAIL_SSL_ENABLE", s)

	config.GetInstance().LoadConf()
}
func resetConfig4EmailCfgFromConfig(cfg *config.Config) {
	os.Setenv("NOTIFICATION_EMAIL_PROTOCOL", cfg.Email.Protocol)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL_HOST", cfg.Email.EmailHost)
	p := strconv.Itoa(int(cfg.Email.Port))
	os.Setenv("NOTIFICATION_EMAIL_PORT", p)
	os.Setenv("NOTIFICATION_EMAIL_DISPLAY_SENDER", cfg.Email.DisplaySender)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL", cfg.Email.Email)
	os.Setenv("NOTIFICATION_EMAIL_PASSWORD", cfg.Email.Password)
	os.Setenv("NOTIFICATION_EMAIL_SSL_ENABLE", strconv.FormatBool(cfg.Email.SSLEnable))

	config.GetInstance().LoadConf()
}

package nf

import (
	"openpitrix.io/notification/pkg/models"
)

//import "nf/pkg/nf/models"

// Service interface describes all functions that must be implemented.
type Service interface {
	SayHello(str string) (string, error)
	CreateNfWaddrs(*models.NotificationCenterPost) (nfPostID string, err error)
	DescribeNfs(nfID string) (*models.NotificationCenterPost, error)

	//FindByEmail(email string) (*models.NotificationCenterPost, error)
	//EncryptPassword(password string) ([]byte, error)
	//ComparePasswords(hashedPassword []byte, password string) bool
	//CreateUser(email, password, firstName, lastName string, status UserStatus) (*models.NotificationCenterPost, error)
	//ActivateUser(email string) error
	//ChangePassword(email, password string) error


	//GetDataFromDB4Test()
	//CreateNfWaddrs(nfPostID string, nfPostType string, title string, content string, shortContent string, exporedDays int64,owner string ) (error)


}

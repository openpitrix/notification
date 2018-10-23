package nf

import "notification/pkg/models"

//import "nf/pkg/nf/models"

// Service interface describes all functions that must be implemented.
type Service interface {
	//FindByEmail(email string) (*models.NotificationCenterPost, error)
	//EncryptPassword(password string) ([]byte, error)
	//ComparePasswords(hashedPassword []byte, password string) bool
	//CreateUser(email, password, firstName, lastName string, status UserStatus) (*models.NotificationCenterPost, error)
	//ActivateUser(email string) error
	//ChangePassword(email, password string) error

	SayHello(str string)(string, error)
	GetDataFromDB4Test()
	CreateNfWaddrs(nfPostID string, nfPostType string, title string, content string, shortContent string, exporedDays int64,owner string ) (error)
	CreateNfWaddrs2(*models.NotificationCenterPost) error
	CreateNfWaddrs3(nf *models.NotificationCenterPost,job *models.Job,task *models.Task) error
}

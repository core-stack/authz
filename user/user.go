package user

import (
	"time"

	"github.com/core-stack/authz/email"
	"github.com/core-stack/authz/zrepository"
)

type UserService struct {
	userRepo    zrepository.IUserRepository
	codeRepo    zrepository.ICodeRepository
	emailSender email.ISender

	defaultCodeDuration time.Duration
}

func NewUserService(userRepo zrepository.IUserRepository, codeRepo zrepository.ICodeRepository, emailSender email.ISender, defaultCodeDuration time.Duration) *UserService {
	return &UserService{userRepo: userRepo, codeRepo: codeRepo, emailSender: emailSender, defaultCodeDuration: defaultCodeDuration}
}

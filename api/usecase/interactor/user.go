package interactor

import (
	"errors"
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/email"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/constructor"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailNotChanged   = errors.New("email not changed")
	ErrEmailAlreadyUsed  = errors.New("email already used")
)

type UserUseCase struct {
	clock       output_port.Clock
	email       output_port.Email
	ulid        output_port.ULID
	transaction output_port.Transaction
	userAuth    output_port.UserAuth
	userRepo    output_port.UserRepository
}

func NewUserUseCase(
	clock output_port.Clock,
	email output_port.Email,
	ulid output_port.ULID,
	transaction output_port.Transaction,
	userAuth output_port.UserAuth,
	userRepo output_port.UserRepository,
) input_port.IUserUseCase {
	return &UserUseCase{
		clock:       clock,
		email:       email,
		ulid:        ulid,
		transaction: transaction,
		userAuth:    userAuth,
		userRepo:    userRepo,
	}
}

func (u *UserUseCase) Login(loginID, password string) (entity.User, string, error) {
	user, err := u.userRepo.FindByLoginID(loginID)
	if err != nil {
		return entity.User{}, "", err
	}

	token, err := u.userAuth.IssueUserToken(user, u.clock.Now())
	if err != nil {
		return entity.User{}, "", err
	}

	return user, token, nil
}

func (u *UserUseCase) Authenticate(token string) (string, error) {
	return u.userAuth.Authenticate(token)
}

func (u *UserUseCase) AuthenticateForUpdateEmail(token string) (string, error) {
	return u.userAuth.AuthenticateForUpdateEmail(token)
}

func (u *UserUseCase) AuthenticateForUpdatePassword(token string) (string, error) {
	return u.userAuth.AuthenticateForUpdatePassword(token)
}

func (u *UserUseCase) FindByID(_ entity.User, userID string) (entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *UserUseCase) FindByIDByAdmin(_ entity.User, userID string) (entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *UserUseCase) Create(userCreate input_port.UserCreate) (entity.User, error) {
	userID := u.ulid.GenerateID()
	user, err := constructor.NewUserCreate(
		userID,
		userCreate.StudentID,
		userCreate.IdmUniv,
		userCreate.IdmBus,
		userCreate.UserType,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to construct user: %w", err)
	}

	if err := u.userRepo.Create(user); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	if ret, err := u.userRepo.FindByID(userID); err != nil {
		return entity.User{}, fmt.Errorf("failed to find user: %w", err)
	} else {
		return ret, nil
	}
}

func (u *UserUseCase) CreateUserWithDetail(user entity.User) error {
	return u.transaction.StartTransaction(func(tx interface{}) error {
		if err := u.userRepo.CreateWithTx(tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
}

func (u *UserUseCase) Search(_ entity.User, query string, userType string, skip int, limit int) ([]entity.User, int, error) {
	if users, total, err := u.userRepo.Search(query, userType, skip, limit); err != nil {
		return nil, 0, fmt.Errorf("failed to search user: %w", err)
	} else {
		return users, total, nil
	}
}

func (u *UserUseCase) Update(_ entity.User, user input_port.UserUpdate) (entity.User, error) {
	// TODO: 更新処理
	if user, err := u.userRepo.FindByID(user.ID); err != nil {
		return entity.User{}, fmt.Errorf("failed to find user: %w", err)
	} else {
		return user, nil
	}
}

func (u *UserUseCase) Delete(_ entity.User, userID string) (entity.User, error) {
	res, err := u.userRepo.FindByID(userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	err = u.userRepo.Delete(userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to delete user: %w", err)
	}

	return res, nil
}

func (u *UserUseCase) SendResetPasswordMail(emailAddress string) error {
	if emailAddress == "" {
		return fmt.Errorf("%s: email address is empty", ErrKind.BadRequest)
	}

	user, err := u.userRepo.FindByEmail(emailAddress)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	token, err := u.userAuth.IssueUserTokenForUpdatePassword(user, u.clock.Now())
	if err != nil {
		return fmt.Errorf("failed to issue user token: %w", err)
	}

	subject, body := email.ContentToResetPassword(token)
	err = u.email.Send([]string{emailAddress}, subject, body, "")
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

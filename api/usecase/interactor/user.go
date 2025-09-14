package interactor

import (
	"errors"
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/authentication"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/email"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/constructor"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

var (
	ErrUserAlreadyExists            = errors.New("user already exists")
	ErrEmailNotChanged              = errors.New("email not changed")
	ErrEmailAlreadyUsed             = errors.New("email already used")
	ErrAuthenticationCodeExpired    = errors.New("authentication code expired")
	ErrAuthenticationCodeInvalid    = errors.New("authentication code invalid")
	ErrRegisterVerificationNotFound = errors.New("register verification not found")
)

type UserUseCase struct {
	clock                    output_port.Clock
	email                    output_port.Email
	ulid                     output_port.ULID
	transaction              output_port.Transaction
	userAuth                 output_port.UserAuth
	userRepo                 output_port.UserRepository
	authCode                 output_port.AuthCode
	registerVerificationRepo output_port.RegisterVerificationRepository
}

func NewUserUseCase(
	clock output_port.Clock,
	email output_port.Email,
	ulid output_port.ULID,
	transaction output_port.Transaction,
	userAuth output_port.UserAuth,
	userRepo output_port.UserRepository,
	authCode output_port.AuthCode,
	registerVerificationRepo output_port.RegisterVerificationRepository,

) input_port.IUserUseCase {
	return &UserUseCase{
		clock:                    clock,
		email:                    email,
		ulid:                     ulid,
		transaction:              transaction,
		userAuth:                 userAuth,
		userRepo:                 userRepo,
		authCode:                 authCode,
		registerVerificationRepo: registerVerificationRepo,
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
		userCreate.Name,
		userCreate.Age,
		userCreate.UserType,
		userCreate.Email,
		userCreate.Password,
		userCreate.GofileToken,
		userCreate.EmailVerified,
		userCreate.IsDeleted,
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

func (u *UserUseCase) Update(user entity.User, userUpdate input_port.UserUpdate) (entity.User, error) {
	// TODO: 更新処理
	updatingUser, err := constructor.NewUserUpdate(
		user.ID,
		userUpdate.Name,
		userUpdate.Bio,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to construct user update: %w", err)
	}

	if updatingUser.ID != user.ID {
		return entity.User{}, fmt.Errorf("cannot update other user")
	}

	user.Name = updatingUser.Name
	user.Bio = updatingUser.Bio
	if err := u.userRepo.Update(user); err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	res, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return res, nil
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

func (u *UserUseCase) Boot(user entity.User) (_ entity.User, token string, err error) {
	if user.ID == "" {
		guestUser, err := constructor.NewUserCreate(
			u.ulid.GenerateID(),
			"ゲストユーザー",
			1,
			"GuestUser",
			nil,
			nil,
			nil,
			false,
			false,
		)
		if err != nil {
			return entity.User{}, "", fmt.Errorf("failed to construct guest user: %w", err)
		}

		if err := u.userRepo.Create(guestUser); err != nil {
			return entity.User{}, "", fmt.Errorf("failed to create guest user: %w", err)
		}

		user, err := u.userRepo.FindByID(guestUser.ID)
		if err != nil {
			return entity.User{}, "", fmt.Errorf("failed to find guest user: %w", err)
		}

		token, err := u.userAuth.IssueUserToken(user, u.clock.Now())
		if err != nil {
			return entity.User{}, "", fmt.Errorf("failed to issue user token: %w", err)
		}
		return user, token, nil
	}

	return user, "", nil
}

func (u *UserUseCase) CreateByMe(create input_port.CreateByMe) error {
	_, err := u.userRepo.FindByEmail(create.Email)
	if err == nil {
		return errors.Join(ErrKind.Conflict, ErrEmailAlreadyUsed)
	}
	if !errors.Is(err, ErrKind.NotFound) {
		return err
	}

	hashedPassword, err := u.userAuth.HashPassword(create.Password)
	if err != nil {
		return err
	}

	authenticationCode := u.authCode.Generate4DigitCode()

	hashedAuthenticationCode, err := authentication.HashBcryptPassword(authenticationCode)
	if err != nil {
		return err
	}

	registerVerification, err := constructor.NewRegisterVerificationCreate(constructor.NewRegisterVerificationCreateArgs{
		RegisterVerificationID:   u.ulid.GenerateID(),
		Email:                    create.Email,
		Password:                 create.Password,
		HashedPassword:           hashedPassword,
		HashedAuthenticationCode: hashedAuthenticationCode,
		ExpiresAt:                u.clock.Now().Add(output_port.TokenRegisterExpireDuration),
	})
	if err != nil {
		return err
	}

	return u.transaction.StartTransaction(func(tx interface{}) error {
		if err = u.registerVerificationRepo.UpsertInTx(tx, registerVerification); err != nil {
			return err
		}

		subject, body := email.ContentToRegister(authenticationCode)
		return u.email.Send([]string{create.Email}, subject, body, "")
	})
}

func (u UserUseCase) VerifyEmail(user entity.User, input input_port.VerifyEmail) (entity.User, string, error) {
	registerVerification, err := u.registerVerificationRepo.FindByEmail(input.Email)
	if err != nil {
		return entity.User{}, "", errors.Join(ErrKind.NotFound, ErrRegisterVerificationNotFound)
	}

	now := u.clock.Now()
	if registerVerification.ExpiresAt.Before(now) {
		return entity.User{}, "", errors.Join(ErrKind.BadRequest, ErrAuthenticationCodeExpired)
	}

	err = u.userAuth.VerifyAuthenticationCode(registerVerification.HashedAuthenticationCode, input.AuthenticationCode)
	if err != nil {
		return entity.User{}, "", errors.Join(ErrKind.BadRequest, ErrAuthenticationCodeInvalid)
	}

	user.UserType = entconst.MemberUser
	user.Email = &registerVerification.Email
	user.HashedPassword = &registerVerification.HashedPassword
	user.EmailVerified = true

	if err := u.transaction.StartTransaction(func(tx interface{}) error {
		if err := u.registerVerificationRepo.DeleteByEmailInTx(tx, registerVerification.Email); err != nil {
			return err
		}

		if err := u.userRepo.UpdateWithTx(tx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return entity.User{}, "", err
	}

	res, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		return entity.User{}, "", err
	}

	token, err := u.userAuth.IssueUserToken(res, now)
	if err != nil {
		return entity.User{}, "", err
	}

	return res, token, nil
}

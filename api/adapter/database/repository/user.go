package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type UserRepository struct {
	db   *gorm.DB
	ulid output_port.ULID
}

func NewUserRepository(
	db *gorm.DB,
	ulid output_port.ULID,
) output_port.UserRepository {
	return &UserRepository{db: db, ulid: ulid}
}

func (r *UserRepository) create(tx *gorm.DB, user entity.User) (err error) {
	defer output_port.WrapDatabaseError(&err)

	m := &model.User{
		ID:          user.ID,
		Name:        user.Name,
		Age:         user.Age,
		UserType:    user.UserType.String(),
		GofileToken: user.GofileToken,
	}
	if err = tx.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Create(user entity.User) error {
	return r.create(r.db, user)
}

func (r *UserRepository) CreateWithTx(tx interface{}, user entity.User) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return output_port.ErrInvalidTransaction
	}

	return r.create(txAsserted, user)
}

func (r *UserRepository) Delete(ID string) (err error) {
	defer output_port.WrapDatabaseError(&err)

	return r.db.Model(&model.User{}).
		Where("id = ?", ID).
		Updates(
			map[string]interface{}{
				"email":           ID + "_deleted",
				"hashed_password": "",
				"is_deleted":      true,
			},
		).Error
}

func (r *UserRepository) FindByEmail(email string) (user entity.User, err error) {
	defer output_port.WrapDatabaseError(&err)

	if users, err := r.listByEmails(r.db, []string{email}); err != nil {
		return entity.User{}, err
	} else if len(users) == 0 {
		return entity.User{}, fmt.Errorf("%w: user", interactor.ErrKind.NotFound)
	} else {
		return users[0], nil
	}
}

func (r *UserRepository) ListByEmails(emails []string) (_ []entity.User, err error) {
	defer output_port.WrapDatabaseError(&err)
	return r.listByEmails(r.db, emails)
}

func (r *UserRepository) FindByID(ID string) (user entity.User, err error) {
	defer output_port.WrapDatabaseError(&err)

	res := model.User{}
	err = r.db.Model(&model.User{}).
		Where("id = ?", ID).
		First(&res).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.User{}, fmt.Errorf("%w: user", interactor.ErrKind.NotFound)
	}
	if err != nil {
		return entity.User{}, err
	}
	return res.Entity(), nil
}

func (r *UserRepository) FindByLoginID(loginID string) (user entity.User, err error) {
	defer output_port.WrapDatabaseError(&err)
	if users, err := r.listByLoginIDs(r.db, []string{loginID}); err != nil {
		return entity.User{}, err
	} else if len(users) == 0 {
		return entity.User{}, fmt.Errorf("%w: user", interactor.ErrKind.NotFound)
	} else {
		return users[0], nil
	}
}

func (r *UserRepository) ListByLoginIDs(loginIDs []string) ([]entity.User, error) {
	return r.listByLoginIDs(r.db, loginIDs)
}

func (r *UserRepository) Search(query string, userType string, skip int, limit int) (users []entity.User, total int, err error) {
	defer output_port.WrapDatabaseError(&err)

	var res []*model.User
	var totalRes int64
	sqlQuery := r.db.Model(&model.User{}).
		Preload("Member.Occupations").
		Preload("Member.ExternalOrganization").
		Preload("Member.Prefecture").
		Preload("Member.OfficePrefecture").
		Where("email LIKE ?", "%"+query+"%").
		Where("is_deleted = false")

	if userType != "" {
		sqlQuery = sqlQuery.Where("user_type = ?", userType)
	}

	err = sqlQuery.
		Group("users.id").
		Count(&totalRes).
		Offset(skip).
		Limit(limit).
		Find(&res).
		Error
	if err != nil {
		return nil, 0, err
	}

	users = make([]entity.User, len(res))
	for i, user := range res {
		users[i] = user.Entity()
	}

	return users, int(totalRes), nil
}

// Update TODO 更新の実装を修正する
func (r *UserRepository) update(tx *gorm.DB, user entity.User) (err error) {
	defer output_port.WrapDatabaseError(&err)
	return tx.Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(
			map[string]interface{}{
				"id":              user.ID,
				"name":            user.Name,
				"age":             user.Age,
				"user_type":       user.UserType.String(),
				"email":           user.Email,
				"hashed_password": user.HashedPassword,
				"gofile_token":    user.GofileToken,
				"email_verified":  user.EmailVerified,
				"is_deleted":      user.IsDeleted,
			},
		).Error
}

func (r *UserRepository) Update(user entity.User) error {
	return r.update(r.db, user)
}

func (r *UserRepository) UpdateWithTx(tx interface{}, user entity.User) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return output_port.ErrInvalidTransaction
	}
	return r.update(txAsserted, user)
}

func (r *UserRepository) listByEmails(tx *gorm.DB, emails []string) ([]entity.User, error) {
	var res []model.User
	err := tx.Model(&model.User{}).
		Where("is_deleted = false").
		Where("email IN ?", emails).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(res))
	for i, user := range res {
		users[i] = user.Entity()
	}

	return users, nil
}

func (r *UserRepository) listByLoginIDs(tx *gorm.DB, loginIDs []string) ([]entity.User, error) {
	var res []model.User
	err := tx.Model(&model.User{}).
		Where("is_deleted = false").
		Where("login_id IN ?", loginIDs).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(res))
	for i, user := range res {
		users[i] = user.Entity()
	}

	return users, nil
}

func (r *UserRepository) FindMaxLoginID() (string, error) {
	return r.findMaxLoginID(r.db)
}

func UserEntityToModel(user entity.User) *model.User {
	return &model.User{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}
}

func (r *UserRepository) findMaxLoginID(tx *gorm.DB) (string, error) {
	var maxLoginID string
	err := tx.Model(&model.User{}).
		Select("login_id").
		Order("login_id DESC").
		Limit(1).
		Find(&maxLoginID).
		Error
	if err != nil {
		return "", err
	}
	return maxLoginID, nil
}

func (r *UserRepository) GetAdminUser() ([]entity.User, error) {
	var res []model.User
	err := r.db.Model(&model.User{}).
		Where("user_type = ?", entconst.SystemAdmin.String()).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(res))
	for i, user := range res {
		users[i] = user.Entity()
	}

	return users, nil
}

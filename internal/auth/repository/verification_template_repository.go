package repository

import (
	"context"

	"github.com/anychat/server/internal/auth/model"
	"gorm.io/gorm"
)

type VerificationTemplateRepository interface {
	GetByPurpose(ctx context.Context, purpose model.VerificationPurpose) (*model.VerificationTemplate, error)
	GetActive(ctx context.Context) ([]*model.VerificationTemplate, error)
	Update(ctx context.Context, template *model.VerificationTemplate) error
	WithTx(tx *gorm.DB) VerificationTemplateRepository
}

type verificationTemplateRepositoryImpl struct {
	db *gorm.DB
}

func NewVerificationTemplateRepository(db *gorm.DB) VerificationTemplateRepository {
	return &verificationTemplateRepositoryImpl{db: db}
}

func (r *verificationTemplateRepositoryImpl) GetByPurpose(ctx context.Context, purpose model.VerificationPurpose) (*model.VerificationTemplate, error) {
	var template model.VerificationTemplate
	err := r.db.WithContext(ctx).
		Where("purpose = ? AND is_active = ?", purpose, true).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *verificationTemplateRepositoryImpl) GetActive(ctx context.Context) ([]*model.VerificationTemplate, error) {
	var templates []*model.VerificationTemplate
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&templates).Error
	return templates, err
}

func (r *verificationTemplateRepositoryImpl) Update(ctx context.Context, template *model.VerificationTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *verificationTemplateRepositoryImpl) WithTx(tx *gorm.DB) VerificationTemplateRepository {
	return &verificationTemplateRepositoryImpl{db: tx}
}

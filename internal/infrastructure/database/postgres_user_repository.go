package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"jdgonzalez907/users-api/internal/domain"
	"jdgonzalez907/users-api/internal/infrastructure/database/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *PostgresUserRepository) FindById(id int) (*domain.User, error) {
	row, err := r.queries.FindUserByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToDomain(
		row.ID, row.IdentificationType, row.IdentificationNumber, row.FirstName, row.LastName,
		row.BirthDate, row.Address, row.PhoneCountryCode, row.PhoneNumber, row.Email,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (r *PostgresUserRepository) FindByPhone(phone domain.Phone) (*domain.User, error) {
	dto := phone.ToDTO()
	row, err := r.queries.FindUserByPhone(context.Background(), sqlc.FindUserByPhoneParams{
		PhoneCountryCode: dto.CountryCode,
		PhoneNumber:      dto.Number,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToDomain(
		row.ID, row.IdentificationType, row.IdentificationNumber, row.FirstName, row.LastName,
		row.BirthDate, row.Address, row.PhoneCountryCode, row.PhoneNumber, row.Email,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (r *PostgresUserRepository) FindByEmail(email domain.Email) (*domain.User, error) {
	dto := email.ToDTO()
	row, err := r.queries.FindUserByEmail(context.Background(), pgtype.Text{
		String: string(dto),
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToDomain(
		row.ID, row.IdentificationType, row.IdentificationNumber, row.FirstName, row.LastName,
		row.BirthDate, row.Address, row.PhoneCountryCode, row.PhoneNumber, row.Email,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (r *PostgresUserRepository) FindAll(pagination domain.Pagination) ([]*domain.User, error) {
	ctx := context.Background()
	limit := int32(pagination.Limit())
	var dbRows []sqlc.FindUsersPaginatedWithCursorRow

	if pagination.LastID() != nil {
		var err error
		dbRows, err = r.queries.FindUsersPaginatedWithCursor(ctx, sqlc.FindUsersPaginatedWithCursorParams{
			ID:    int64(*pagination.LastID()),
			Limit: limit,
		})
		if err != nil {
			return nil, err
		}
	} else {
		rows, err := r.queries.FindUsersPaginatedWithoutCursor(ctx, limit)
		if err != nil {
			return nil, err
		}
		dbRows = make([]sqlc.FindUsersPaginatedWithCursorRow, len(rows))
		for i, row := range rows {
			dbRows[i] = sqlc.FindUsersPaginatedWithCursorRow(row)
		}
	}

	users := make([]*domain.User, len(dbRows))
	for i, row := range dbRows {
		var err error
		users[i], err = mapToDomain(
			row.ID, row.IdentificationType, row.IdentificationNumber, row.FirstName, row.LastName,
			row.BirthDate, row.Address, row.PhoneCountryCode, row.PhoneNumber, row.Email,
			row.CreatedAt, row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	birthDate, err := toPgDate(user.BirthDate())
	if err != nil {
		return err
	}

	addressBytes, err := toJSONB(user.Address())
	if err != nil {
		return err
	}

	id, err := r.queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		IdentificationType:   string(user.Identification().ToDTO().Type),
		IdentificationNumber: user.Identification().ToDTO().Number,
		FirstName:            user.FirstName(),
		LastName:             user.LastName(),
		BirthDate:            birthDate,
		Address:              addressBytes,
		PhoneCountryCode:     user.Phone().ToDTO().CountryCode,
		PhoneNumber:          user.Phone().ToDTO().Number,
		Email:                toPgText(user.Email()),
		CreatedAt:            pgtype.Timestamptz{Time: user.CreatedAt(), Valid: true},
		UpdatedAt:            pgtype.Timestamptz{Time: user.UpdatedAt(), Valid: true},
	})
	if err != nil {
		return err
	}

	user.AssignID(int(id))
	return nil
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	birthDate, err := toPgDate(user.BirthDate())
	if err != nil {
		return err
	}

	addressBytes, err := toJSONB(user.Address())
	if err != nil {
		return err
	}

	return r.queries.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		IdentificationType:   string(user.Identification().ToDTO().Type),
		IdentificationNumber: user.Identification().ToDTO().Number,
		FirstName:            user.FirstName(),
		LastName:             user.LastName(),
		BirthDate:            birthDate,
		Address:              addressBytes,
		PhoneCountryCode:     user.Phone().ToDTO().CountryCode,
		PhoneNumber:          user.Phone().ToDTO().Number,
		Email:                toPgText(user.Email()),
		UpdatedAt:            pgtype.Timestamptz{Time: user.UpdatedAt(), Valid: true},
		ID:                   int64(user.ID()),
	})
}

func (r *PostgresUserRepository) Delete(id int) error {
	return r.queries.DeleteUser(context.Background(), int64(id))
}

func toPgDate(bd *domain.BirthDate) (pgtype.Date, error) {
	if bd == nil {
		return pgtype.Date{}, nil
	}
	t, err := time.Parse("2006-01-02", string(bd.ToDTO()))
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("error parsing birth date: %w", err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func toPgText(email *domain.Email) pgtype.Text {
	if email == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: string(email.ToDTO()), Valid: true}
}

func toJSONB(addr *domain.Address) ([]byte, error) {
	if addr == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(addr.ToDTO())
	if err != nil {
		return nil, fmt.Errorf("error marshaling address: %w", err)
	}
	return bytes, nil
}

func mapToDomain(
	id int64,
	idType, idNumber, firstName, lastName string,
	birthDate pgtype.Date,
	addressBytes []byte,
	phoneCountryCode, phoneNumber string,
	email pgtype.Text,
	createdAt, updatedAt pgtype.Timestamptz,
) (*domain.User, error) {
	identification, err := domain.NewIdentification(domain.IdentificationType(idType), idNumber)
	if err != nil {
		return nil, err
	}

	var address *domain.Address
	if len(addressBytes) > 0 {
		var dto domain.AddressDTO
		if err := json.Unmarshal(addressBytes, &dto); err != nil {
			return nil, fmt.Errorf("error deserializing address: %w", err)
		}
		addr, err := domain.NewAddress(
			dto.Street, dto.City, dto.State, dto.Country, dto.PostalCode, dto.Description,
		)
		if err != nil {
			return nil, err
		}
		address = &addr
	}

	var birthDateVO *domain.BirthDate
	if birthDate.Valid {
		bd, err := domain.NewBirthDate(birthDate.Time.Format("2006-01-02"))
		if err != nil {
			return nil, err
		}
		birthDateVO = &bd
	}

	personalInfo, err := domain.NewPersonalInformation(
		identification, firstName, lastName, address, birthDateVO,
	)
	if err != nil {
		return nil, err
	}

	phone, err := domain.NewPhone(phoneCountryCode, phoneNumber)
	if err != nil {
		return nil, err
	}

	var emailVO *domain.Email
	if email.Valid {
		e, err := domain.NewEmail(email.String)
		if err != nil {
			return nil, err
		}
		emailVO = &e
	}

	return domain.NewUser(domain.UserParams{
		ID:                  int(id),
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               emailVO,
		CreatedAt:           createdAt.Time,
		UpdatedAt:           updatedAt.Time,
	})
}

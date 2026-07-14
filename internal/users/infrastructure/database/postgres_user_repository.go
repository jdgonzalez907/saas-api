package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	queries *postgres.Queries
	pool    *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		queries: postgres.New(pool),
		pool:    pool,
	}
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.queries.FindUserByID(ctx, id)
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

func (r *PostgresUserRepository) FindByPhone(ctx context.Context, phone domain.Phone) (*domain.User, error) {
	row, err := r.queries.FindUserByPhone(ctx, postgres.FindUserByPhoneParams{
		PhoneCountryCode: phone.CountryCode(),
		PhoneNumber:      phone.Number(),
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

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	row, err := r.queries.FindUserByEmail(ctx, pgtype.Text{
		String: email.Value(),
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

func (r *PostgresUserRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]*domain.User, error) {
	limit := pagination.Limit()
	var dbRows []postgres.FindUsersPaginatedWithCursorRow

	if pagination.LastID() != nil {
		var err error
		dbRows, err = r.queries.FindUsersPaginatedWithCursor(ctx, postgres.FindUsersPaginatedWithCursorParams{
			ID:    *pagination.LastID(),
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
		dbRows = make([]postgres.FindUsersPaginatedWithCursorRow, len(rows))
		for i, row := range rows {
			dbRows[i] = postgres.FindUsersPaginatedWithCursorRow(row)
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

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	addressBytes, err := toJSONB(user.Address())
	if err != nil {
		return err
	}

	id, err := r.queries.CreateUser(ctx, postgres.CreateUserParams{
		IdentificationType:   string(user.Identification().Type()),
		IdentificationNumber: user.Identification().Number(),
		FirstName:            user.FirstName(),
		LastName:             user.LastName(),
		BirthDate:            toPgDate(user.BirthDate()),
		Address:              addressBytes,
		PhoneCountryCode:     user.Phone().CountryCode(),
		PhoneNumber:          user.Phone().Number(),
		Email:                toPgText(user.Email()),
		CreatedAt:            pgtype.Timestamptz{Time: user.CreatedAt(), Valid: true},
		UpdatedAt:            pgtype.Timestamptz{Time: user.UpdatedAt(), Valid: true},
	})
	if err != nil {
		return err
	}

	user.AssignID(id)
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	addressBytes, err := toJSONB(user.Address())
	if err != nil {
		return err
	}

	return r.queries.UpdateUser(ctx, postgres.UpdateUserParams{
		IdentificationType:   string(user.Identification().Type()),
		IdentificationNumber: user.Identification().Number(),
		FirstName:            user.FirstName(),
		LastName:             user.LastName(),
		BirthDate:            toPgDate(user.BirthDate()),
		Address:              addressBytes,
		PhoneCountryCode:     user.Phone().CountryCode(),
		PhoneNumber:          user.Phone().Number(),
		Email:                toPgText(user.Email()),
		UpdatedAt:            pgtype.Timestamptz{Time: user.UpdatedAt(), Valid: true},
		ID:                   user.ID(),
	})
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}

func toPgDate(bd *domain.BirthDate) pgtype.Date {
	if bd == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: bd.Time(), Valid: true}
}

func toPgText(email *domain.Email) pgtype.Text {
	if email == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: email.Value(), Valid: true}
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
		ID:                  id,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               emailVO,
		CreatedAt:           createdAt.Time.UTC(),
		UpdatedAt:           updatedAt.Time.UTC(),
	})
}

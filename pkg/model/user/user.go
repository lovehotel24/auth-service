package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/sys/database"
)

// Store manages the set of API's for user access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs a user store for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Authenticate finds a user by their phone and verifies their password. On
// success, it returns a UserId representing this user.
func (s Store) Authenticate(ctx context.Context, now time.Time, phone, password string) (string, error) {
	data := struct {
		Phone string `db:"phone"`
	}{
		Phone: phone,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		phone = :phone`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if err == database.ErrNotFound {
			return "", database.ErrNotFound
		}
		return "", fmt.Errorf("selecting user[%q]: %w", phone, err)
	}

	if err := bcrypt.CompareHashAndPassword(dbUsr.PasswordHash, []byte(password)); err != nil {
		return "", database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	//claims := auth.Claims{
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		Issuer:    "service project",
	//		Subject:   dbUsr.ID,
	//		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	//		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
	//	},
	//	Role: dbUsr.Role,
	//}

	return dbUsr.ID, nil
}

// Query retrieves a list of existing users from database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]User, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		users
	ORDER BY 
	    user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var dbUsrs []dbUser
	if err := database.NameQuerySlice(ctx, s.log, s.db, q, data, &dbUsrs); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return toUserSlice(dbUsrs), nil
}

// QueryByID gets the specified user from the database by ID.
func (s Store) QueryByID(ctx context.Context, userID string) (User, error) {
	if err := validate.CheckID(userID); err != nil {
		return User{}, database.ErrInvalidID
	}

	// If you are not an admin and looking to delete someone other than yourself.
	//if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
	//	return User{}, database.ErrForbidden
	//}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	SELECT
	    *
	FROM
		users
	WHERE 
	    user_id = :user_id`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("selecting userID[%s]: %w", userID, err)
	}

	//if dbUsr.Role != "ADMIN" {
	//	return User{}, database.ErrForbidden
	//}

	return toUser(dbUsr), nil
}

// QueryByPhone gets the specified user from database by phone number.
func (s Store) QueryByPhone(ctx context.Context, phone string) (User, error) {

	data := struct {
		Phone string `db:"phone"`
	}{
		Phone: phone,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
	    phone = :phone`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("selecting phone[%q]: %w", phone, err)
	}

	//if dbUsr.Role != "ADMIN" {
	//	return User{}, database.ErrForbidden
	//}

	return toUser(dbUsr), nil
}

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	dbUsr := dbUser{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Phone:        nu.Phone,
		PasswordHash: hash,
		Role:         nu.Role,
		DateCreated:  now,
		DateUpdated:  now,
	}

	const q = `
	INSERT INTO users
		(user_id, name, phone, password_hash, role, date_created, date_updated)
	VALUES 
	    (:user_id, :name, :phone, :password_hash, :role, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, dbUsr); err != nil {
		return User{}, fmt.Errorf("inserting user: %w", err)
	}

	return toUser(dbUsr), nil
}

// Update replaces a user document in the database.
func (s Store) Update(ctx context.Context, userID string, uu UpdateUser, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return database.ErrInvalidID
	}

	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	usr, err := s.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("updating user UserID[%s]: %w", userID, err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}

	if uu.Phone != nil {
		usr.Phone = *uu.Phone
	}

	if uu.Role != nil {
		usr.Role = *uu.Role
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	const q = `
	UPDATE
		users
	SET
		"name" = :name,
		"phone" = :phone,
		"role" = :role,
		"password_hash" = :password_hash,
		"date_updated" = :date_updated
	WHERE 
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		return fmt.Errorf("updating userID[%s]: %w", userID, err)
	}
	return nil
}

// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, userID string) error {
	if err := validate.CheckID(userID); err != nil {
		return database.ErrInvalidID
	}

	// If you are not an admin and looking to delete someone other than yourself.
	//if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
	//	return database.ErrForbidden
	//}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	DELETE FROM
		users
	WHERE 
	    user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", userID, err)
	}

	return nil
}

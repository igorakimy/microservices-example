package data

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbTimeout = time.Second * 3
)

var db *sqlx.DB

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in teh New function
type Models struct {
	User User
}

// New is the function used to create an instance of the data package. It returns the type
// Models, which embeds all the types we want to be available to our application.
func New(dbPool *sqlx.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name,omitempty" db:"first_name"`
	LastName  string    `json:"last_name,omitempty" db:"last_name"`
	Password  string    `json:"-" db:"password"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetAll returns all users.
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var users []*User

	if err := db.SelectContext(ctx, &users, getAllUsersQuery); err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}

	return users, nil
}

// GetByEmail returns one user by Email.
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	if err := db.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by ID
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	if err := db.GetContext(ctx, &user, getUserByIDQuery, id); err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	args := []any{
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	}
	if _, err := db.ExecContext(ctx, updateUserByIDQuery, args...); err != nil {
		log.Printf("Query failed: %v\n", err)
		return err
	}

	return nil
}

// Delete deletes one user from database, by User.ID.
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := db.ExecContext(ctx, deleteUserByIDQuery, u.ID); err != nil {
		log.Printf("Query failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteOne deletes user from database, by id.
func (u *User) DeleteOne(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := db.ExecContext(ctx, deleteUserByIDQuery, id); err != nil {
		log.Printf("Query failed: %v\n", err)
		return err
	}

	return nil
}

// Insert puts new user to the database and returns id of inserted user.
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Printf("Password hashing failed: %v\n", err)
		return 0, err
	}

	args := []any{
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		u.Active,
		time.Now(),
		time.Now(),
	}

	ret, err := db.ExecContext(ctx, insertNewUserQuery, args...)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return 0, err
	}

	insertedID, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(insertedID), nil
}

// ResetPassword changes user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("Password hashing failed: %v\n", err)
		return err
	}

	_, err = db.ExecContext(ctx, updateUserPasswordQuery, hashedPassword, u.ID)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

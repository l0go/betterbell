package internal

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type AuthState struct {
	DB       *sql.DB
	Sessions map[[16]byte]struct{}
}

// This is equivalent to the rows in the Users table
type User struct {
	ID       int
	Username string
	// The hash and salt are in base64 format
	Hash string
	Salt string
}

// Status for the login
type LoginStatus uint8

const (
	LoginSuccess LoginStatus = iota
	UsernameInvalid
	PasswordInvalid
	UnknownLoginError
)

// Status for the registration
type RegisterStatus uint8

const (
	RegisterSuccess RegisterStatus = iota
	UsernameTaken
	UnknownRegisterError
)

// Checks if the username and password are valid
func (a AuthState) Check(username, password string) (LoginStatus, error) {
	// Return if the user does not exist
	exists, err := a.Exists(username)
	if err != nil || !exists {
		return UsernameInvalid, err
	}

	// Else, check the password to make sure it is valid
	rows, err := a.DB.Query("SELECT * FROM Users WHERE username=?;", strings.ToLower(username))
	if err != nil {
		return UnknownLoginError, err
	}
	defer rows.Close()

	// Get the user
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.Hash, &user.Salt); err != nil {
			return UnknownLoginError, err
		}
		break
	}

	// Hash with the same parameters that the original password used on the password we are testing
	decodedSalt, _ := base64.StdEncoding.DecodeString(user.Salt)
	checkHash := base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), decodedSalt, 1, 64*1024, 1, 32))

	var status LoginStatus
	if checkHash == user.Hash {
		status = LoginSuccess
	} else {
		status = PasswordInvalid
	}
	return status, nil
}

// Registers an user with the provided username and password
func (a AuthState) Register(username, password string) (RegisterStatus, error) {
	// Check if user already exists
	if exists, err := a.Exists(username); exists || err != nil {
		return UsernameTaken, err
	}

	// Hash and salt password
	hash, salt, err := hash(password)
	if err != nil {
		return UnknownRegisterError, err
	}

	encodedHash := base64.StdEncoding.EncodeToString(hash)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	_, err = a.DB.Exec("INSERT INTO Users (username, hash, salt) VALUES(?, ?, ?);", strings.ToLower(username), encodedHash, encodedSalt)
	if err != nil {
		return UnknownRegisterError, err
	}

	return RegisterSuccess, nil
}

// Returns if the username exists
// If for some reason there is an error, it will return true along with the error
func (a AuthState) Exists(username string) (bool, error) {
	rows, err := a.DB.Query("SELECT * FROM Users WHERE username=?;", strings.ToLower(username))
	if err != nil {
		return true, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

// Based on the status, return a human-readable string
func FormatLoginStatus(status LoginStatus) string {
	switch status {
	case LoginSuccess:
		return "Login was successful!"
	case UsernameInvalid:
		fallthrough
	case PasswordInvalid:
		return "Username and/or password is invalid"
	case UnknownLoginError:
		return "Unknown error: Please try again!"
	}
	return ""
}

// Based on the status, return a human-readable string
func FormatRegisterStatus(status RegisterStatus) string {
	switch status {
	case RegisterSuccess:
		return "Registration was successful!"
	case UsernameTaken:
		return "Username is taken"
	case UnknownRegisterError:
		return "Unknown error: Please try again!"
	}
	return ""
}

// Adds a session key
func (a AuthState) GrantSession(r *http.Request) (http.Cookie, error) {
	var sessionKey, err = generateKey(16)

	// Check if the sessionKey exists, if so we will regenerate another one
	// The probability of this happening is VERY low, but better to be safe
	for {
		if _, ok := a.Sessions[[16]byte(sessionKey)]; ok {
			sessionKey, err = generateKey(16)
			continue
		}
		break
	}

	// Add sessionKey to sessions array
	a.Sessions[[16]byte(sessionKey)] = struct{}{}

	sessionKeyHash := base64.StdEncoding.EncodeToString(sessionKey)

	// Set cookie
	cookie := http.Cookie{
		Name:     "session_key",
		Value:    sessionKeyHash,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		Secure:   false, // TODO: Set this to true
		HttpOnly: true,
		Path:     "/",
	}

	return cookie, err
}

// Checks if the session key stored in the sessionKey has permission to access that page
func (a AuthState) HasPermission(r *http.Request) (bool, error) {
	sessionKey, err := r.Cookie("session_key")
	if sessionKey == nil {
		return false, nil
	}
	decodedKey, err2 := base64.StdEncoding.DecodeString(sessionKey.Value)
	_, ok := a.Sessions[[16]byte(decodedKey)]
	return ok, errors.Join(err, err2)
}

// From a password string will provide a hash and salt
func hash(password string) (hash, salt []byte, err error) {
	// Generate a salt
	salt, err = generateKey(16)
	if err != nil {
		return nil, nil, err
	}

	// Generate the hash
	hash = argon2.IDKey([]byte(password), salt, 1, 64*1024, 1, 32)
	return hash, salt, nil
}

// Generates a cryptographically secure random key, used for the salt and session key
func generateKey(length int) ([]byte, error) {
	key := make([]byte, length) // Create an array that is 16 characters long
	_, err := rand.Read(key)
	return key, err
}

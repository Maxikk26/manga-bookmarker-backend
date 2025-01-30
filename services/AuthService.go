package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/models"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
	"time"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

//Core services

func Login(login dtos.Login) (tokenString string, err error) {
	filter := bson.M{"username": login.Username}
	user, code, err := repository.FindUser(filter)
	if err != nil {
		fmt.Println("Error obtaining user: ", err.Error())
		return "", errors.New("User or password invalid")
	}

	if code == constants.NoDocumentFound {
		return "", errors.New("User or password invalid")
	}

	if user.Attempts == 3 {
		return "", errors.New("User blocked")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		log.Println("err", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			filter = bson.M{"_id": user.Id}
			update := bson.D{
				{"$set", bson.M{
					"attempts": user.Attempts + 1,
				}},
			}

			_, err = repository.UpdateUser(filter, update)
			if err != nil {
				log.Println("err", err)
			}

			return "", errors.New("User or password invalid")
		}
		return "", errors.New("Error interno")
	}

	//Reset attempts
	go func() {
		if user.Attempts > 0 {
			filter = bson.M{"_id": user.Id}
			update := bson.D{
				{"$set", bson.M{
					"attempts": 0,
				}},
			}

			_, _ = repository.UpdateUser(filter, update)
		}
	}()

	var jwtKey = []byte("my_secret_key")

	// Set the JWT claims
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		UserId: user.Id.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing token: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func GetAllUsers() (users []models.User, err error) {
	users, err = repository.GetUsers()
	if err != nil {
		log.Fatal("GetAllUsers Error: ", err)
		return nil, err
	}
	return users, nil
}

func CreateUser(user dtos.UserCreate) (err error) {
	//TODO user, email and password validation

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}
	fmt.Println(hashedPassword)

	user.Password = hashedPassword

	var userModel models.User
	// Use dto-mapper to map the data to the struct
	err = utils.Mapper.Map(&userModel, &user)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return err
	}

	userModel.Rol = "master"
	userModel.Status = 1
	userModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	err = repository.CreateUser(userModel)
	if err != nil {
		return err
	}
	return nil
}

func GetUserIdFromClaims(tokenString string) (userId string, err error, errorCode int) {
	// Secret key used for signing the JWT (replace with your actual key)
	secretKey := []byte("my_secret_key")

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		// Check if the error is due to the token being expired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", fmt.Errorf("token is expired"), 1
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			return "", fmt.Errorf("token is malformed"), 1
		}

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", fmt.Errorf("token is invalid"), 1
		}
		// Handle other errors
		return "", err, 1
	}

	// Validate the token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserId, nil, 0
	} else {
		return "", err, 2
	}
}

func ValidUser(userId string) (valid bool, err error) {
	if userId == "" {
		return false, errors.New("invalid user")
	}

	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false, errors.New("invalid user")
	}

	// Define conditions for finding the user's bookmarks
	filter := bson.M{"_id": objectID}

	_, _, err = repository.FindUser(filter)
	if err != nil {
		return false, errors.New("invalid user")
	}

	return true, nil
}

// Helpers
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

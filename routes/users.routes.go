package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ItsDavidHdez/go-sql-restapi/db"
	"github.com/ItsDavidHdez/go-sql-restapi/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const signKey = "changeMeInProduction"

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Credentials struct {
	email string `json:"email"`
	Password string `json:"password"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.Write([]byte("User not found"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&user)

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Unscoped().Delete(&user)
	w.WriteHeader(http.StatusOK)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	buf, errBuf := io.ReadAll(r.Body)
    if errBuf != nil {
        panic(errBuf)
    }

	json.Unmarshal(buf, &body)

	hash, errorHash := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if errorHash != nil {
		w.Write([]byte(errorHash.Error()))
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}
	json.NewDecoder(r.Body).Decode(&body)

	
	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} 

	json.NewEncoder(w).Encode(&user)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	buf, errBuf := io.ReadAll(r.Body)
    if errBuf != nil {
        panic(errBuf)
    }

	json.Unmarshal(buf, &body)

	var user models.User
	db.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		w.Write([]byte("Without user"))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	// Generate a new jwt
	var credentials Credentials
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Email: credentials.email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errTokenString := token.SignedString([]byte(os.Getenv("SECRET")))

	if errTokenString != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	json.NewEncoder(w).Encode(&tokenString)

}

func Test(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome, you are logged")))
}

package service
import (
    "errors"
    "os"
    "time"

    "github.com/ankan8/swapsync/backend/services/auth-service/models"
    "github.com/ankan8/swapsync/backend/services/auth-service/repository"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)



func RegisterUser(email, password string) (string, error) {
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    newUser := models.User{
        Email:    email,
        Password: string(hashedPassword),
    }
    
    newUser.ID = uuid.NewString()
    if err := repository.InsertUser(&newUser); err != nil {
        return "", err
    }
    return newUser.ID, nil
}

func AuthenticateUser(email, password string) (string, error) {
    user, err := repository.FindUserByEmail(email)
    if err != nil {
        return "", errors.New("user not found")
    }
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return "", errors.New("invalid password")
    }
    token, err := GenerateJWT(email)
    if err != nil {
        return "", err
    }
    return token, nil
}

func GenerateJWT(email string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET not set")
    }
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

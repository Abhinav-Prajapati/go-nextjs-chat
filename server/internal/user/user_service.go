package user

import (
	"context"
	"server/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TODO : move key to yaml file
const (
	jwtSecretKey = "mykey3984jshalfjslfdjaslkdjflaskjdfalksjdfklj"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}
func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	user, err := s.Repository.CreateUser(ctx, u)

	if err != nil {
		return nil, err
	}

	response := &CreateUserRes{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		Email:    user.Email,
	}
	return response, nil
}

type jwtClame struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}
	err = utils.CheakPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClame{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signedString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{
		accessToken: signedString,
		ID:          strconv.Itoa(int(u.ID)),
		Username:    u.Username,
	}, err
}

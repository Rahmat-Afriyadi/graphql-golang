package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"product-golang-graphql/configs"
	"product-golang-graphql/graph/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


type Resolver struct{}


var SecretKey string = configs.EnvJWTSecret()


var db = configs.Connect()

func GenerateToken(id string, name string, email string) (string, error) {
	claims := jwt.MapClaims{
		"id": id,
		"name": name,
		"email": email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string) (*model.AuthPayload, error) {
	// Koneksi ke MongoDB

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := db.CreateUser(model.User{Email: email, Name: name, Password: string(hashedPassword)})
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, err
	}


	// Kembalikan token dan user
	return &model.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, name *string, email *string) (*model.User, error) {
	return db.UpdateUser(model.User{ID: id, Name: *name, Email: *email})
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return db.DeleteUser(id)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	user:=db.GetUser(email)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("ini error ", err)
		return &model.AuthPayload{}, errors.New("Password salah")
	}
	token, err := GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, err
	}


	return &model.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value("user").(*model.User)
	if !ok {
		return nil, errors.New("unauthenticated")
	}
	return user, nil
}

func (r *queryResolver) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	return db.GetAllProducts()
}

func (r *mutationResolver) CreateProduct(ctx context.Context, name string, price float64, stock int32) (*model.Product, error) {
	user, _ := ctx.Value("user").(*model.User)
	return db.CreateProduct(model.Product{Name: name, Price: price, Stock: stock, CreatedBy: user.Name})
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, name *string, price *float64, stock *int32) (*model.Product, error) {
	user, _ := ctx.Value("user").(*model.User)
	
	return db.UpdateProduct(model.Product{ID: id, Name: *name, Price: *price, Stock: *stock, UpdatedBy: user.Name  })
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	return db.DeleteProduct(id)
}

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return db.GetAllUsers()
}
func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	return db.GetUserById(id)
}

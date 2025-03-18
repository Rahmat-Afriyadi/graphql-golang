package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.68

import (
	"context"
	"product-golang-graphql/graph/model"
)

// Register is the resolver for the register field.
// func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string) (*model.AuthPayload, error) {
// 	// Koneksi ke MongoDB

// 	// Hash password sebelum disimpan
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := db.CreateUser(model.User{Email: email, Name: name, Password: string(hashedPassword)})

// 	// Generate JWT token
// 	token, err := GenerateToken(user.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Kembalikan token dan user
// 	return &model.AuthPayload{
// 		Token: token,
// 		User:  user,
// 	}, nil
// }

// Login is the resolver for the login field.
// func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
// 	user := db.GetUser(email)

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		fmt.Println("ini error ", err)
// 		return &model.AuthPayload{}, errors.New("Password salah")
// 	}
// 	if users[email] != password {
// 		return nil, errors.New("invalid credentials")
// 	}

// 	token, err := GenerateToken(email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.AuthPayload{
// 		Token: token,
// 		User:  user,
// 	}, nil
// }

// UpdateUser is the resolver for the updateUser field.
// func (r *mutationResolver) UpdateUser(ctx context.Context, id string, name *string, email *string) (*model.User, error) {
// 	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
// }

// DeleteUser is the resolver for the deleteUser field.
// func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
// 	panic(fmt.Errorf("not implemented: DeleteUser - deleteUser"))
// }

// CreateProduct is the resolver for the createProduct field.
// func (r *mutationResolver) CreateProduct(ctx context.Context, name string, price float64, stock int32) (*model.Product, error) {
// 	panic(fmt.Errorf("not implemented: CreateProduct - createProduct"))
// }

// UpdateProduct is the resolver for the updateProduct field.
// func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, name *string, price *float64, stock *int32) (*model.Product, error) {
// 	panic(fmt.Errorf("not implemented: UpdateProduct - updateProduct"))
// }

// DeleteProduct is the resolver for the deleteProduct field.
// func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
// 	panic(fmt.Errorf("not implemented: DeleteProduct - deleteProduct"))
// }

// GetUser is the resolver for the getUser field.
// func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
// 	panic(fmt.Errorf("not implemented: GetUser - getUser"))
// }

// GetAllUsers is the resolver for the getAllUsers field.
// func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
// 	panic(fmt.Errorf("not implemented: GetAllUsers - getAllUsers"))
// }

// Me is the resolver for the me field.
// func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
// 	user, ok := ctx.Value("user").(*model.User)
// 	if !ok {
// 		return nil, errors.New("unauthenticated")
// 	}
// 	return user, nil
// }

// GetProduct is the resolver for the getProduct field.
func (r *queryResolver) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	return db.GetProductById(id)
}

// GetAllProducts is the resolver for the getAllProducts field.
// func (r *queryResolver) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
// 	panic(fmt.Errorf("not implemented: GetAllProducts - getAllProducts"))
// }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

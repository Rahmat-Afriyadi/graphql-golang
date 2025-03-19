package configs

import (
	"context"
	"fmt"
	"log"
	"product-golang-graphql/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var connectionString string = EnvMongoURI()

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	createUniqueIndex(client, "user", "email")

	return &DB{
		client: client,
	}
}

func createUniqueIndex(Client *mongo.Client,collectionName, fieldName string) {
	collection := Client.Database("product_db").Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{fieldName: 1}, // 1 for ascending order
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Could not create index for %s: %v", fieldName, err)
	}

}


func (db *DB) CreateUser(user model.User) (*model.User, error) {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := userCollec.InsertOne(ctx, bson.M{"name": user.Name, "email": user.Email,"password":user.Password})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("email %s already exists", user.Email)
		}
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
	}
	
	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.User{ID: insertedID, Email: user.Email, Name: user.Name}
	return &returnJobListing, nil
}

func (db *DB) CreateProduct(product model.Product) (*model.Product, error) {
	productCollec := db.client.Database("product_db").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := productCollec.InsertOne(ctx, bson.M{"name": product.Name, "price": product.Price,"stock":product.Stock,"created_by":product.CreatedBy})
	if err != nil {
		log.Fatal(err)
	}
	
	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnProduct := model.Product{ID: insertedID, Name: product.Name, Price: product.Price, Stock: product.Stock}
	return &returnProduct, nil
}

func (db *DB) UpdateUser(user model.User) (*model.User, error) {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"name": user.Name, "email":user.Email}}
	
	opts := options.Update().SetUpsert(false) 
	
	result, err := userCollec.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no document found with ID: %s", user.ID)
	}
	
	// upsertedID := result.UpsertedID.(primitive.ObjectID).Hex()
	returnUser := model.User{ID: user.ID, Name: user.Name, Email: user.Email}
	return &returnUser, nil
}

func (db *DB) DeleteUser(id string) (bool, error) {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid ID format")
	}
	
	filter := bson.M{"_id": objectID}
	
	result, err := userCollec.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, fmt.Errorf("no document found with ID: %s", id)
	}
	
	return true, nil
}

func (db *DB) UpdateProduct(product model.Product) (*model.Product, error) {
	productCollec := db.client.Database("product_db").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	objectID, err := primitive.ObjectIDFromHex(product.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"name": product.Name, "price":product.Price,"stock":product.Stock, "updated_by":product.UpdatedBy}}
	
	opts := options.Update().SetUpsert(false) 
	
	result, err := productCollec.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no document found with ID: %s", product.ID)
	}
	
	// upsertedID := result.UpsertedID.(primitive.ObjectID).Hex()
	returnProduct := model.Product{ID: product.ID, Name: product.Name, Price: product.Price, Stock: product.Stock}
	return &returnProduct, nil
}

func (db *DB) DeleteProduct(id string) (bool, error) {
	productCollec := db.client.Database("product_db").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid ID format")
	}
	
	filter := bson.M{"_id": objectID}
	
	result, err := productCollec.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, fmt.Errorf("no document found with ID: %s", id)
	}
	
	return true, nil
}

func (db *DB) GetUser(email string) *model.User {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// _id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"email": email}
	var userListing model.User
	err := userCollec.FindOne(ctx, filter).Decode(&userListing)
	if err != nil {
		log.Fatal(err)
	}
	return &userListing
}

func (db *DB) GetAllProducts(search string) ([]*model.Product, error) {
	productCollec := db.client.Database("product_db").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"name": bson.M{"$regex": search, "$options": "i"}}

	cursor, err := productCollec.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Store results in a slice
	var products []*model.Product
	for cursor.Next(ctx) {
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		products = append(products, &product)
	}
	return products, nil
}

func (db *DB) GetProductById(id string) (*model.Product, error) {
	productCollec := db.client.Database("product_db").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var productListing model.Product
	err := productCollec.FindOne(ctx, filter).Decode(&productListing)
	if err != nil {
		return nil, err
	}
	return &productListing, nil
}

func (db *DB) GetUserById(id string) (*model.User, error) {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var userListing model.User
	err := userCollec.FindOne(ctx, filter).Decode(&userListing)
	if err != nil {
		return nil, err
	}
	return &userListing, nil
}

func (db *DB) GetAllUsers() ([]*model.User, error) {
	userCollec := db.client.Database("product_db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := userCollec.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Store results in a slice
	var users []*model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, &user)
	}
	return users, nil
}


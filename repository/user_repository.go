package repository

import (
	"context"
	"loan-tracker-api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepositoryImpl(coll *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{collection: coll}
}


func (u *UserRepositoryImpl) Register(user domain.User) error{
	_, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}



func (u *UserRepositoryImpl) GetUserByUsernameOrEmail(username, email string) (domain.User, error) {
	var user domain.User
	err := u.collection.FindOne(context.Background(), map[string]string{"username": username, "email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}


func (u *UserRepositoryImpl) AccountActivation(email string) error {
	
	
	_, err := u.collection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"is_active": true}, "$unset": bson.M{"activation_token": ""}, "$currentDate": bson.M{"updated_at": true}})
	if err != nil {
		return err
	}


	return nil

}

func (u *UserRepositoryImpl) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	err := u.collection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UserRepositoryImpl) UpdateUser(user *domain.User) error {
	_, err := u.collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": user})
	if err != nil {
		return err
	}

	return nil
}


func (ur *UserRepositoryImpl) Login(user *domain.User) (*domain.User, error) {
	var existingUser domain.User
	err := ur.collection.FindOne(context.Background(), map[string]string{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		return &domain.User{}, err
	}
	return &existingUser, nil
	
}




func (ur *UserRepositoryImpl) GetUserByID(id string) (domain.User, error) {
	var user domain.User
	objID, err:= primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	err = ur.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}



func (ur *UserRepositoryImpl) DeleteRefreshToken(user *domain.User, token string) error {
	objID, err := primitive.ObjectIDFromHex(user.ID.Hex())
    if err != nil {
        return err
    }
    _, err = ur.collection.UpdateOne(
        context.Background(),
        bson.M{"_id": objID},
        bson.M{"$pull": bson.M{"refresh_tokens": bson.M{"token": token}}},
    )
    return err
}


func (ur *UserRepositoryImpl) DeleteAllRefreshTokens(user *domain.User) error {
	_, err := ur.collection.UpdateOne(context.Background(), map[string]string{"username": user.Username}, bson.M{"$set": bson.M{"refresh_tokens": []domain.RefreshToken{}}})
	return err
}



func (ur *UserRepositoryImpl) GetMyProfile(userID string) (domain.User, error) {
	var user domain.User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, err
	}

	err = ur.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}


func (ur *UserRepositoryImpl) GetUsers()([]domain.User, error){
	var users []domain.User
	cursor, err := ur.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []domain.User{}, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user domain.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}


func (ur *UserRepositoryImpl) DeleteUser(id string) (domain.User, error) {
	var user domain.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	err = ur.collection.FindOneAndDelete(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}


func (ur *UserRepositoryImpl) GetUserByResetToken(token string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"password_reset_token": token}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}



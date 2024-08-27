package repository

import (
	"context"
	"loan-tracker-api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLoanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepositoryImpl(LoanCollection *mongo.Collection) domain.LoanRepository {
	return &MongoLoanRepository{
		collection:        LoanCollection,
		
	}
}

func (m *MongoLoanRepository) CreateLoan(loan domain.Loan) (domain.Loan, error) {
	var newLoan domain.Loan
	new_, err := m.collection.InsertOne(context.Background(), loan)
	if err != nil {
		return domain.Loan{}, err
	}
	err = m.collection.FindOne(context.Background(), bson.M{"_id": new_.InsertedID}).Decode(&newLoan)

	return newLoan, nil
	
}


func (m *MongoLoanRepository) GetLoanByID(id string) (domain.Loan, error) {
	var loan domain.Loan
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Loan{}, err
	}
	err = m.collection.FindOne(context.Background(),bson.M{"_id":objID}).Decode(&loan)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

func (m *MongoLoanRepository) GetLoans() ([]domain.Loan, error) {
	var loans []domain.Loan
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var loan domain.Loan
		cursor.Decode(&loan)
		loans = append(loans, loan)
	}
	return loans, nil
}

func (m *MongoLoanRepository) UpdateLoanStatus(newStatus, loanID string) (domain.Loan, error) {
	var loan domain.Loan
	objID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return domain.Loan{}, err
	}
	err = m.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": newStatus}}).Decode(&loan)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}


func (m *MongoLoanRepository) DeleteLoan(id string) (domain.Loan, error) {
	var loan domain.Loan
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Loan{}, err
	}
	err = m.collection.FindOneAndDelete(context.Background(), bson.M{"_id": objID}).Decode(&loan)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}
package service

import (
	"context"
	"time"

	"user-service/models"
	"user-service/pb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	UserCollection *mongo.Collection
}

func NewUserService(userCollection *mongo.Collection) *UserServiceServer {
	return &UserServiceServer{
		UserCollection: userCollection,
	}
}

func (s *UserServiceServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	idHex := req.GetId()

	// แปลง string -> ObjectID
	objectID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	// ค้นหาผู้ใช้
	var user models.User
	err = s.UserCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileResponse{
		Id:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

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
func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	now := time.Now()

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"email":      req.GetEmail(),
			"updated_at": now,
		},
	}

	_, err = s.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var updatedUser models.User
	err = s.UserCollection.FindOne(ctx, filter).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{
		Id:        updatedUser.ID.Hex(),
		Email:     updatedUser.Email,
		UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (s *UserServiceServer) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	now := time.Now()

	filter := bson.M{"_id": userID, "deleted_at": bson.M{"$exists": false}} // ตรวจสอบว่าไม่ถูกลบแล้ว
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
		},
	}

	result, err := s.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return &pb.DeleteProfileResponse{
			Success: false,
			Message: "User not found or already deleted",
		}, nil
	}

	return &pb.DeleteProfileResponse{
		Success: true,
		Message: "User deleted successfully (soft delete)",
	}, nil
}

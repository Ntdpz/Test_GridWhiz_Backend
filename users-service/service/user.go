package service

import (
	"context"
	"time"
	"user-service/models"
	"user-service/pb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// Apply filters
	if req.Email != "" {
		filter["email"] = bson.M{"$regex": req.Email, "$options": "i"} // case-insensitive match
	}
	if req.Username != "" {
		filter["username"] = bson.M{"$regex": req.Username, "$options": "i"}
	}

	// Pagination
	page := req.GetPage()
	limit := req.GetLimit()
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	skip := int64((page - 1) * limit)

	// Query
	cursor, err := s.UserCollection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &[]int64{int64(limit)}[0],
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*pb.UserInfo
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		users = append(users, &pb.UserInfo{
			Id:        user.ID.Hex(),
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	// Total count
	count, _ := s.UserCollection.CountDocuments(ctx, filter)

	return &pb.ListUsersResponse{
		Users: users,
		Total: int32(count),
	}, nil
}

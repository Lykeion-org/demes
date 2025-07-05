package grpc

import (
	"github.com/google/uuid"
	pb "github.com/Lykeion-org/go-shared/pkg/grpc/generated/demes"
	models "github.com/Lykeion-org/demes/internal/models"
)


func ConvertUserToProto(user *models.LykeionUser) *pb.User {
	
	return &pb.User{
		Uid:        		user.UID.String(),
		EmailAddress: 		user.EmailAddress,
		PasswordHash: 		user.PasswordHash,
		EmailAddressHash: 	user.EmailAddressHash,
		AvatarSource: 		user.AvatarSource,
		NativeLanguage: 	user.NativeLanguage,
		ReferenceLanguage: 	user.ReferenceLanguage,
		TargetLanguage: 	user.TargetLanguage,
		FirstName: 			user.FirstName,
		LastName: 			user.LastName,
		ExtraName: 			user.ExtraName,
		Role:   			user.Role,
		Verified: 			user.Verified,
		}
	}


func ConvertProtoToUser(user *pb.User) (*models.LykeionUser,error) {
	userUID, err := uuid.Parse(user.Uid)
	if err != nil {
		return nil, err
	}

	return &models.LykeionUser{
		UID: userUID,
		EmailAddress: user.EmailAddress,
		PasswordHash: user.PasswordHash,
		EmailAddressHash: user.EmailAddressHash,
		AvatarSource: user.AvatarSource,
		NativeLanguage: user.NativeLanguage,
		ReferenceLanguage: user.ReferenceLanguage,
		TargetLanguage: user.TargetLanguage,
		FirstName: strPtr(user.GetFirstName()),
		LastName: strPtr(user.GetLastName()),
		ExtraName: strPtr(user.GetExtraName()),
		Role: user.Role,
		Verified: user.Verified,
	}, nil
}

func ConvertUserListToProto(users []*models.LykeionUser) []*pb.User {
	protoUsers := make([]*pb.User, 0, len(users))
	for _, u := range users {
		protoUsers = append(protoUsers, ConvertUserToProto(u))
	}
	return protoUsers
}

func ConvertProtoToUsers(protoUsers []*pb.User) []*models.LykeionUser {
	users := make([]*models.LykeionUser, 0, len(protoUsers))
	for _, pu := range protoUsers {
		user, _ := ConvertProtoToUser(pu)
		if user != nil {
			users = append(users, user)
		}
	}
	return users
}



func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int32Ptr(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}
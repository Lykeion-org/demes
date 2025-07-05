package models

import(
	"time"
	"github.com/google/uuid"
)

type LykeionUser struct {
	UID					uuid.UUID		`gorm:"type:uuid;index;default:uuid_generate_v4();primaryKey" json:"uid"`			
	EmailAddress		string			`gorm:"not null" json:"emailAddress"`		
	PasswordHash		string			`gorm:"not null"`		
	EmailAddressHash	string			`gorm:"not null"`		
	AvatarSource		string			`gorm:"not null" json:"avatarSource"`		
	NativeLanguage		int32			`gorm:"not null" json:"nativeLanguage"`	
	ReferenceLanguage	int32			`gorm:"not null" json:"referenceLanguage"`	
	TargetLanguage		int32			`gorm:"not null" json:"targetLanguage"`	
	FirstName			*string			`gorm:"" json:"firstName"`		
	LastName			*string			`gorm:"" json:"lastName"`		
	ExtraName			*string			`gorm:"" json:"extraName"`		
	DisplayName			*string			`gorm:"" json:"displayName"`		
	Role				int32			`gorm:"not null" json:"role"`	
	Verified			bool			`gorm:"not null" json:"isVerfied"`		
	
	CreatedAt			time.Time
	UpdatedAt			time.Time		

}
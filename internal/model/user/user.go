package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                 string               `bson:"_id,omitempty"`
	Username           string               `bson:"username"`
	Password           string               `bson:"password"`
	Email              string               `bson:"email"`
	RefreshToken       []string             `bson:"refreshToken,omitempty"`
	AccessToken        string               `bson:"accessToken,omitempty"`
	VerificationToken  string               `bson:"verificationToken,omitempty"`
	PasswordResetToken string               `bson:"passwordResetToken,omitempty"`
	RecoveryToken      string               `bson:"recoveryToken,omitempty"`
	Verified           bool                 `bson:"verified"`
	Roles              []string             `bson:"roles"`
	ImageURL           string               `bson:"profileImageURL"`
	GoogleID           string               `bson:"googleId"`
	UsedStorage        primitive.Decimal128 `bson:"usedStorage"`
	Created            time.Time            `bson:"created_at"`
}

package user

type Model struct {
	ID                 string   `bson:"_id,omitempty"`
	Username           string   `bson:"username"`
	Password           string   `bson:"password"`
	Email              string   `bson:"email"`
	RefreshToken       []string `bson:"refreshToken"`
	VerificationToken  string   `bson:"verificationToken"`
	PasswordResetToken string   `bson:"passwordResetToken"`
	RecoveryToken      string   `bson:"recoveryToken"`
	Verified           bool     `bson:"verified"`
	Roles              []string `bson:"roles"`
	ImageURL           string   `bson:"profileImageURL"`
	GoogleID           string   `bson:"googleId"`
	UsedStorage        int32    `bson:"usedStorage"`
}

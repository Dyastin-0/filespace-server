package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	usr "filespace/internal/model/user"
)

func Signout(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("rt")

		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		token := cookie.Value

		claims := types.Claims{}
		jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
		})

		collection := client.Database("test").Collection("users")
		filter := bson.M{"email": claims.User.Email}
		user := usr.Model{}
		err = collection.FindOne(r.Context(), filter).Decode(&user)

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "rt",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
				MaxAge:   -1,
				Domain:   os.Getenv("DOMAIN"),
				Path:     "/",
			})

			w.WriteHeader(http.StatusOK)
		}

		newRefeshToken := user.RefreshToken

		for i, rt := range user.RefreshToken {
			if rt == token {
				newRefeshToken = append(user.RefreshToken[:i], user.RefreshToken[i+1:]...)
				break
			}
		}

		collection.UpdateOne(r.Context(), filter, bson.M{"$set": bson.M{"refreshToken": newRefeshToken}})

		http.SetCookie(w, &http.Cookie{
			Name:     "rt",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   -1,
			Domain:   os.Getenv("DOMAIN"),
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
	}
}

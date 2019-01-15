package basicserver

import (
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
)

// RequireAuth is a middleware used by routes which require authentication.
//
// It checks request `Authorization` header and tries to parse it.
//
// If everything goes well with parsing, then the "uid" value is passed to Next().
//
// In case of invalid/expired token, this returns status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) RequireAuth() iris.Handler {
	return func(ctx iris.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			app.HandleError(errors.New("No Authorization Header"), ctx, iris.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return app.Settings.Secret, nil
		})
		if err != nil {
			app.HandleError(err, ctx, iris.StatusUnauthorized)
			return
		}

		uid := token.Claims.(jwt.MapClaims)["uid"]
		if uid == nil {
			err := errors.New("Incorrect Authorization Header")
			app.HandleError(err, ctx, iris.StatusUnauthorized)
			return
		}

		objectUID := bson.ObjectIdHex(uid.(string))
		var user User
		err = app.Coll.Users.FindId(objectUID).One(&user)
		if err != nil {
			if err.Error() == "not found" {
				err := errors.New("No Such User")
				app.HandleError(err, ctx, iris.StatusUnauthorized)
				ctx.WriteString("No Such User")
			} else {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
			return
		}

		// pass on the "uid"
		ctx.Values().Set("uid", uid)
		ctx.Next()
	}
}

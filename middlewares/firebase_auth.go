package middlewares

import (
	"net/http"

	"firebase.google.com/go/v4/appcheck"
	"github.com/gin-gonic/gin"
)

const (
	FirebaseAppCheckToken = "X-Firebase-AppCheck"
)

type firebaseAuthMiddleware struct {
	appCheck *appcheck.Client
}

type IFirebaseAuthMiddleware interface {
	RequireAppCheck() gin.HandlerFunc
}

var _ IFirebaseAuthMiddleware = &firebaseAuthMiddleware{}

// RequireAppCheck refer to https://firebase.google.com/docs/app-check/custom-resource-backend#go. Check if request has valid app check token
func (f firebaseAuthMiddleware) RequireAppCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		w := c.Writer

		appCheckToken, ok := r.Header[http.CanonicalHeaderKey(FirebaseAppCheckToken)]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized."))
			return
		}

		_, err := f.appCheck.VerifyToken(appCheckToken[0])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized."))
			return
		}

		// before request
		c.Next()
	}
}

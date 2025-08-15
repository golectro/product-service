package interceptor

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"golectro-product/internal/model"
	"golectro-product/internal/utils"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(viper *viper.Viper) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

		jwksURL := viper.GetString("KEYCLOAK_URL") + "/realms/golectro/protocol/openid-connect/certs"
		jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{RefreshInterval: time.Hour})
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to load JWKS")
		}

		parsedToken, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !parsedToken.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "failed to parse token claims")
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "user ID (sub) not found in token claims")
		}

		uid, err := utils.ParseUUID(userID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid user ID")
		}

		var roles []string
		if realmAccess, ok := claims["realm_access"].(map[string]any); ok {
			if roleList, ok := realmAccess["roles"].([]any); ok {
				roles = make([]string, 0, len(roleList))
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}

		rolesJSON, err := json.Marshal(roles)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to marshal roles")
		}

		username, _ := claims["preferred_username"].(string)
		email, _ := claims["email"].(string)

		auth := &model.Auth{
			ID:       uid,
			Username: username,
			Email:    email,
			Roles:    rolesJSON,
		}

		ctx = context.WithValue(ctx, UserContextKey, auth)

		return handler(ctx, req)
	}
}

func GetUserFromContext(ctx context.Context) *model.Auth {
	if auth, ok := ctx.Value(UserContextKey).(*model.Auth); ok {
		return auth
	}
	return nil
}

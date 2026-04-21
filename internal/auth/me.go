package auth

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	Name     string `json:"name"`
	Email    string `json:"upn"`
	TenantID string `json:"tid"`
	ObjectID string `json:"oid"`
	jwt.RegisteredClaims
}

// AuthedUser is the signed-in identity extracted from an Azure access token.
// ObjectID is the Entra principal ID used for role activation requests;
// TenantID and Email populate the UI header.
type AuthedUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	TenantID string `json:"tenantId"`
	ObjectID string `json:"objectId"`
}

// GetIdentity decodes (without verifying) the access token returned by cred
// and extracts the signed-in user's identity claims. The token is treated as
// trusted input — Azure is the issuer and the SDK has already validated it,
// so a local signature check would only add a JWKS fetch for no gain.
func GetIdentity(ctx context.Context, cred azcore.TokenCredential) (*AuthedUser, error) {
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("get azure token: %w", err)
	}

	claims := &userClaims{}
	jwtParser := jwt.NewParser()

	_, _, err = jwtParser.ParseUnverified(token.Token, claims)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	return convertClaims(claims), nil
}

func convertClaims(claims *userClaims) *AuthedUser {
	return &AuthedUser{
		Name:     claims.Name,
		Email:    claims.Email,
		TenantID: claims.TenantID,
		ObjectID: claims.ObjectID,
	}
}

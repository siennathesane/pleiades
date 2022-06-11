package testdata

import (
	"context"
)

type TestCookieMonsterServer struct {
	DRPCCookieMonsterUnimplementedServer
	// struct fields
}

// EatCookie turns a cookie into crumbs.
func (s *TestCookieMonsterServer) EatCookie(ctx context.Context, cookie *Cookie) (*Crumbs, error) {
	return &Crumbs{
		Cookie: cookie,
	}, nil
}

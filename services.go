package main

import (
	"context"
	"fmt"
)

type services struct {
	policy policyService
	user   *userService
}

func createServices(ctx context.Context) (services, error) {
	policy, err := newOPASDK(ctx)
	if err != nil {
		return services{}, fmt.Errorf("create policy: %w", err)
	}

	user, err := newUserService(ctx)
	if err != nil {
		return services{}, fmt.Errorf("create user: %w", err)
	}

	return services{
		policy: policy, user: user,
	}, nil
}

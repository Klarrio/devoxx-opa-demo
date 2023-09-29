package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/sdk"
)

type policyService interface {
	eval(context.Context, policyRequest) (bool, error)
}

// opaSDK implements policyService
type opaSDK struct {
	opa *sdk.OPA
}

func newOPASDK(ctx context.Context) (*opaSDK, error) {
	f, err := os.Open("./resources/opa_config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sdk, err := sdk.New(ctx, sdk.Options{
		Logger: logging.New(),
		Config: f,
	})
	if err != nil {
		return nil, fmt.Errorf("creating opa sdk instance: %w", err)
	}

	return &opaSDK{opa: sdk}, nil
}

func (o *opaSDK) eval(ctx context.Context, r policyRequest) (bool, error) {
	res, err := o.opa.Decision(ctx, sdk.DecisionOptions{
		Input: r,
		Path:  "policy/authz",
	})
	if err != nil {
		return false, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	authz, ok := res.Result.(bool)
	if !ok {
		return false, fmt.Errorf("authz not bool, but %T (%#v)", ok, res.Result)
	}
	return authz, nil
}

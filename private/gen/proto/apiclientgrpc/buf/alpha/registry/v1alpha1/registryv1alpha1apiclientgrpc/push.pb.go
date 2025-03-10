// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-apiclientgrpc. DO NOT EDIT.

package registryv1alpha1apiclientgrpc

import (
	context "context"
	v1alpha11 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/module/v1alpha1"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
)

type pushService struct {
	logger          *zap.Logger
	client          v1alpha1.PushServiceClient
	contextModifier func(context.Context) context.Context
}

// Push pushes.
func (s *pushService) Push(
	ctx context.Context,
	owner string,
	repository string,
	branch string,
	module *v1alpha11.Module,
	tags []string,
	tracks []string,
) (localModulePin *v1alpha1.LocalModulePin, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.Push(
		ctx,
		&v1alpha1.PushRequest{
			Owner:      owner,
			Repository: repository,
			Branch:     branch,
			Module:     module,
			Tags:       tags,
			Tracks:     tracks,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.LocalModulePin, nil
}

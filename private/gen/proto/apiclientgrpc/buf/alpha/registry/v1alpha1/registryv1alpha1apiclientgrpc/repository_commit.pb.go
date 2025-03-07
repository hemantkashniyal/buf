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
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
)

type repositoryCommitService struct {
	logger          *zap.Logger
	client          v1alpha1.RepositoryCommitServiceClient
	contextModifier func(context.Context) context.Context
}

// ListRepositoryCommitsByBranch lists the repository commits associated
// with a repository branch on a repository, ordered by their create time.
func (s *repositoryCommitService) ListRepositoryCommitsByBranch(
	ctx context.Context,
	repositoryOwner string,
	repositoryName string,
	repositoryBranchName string,
	pageSize uint32,
	pageToken string,
	reverse bool,
) (repositoryCommits []*v1alpha1.RepositoryCommit, nextPageToken string, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.ListRepositoryCommitsByBranch(
		ctx,
		&v1alpha1.ListRepositoryCommitsByBranchRequest{
			RepositoryOwner:      repositoryOwner,
			RepositoryName:       repositoryName,
			RepositoryBranchName: repositoryBranchName,
			PageSize:             pageSize,
			PageToken:            pageToken,
			Reverse:              reverse,
		},
	)
	if err != nil {
		return nil, "", err
	}
	return response.RepositoryCommits, response.NextPageToken, nil
}

// ListRepositoryCommitsByReference returns repository commits up-to and including
// the provided reference.
func (s *repositoryCommitService) ListRepositoryCommitsByReference(
	ctx context.Context,
	repositoryOwner string,
	repositoryName string,
	reference string,
	pageSize uint32,
	pageToken string,
	reverse bool,
) (repositoryCommits []*v1alpha1.RepositoryCommit, nextPageToken string, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.ListRepositoryCommitsByReference(
		ctx,
		&v1alpha1.ListRepositoryCommitsByReferenceRequest{
			RepositoryOwner: repositoryOwner,
			RepositoryName:  repositoryName,
			Reference:       reference,
			PageSize:        pageSize,
			PageToken:       pageToken,
			Reverse:         reverse,
		},
	)
	if err != nil {
		return nil, "", err
	}
	return response.RepositoryCommits, response.NextPageToken, nil
}

// GetRepositoryCommitByReference returns the repository commit matching
// the provided reference, if it exists.
func (s *repositoryCommitService) GetRepositoryCommitByReference(
	ctx context.Context,
	repositoryOwner string,
	repositoryName string,
	reference string,
) (repositoryCommit *v1alpha1.RepositoryCommit, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetRepositoryCommitByReference(
		ctx,
		&v1alpha1.GetRepositoryCommitByReferenceRequest{
			RepositoryOwner: repositoryOwner,
			RepositoryName:  repositoryName,
			Reference:       reference,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.RepositoryCommit, nil
}

// GetRepositoryCommitBySequenceId returns the repository commit matching
// the provided sequence ID and branch, if it exists.
func (s *repositoryCommitService) GetRepositoryCommitBySequenceId(
	ctx context.Context,
	repositoryOwner string,
	repositoryName string,
	repositoryBranchName string,
	commitSequenceId int64,
) (repositoryCommit *v1alpha1.RepositoryCommit, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetRepositoryCommitBySequenceId(
		ctx,
		&v1alpha1.GetRepositoryCommitBySequenceIdRequest{
			RepositoryOwner:      repositoryOwner,
			RepositoryName:       repositoryName,
			RepositoryBranchName: repositoryBranchName,
			CommitSequenceId:     commitSequenceId,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.RepositoryCommit, nil
}

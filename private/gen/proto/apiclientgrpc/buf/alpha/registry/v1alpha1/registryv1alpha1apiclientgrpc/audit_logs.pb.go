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
	v1alpha11 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/audit/v1alpha1"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type auditLogsService struct {
	logger          *zap.Logger
	client          v1alpha1.AuditLogsServiceClient
	contextModifier func(context.Context) context.Context
}

// ListAuditLogs lists audit logs matching the filters specified.
func (s *auditLogsService) ListAuditLogs(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
	reverse bool,
	startTime *timestamppb.Timestamp,
	endTime *timestamppb.Timestamp,
) (auditLogs []*v1alpha11.Event, nextPageToken string, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.ListAuditLogs(
		ctx,
		&v1alpha1.ListAuditLogsRequest{
			PageSize:  pageSize,
			PageToken: pageToken,
			Reverse:   reverse,
			StartTime: startTime,
			EndTime:   endTime,
		},
	)
	if err != nil {
		return nil, "", err
	}
	return response.AuditLogs, response.NextPageToken, nil
}

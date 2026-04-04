// Copyright 2019 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

/*
Package blobs contains a gRPC service to be used for remote file access.

It is used for bulk file reads and writes to files on any CockroachDB node.
Each node will run a blob service, which serves the file access for files on
that node. Each node will also have a blob client, which uses the nodedialer
to connect to another node's blob service, and access its files. The blob client
is the point of entry to this service and it supports the `BlobClient` interface,
which includes the following functionalities:
  - ReadFile
  - WriteFile
  - List
  - Delete
  - Stat
*/
package blobs

import (
	"context"
	"io"

	"github.com/cockroachdb/cockroach/pkg/blobs/blobspb"
	"github.com/cockroachdb/cockroach/pkg/util/ioctx"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/oserror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Service implements the gRPC BlobService which exchanges bulk files between different nodes.
type Service struct {
	localStorage *LocalStorage
}

var _ blobspb.BlobServer = &Service{}

// NewBlobService instantiates a blob service server.
func NewBlobService(externalIODir string) (*Service, error) {
	localStorage, err := NewLocalStorage(externalIODir)
	return &Service{localStorage: localStorage}, err
}

// GetStream implements the gRPC service.
func (s *Service) GetStream(req *blobspb.GetRequest, stream blobspb.Blob_GetStreamServer) error {
	content, _, err := s.localStorage.ReadFile(req.Filename, req.Offset)
	if err != nil {
		return err
	}
	defer content.Close(stream.Context())
	return streamContent(stream.Context(), stream, content)
}

// PutStream implements the gRPC service.
func (s *Service) PutStream(stream blobspb.Blob_PutStreamServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("could not fetch metadata")
	}
	filename := md.Get("filename")
	if len(filename) < 1 || filename[0] == "" {
		return errors.New("no filename in metadata")
	}
	reader := newPutStreamReader(stream)
	defer reader.Close(stream.Context())
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	w, err := s.localStorage.Writer(ctx, filename[0])
	if err != nil {
		cancel()
		return err
	}
	if _, err := io.Copy(w, ioctx.ReaderCtxAdapter(stream.Context(), reader)); err != nil {
		cancel()
		return errors.CombineErrors(w.Close(), err)
	}
	err = w.Close()
	cancel()
	return err
}

// List implements the gRPC service.
func (s *Service) List(
	ctx context.Context, req *blobspb.GlobRequest,
) (*blobspb.GlobResponse, error) {
	matches, err := s.localStorage.List(req.Pattern)
	return &blobspb.GlobResponse{Files: matches}, err
}

// Delete implements the gRPC service.
func (s *Service) Delete(
	ctx context.Context, req *blobspb.DeleteRequest,
) (*blobspb.DeleteResponse, error) {
	return &blobspb.DeleteResponse{}, s.localStorage.Delete(req.Filename)
}

// Stat implements the gRPC service.
func (s *Service) Stat(ctx context.Context, req *blobspb.StatRequest) (*blobspb.BlobStat, error) {
	resp, err := s.localStorage.Stat(req.Filename)
	if oserror.IsNotExist(err) {
		// gRPC hides the underlying golang ErrNotExist error, so we send back an
		// equivalent gRPC error which can be handled gracefully on the client side.
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return resp, err
}

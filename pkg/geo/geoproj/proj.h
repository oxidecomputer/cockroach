// Copyright 2020 The Cockroach Authors.
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

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

// CR_PROJ_Slice contains data that does not need to be freed.
// It can be either a Go or C pointer (which indicates who allocated the
// memory).
typedef struct {
  char* data;
  size_t len;
} CR_PROJ_Slice;

typedef CR_PROJ_Slice CR_PROJ_Status;

// CR_PROJ_Transform converts the given x/y/z coordinates to a new project specification.
// Note points (x[i], y[i], z[i]) are in the range 0 <= i < point_coint.
CR_PROJ_Status CR_PROJ_Transform(char* fromSpec, char* toSpec, long point_count, double* x,
                                 double* y, double* z);

// CR_PROJ_GetProjMetadata gets the metadata for a given spec in relation to
// the spheroid it represents.
CR_PROJ_Status CR_PROJ_GetProjMetadata(char* spec, int* retIsLatLng, double* retMajorAxis,
                                       double* retEccentricitySquared);
#ifdef __cplusplus
}  // extern "C"
#endif

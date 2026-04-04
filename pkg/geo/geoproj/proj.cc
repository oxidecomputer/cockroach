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

#include "proj.h"
#include <cstring>
#include <proj_api.h>
#include <stdlib.h>
#include <string>

const char* DEFAULT_ERROR_MSG = "PROJ could not parse proj4text";

namespace {
CR_PROJ_Status CR_PROJ_ErrorFromErrorCode(int code) {
  char* err = pj_strerrno(code);
  if (err == nullptr) {
    err = (char*)DEFAULT_ERROR_MSG;
  }
  return {.data = err, .len = strlen(err)};
}

}  // namespace

CR_PROJ_Status CR_PROJ_Transform(char* fromSpec, char* toSpec, long point_count, double* x,
                                 double* y, double* z) {
  CR_PROJ_Status err = {.data = NULL, .len = 0};
  auto ctx = pj_ctx_alloc();
  auto fromPJ = pj_init_plus_ctx(ctx, fromSpec);
  if (fromPJ == nullptr) {
    err = CR_PROJ_ErrorFromErrorCode(pj_ctx_get_errno(ctx));
    pj_ctx_free(ctx);
    return err;
  }
  auto toPJ = pj_init_plus_ctx(ctx, toSpec);
  if (toPJ == nullptr) {
    err = CR_PROJ_ErrorFromErrorCode(pj_ctx_get_errno(ctx));
    pj_free(fromPJ);
    pj_ctx_free(ctx);
    return err;
  }
  // If we have a latlng from, transform to radians.
  if (pj_is_latlong(fromPJ)) {
    for (auto i = 0; i < point_count; i++) {
      x[i] = x[i] * DEG_TO_RAD;
      y[i] = y[i] * DEG_TO_RAD;
    }
  }
  pj_transform(fromPJ, toPJ, point_count, 0, x, y, z);
  int errCode = pj_ctx_get_errno(ctx);
  if (errCode != 0) {
    err = CR_PROJ_ErrorFromErrorCode(errCode);
    pj_free(toPJ);
    pj_free(fromPJ);
    pj_ctx_free(ctx);
    return err;
  }

  // If we have a latlng to, transform to degrees.
  if (pj_is_latlong(toPJ)) {
    for (auto i = 0; i < point_count; i++) {
      x[i] = x[i] * RAD_TO_DEG;
      y[i] = y[i] * RAD_TO_DEG;
    }
  }
  pj_free(toPJ);
  pj_free(fromPJ);
  pj_ctx_free(ctx);
  return err;
}

CR_PROJ_Status CR_PROJ_GetProjMetadata(char* spec, int* retIsLatLng, double* retMajorAxis,
                                       double* retEccentricitySquared) {
  CR_PROJ_Status err = {.data = NULL, .len = 0};
  auto ctx = pj_ctx_alloc();
  auto pj = pj_init_plus_ctx(ctx, spec);
  if (pj == nullptr) {
    err = CR_PROJ_ErrorFromErrorCode(pj_ctx_get_errno(ctx));
    pj_ctx_free(ctx);
    return err;
  }

  *retIsLatLng = pj_is_latlong(pj);
  pj_get_spheroid_defn(pj, retMajorAxis, retEccentricitySquared);

  pj_free(pj);
  pj_ctx_free(ctx);
  return err;
}

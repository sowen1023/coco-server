// Copyright (C) INFINI Labs & INFINI LIMITED.
//
// The INFINI Console is offered under the GNU Affero General Public License v3.0
// and as commercial software.
//
// For commercial licensing, contact us at:
//   - Website: infinilabs.com
//   - Email: hello@infini.ltd
//
// Open Source licensed under AGPL V3:
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

/* Copyright © INFINI Ltd. All rights reserved.
 * web: https://infinilabs.com
 * mail: hello#infini.ltd */

package security

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"infini.sh/coco/core"
	"infini.sh/coco/modules/common"
	"infini.sh/framework/core/api"
	httprouter "infini.sh/framework/core/api/router"
	"infini.sh/framework/core/kv"
	"infini.sh/framework/core/security"
	"infini.sh/framework/core/util"
)

func (h APIHandler) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	api.DestroySession(w, r)
	h.WriteOKJSON(w, util.MapStr{
		"status": "ok",
	})
}

func (h APIHandler) Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if !api.IsAuthEnable() {
		panic("auth is not enabled")
	}

	reqUser, err := security.GetUserFromContext(r.Context())
	if err != nil || reqUser == nil {
		panic("invalid user")
	}

	var data []byte
	cfg, _ := common.AppConfigFromFile()
	if cfg != nil && cfg.ServerInfo != nil && cfg.ServerInfo.Managed {
		data, err = kv.GetValue(core.UserProfileKey, []byte(reqUser.UserID))
	} else {
		//TODO to be removed
		data, err = kv.GetValue(core.DefaultSettingBucketKey, []byte(core.DefaultUserProfileKey))
	}

	h.WriteBytes(w, data, 200)
}

func (h APIHandler) UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cfg, _ := common.AppConfigFromFile()
	if cfg.ServerInfo.Managed {
		panic("should not be invoked as in managed mode")
	}

	reqUser, err := security.GetUserFromContext(r.Context())
	if err != nil {
		panic(err)
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	err = h.DecodeJSON(r, &req)
	if err != nil {
		h.ErrorInternalServer(w, err.Error())
		return
	}

	err, success := h.checkPassword(req.OldPassword)
	if !success {
		h.WriteError(w, "failed to login", 403)
		return
	}

	err = SavePassword(req.NewPassword)
	if err != nil {
		h.ErrorInternalServer(w, err.Error())
		return
	}
	h.WriteOKJSON(w, api.UpdateResponse(reqUser.UserID))
	return
}

func SavePassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = kv.AddValue(core.DefaultSettingBucketKey, []byte(core.DefaultUserPasswordKey), hash)
	return err
}

func (h APIHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cfg, _ := common.AppConfigFromFile()
	if cfg.ServerInfo.Managed {
		panic("should not be invoked as in managed mode")
	}

	var req struct {
		Password string `json:"password"`
	}

	var fromForm = false
	var requestID = h.GetParameter(r, "request_id")

	// Check content type and parse accordingly
	contentType := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "application/json"):
		// Handle JSON input
		err := h.DecodeJSON(r, &req)
		if err != nil {
			h.ErrorInternalServer(w, "invalid JSON format")
			return
		}

	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"),
		strings.HasPrefix(contentType, "multipart/form-data"):
		// Handle form input
		if err := r.ParseForm(); err != nil {
			h.ErrorInternalServer(w, "failed to parse form data")
			return
		}
		fromForm = true
		req.Password = r.PostFormValue("password")

	default:
		h.WriteError(w, "unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	// Validate password exists
	if req.Password == "" {
		h.WriteError(w, "password is required", http.StatusBadRequest)
		return
	}

	// Rest of your existing logic
	err, success := h.checkPassword(req.Password)
	if !success {
		h.WriteError(w, "failed to login", http.StatusForbidden)
		return
	}

	var user = &security.UserProfile{
		Name: core.DefaultUserLogin,
	}
	user.ID = core.DefaultUserLogin

	sessionInfo := security.UserSessionInfo{}
	sessionInfo.Provider = "simple"
	sessionInfo.Login = core.DefaultUserLogin

	sessionInfo.TenantID = "LOCAL"
	sessionInfo.UserID = user.ID
	//sessionInfo.Profile = user
	sessionInfo.Roles = []string{security.RoleAdmin}

	err, token := AddUserAccessTokenToSession(w, r, &sessionInfo)
	if err != nil {
		h.ErrorInternalServer(w, "failed to authorize user")
		return
	}

	if fromForm {
		h.Redirect(w, r, fmt.Sprintf("/login/success?request_id=%v&code=%v", requestID, token["access_token"]))
	} else {
		h.WriteOKJSON(w, token)
	}
}

func AddUserAccessTokenToSession(w http.ResponseWriter, r *http.Request, user *security.UserSessionInfo) (error, map[string]interface{}) {

	if user == nil {
		panic("invalid user")
	}

	// Generate access token
	token, err := GenerateJWTAccessToken(user)
	if err != nil {
		return err, nil
	}

	api.SetSession(w, r, core.UserAccessTokenSessionName, token["access_token"])
	return nil, token
}

func (h APIHandler) checkPassword(password string) (error, bool) {
	savedPassword, err := kv.GetValue(core.DefaultSettingBucketKey, []byte(core.DefaultUserPasswordKey))
	if err != nil {
		return err, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if err != nil {
		return err, false
	}
	return nil, true
}

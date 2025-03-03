// +build coap,!http

/*
 * Copyright 2020-2022 ForgeRock AS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import "errors"

var errHTTPNotBuilt = errors.New("http(s) scheme is unsupported")

func (c amConnection) Initialise() error {
	return errHTTPNotBuilt
}

func (c amConnection) Authenticate(payload AuthenticatePayload) (reply AuthenticatePayload, err error) {
	return reply, errHTTPNotBuilt
}

func (c amConnection) AMInfo() (info AMInfoResponse, err error) {
	return info, errHTTPNotBuilt
}

func (c amConnection) ValidateSession(tokenID string, content ContentType, payload string) (ok bool, err error) {
	return ok, errHTTPNotBuilt
}

func (c amConnection) LogoutSession(tokenID string, content ContentType, payload string) (err error) {
	return errHTTPNotBuilt
}

func (c amConnection) AccessToken(tokenID string, content ContentType, payload string) (reply []byte, err error) {
	return reply, errHTTPNotBuilt
}

func (c *amConnection) IntrospectAccessToken(tokenID string, content ContentType, payload string) (introspection []byte, err error) {
	return introspection, errHTTPNotBuilt
}

func (c amConnection) Attributes(tokenID string, content ContentType, payload string, names []string) (reply []byte, err error) {
	return reply, errHTTPNotBuilt
}

func (c *amConnection) UserCode(tokenID string, content ContentType, payload string) (reply []byte, err error) {
	return reply, errHTTPNotBuilt
}

func (c *amConnection) UserToken(tokenID string, content ContentType, payload string) (reply []byte, err error) {
	return reply, errHTTPNotBuilt
}

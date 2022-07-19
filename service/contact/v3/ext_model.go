/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkcontact

import larkevent "github.com/larksuite/oapi-sdk-go/v3/event"

type P1UserChangedV3Data struct {
	Type       string `json:"type"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
	OpenID     string `json:"open_id,omitempty"`
	EmployeeId string `json:"employee_id"`
	UnionId    string `json:"union_id,omitempty"`
}

type P1UserChangedV3 struct {
	*larkevent.EventBase
	*larkevent.EventReq
	Event *P1UserChangedV3Data `json:"event"`
}

func (m *P1UserChangedV3) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1UserStatusV3 struct {
	IsActive   bool `json:"is_active"`
	IsFrozen   bool `json:"is_frozen"`
	IsResigned bool `json:"is_resigned"`
}
type P1UserStatusChangedV3Data struct {
	Type          string          `json:"type"`
	AppID         string          `json:"app_id"`
	TenantKey     string          `json:"tenant_key"`
	OpenID        string          `json:"open_id,omitempty"`
	EmployeeId    string          `json:"employee_id"`
	UnionId       string          `json:"union_id,omitempty"`
	BeforeStatus  *P1UserStatusV3 `json:"before_status"`
	CurrentStatus *P1UserStatusV3 `json:"current_status"`
	ChangeTime    string          `json:"change_time"`
}

type P1UserStatusChangedV3 struct {
	*larkevent.EventBase
	*larkevent.EventReq
	Event *P1UserStatusChangedV3Data `json:"event"`
}

func (m *P1UserStatusChangedV3) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1DepartmentChangedV3Data struct {
	Type             string `json:"type"`
	AppID            string `json:"app_id"`
	TenantKey        string `json:"tenant_key"`
	OpenID           string `json:"open_id,omitempty"`
	OpenDepartmentId string `json:"open_department_id"`
}

type P1DepartmentChangedV3 struct {
	*larkevent.EventBase
	*larkevent.EventReq
	Event *P1DepartmentChangedV3Data `json:"event"`
}

func (m *P1DepartmentChangedV3) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1ContactScopeChangedV3Data struct {
	Type      string `json:"type"`
	AppID     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
}

type P1ContactScopeChangedV3 struct {
	*larkevent.EventBase
	*larkevent.EventReq
	Event *P1ContactScopeChangedV3Data `json:"event"`
}

func (m *P1ContactScopeChangedV3) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

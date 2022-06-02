package messages

/*
 * This file is part of the ObjectVault Project.
 * Copyright (C) 2020-2022 Paulo Ferreira <vault at sourcenotes.org>
 *
 * This work is published under the GNU AGPLv3.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/objectvault/queue-interface/shared"
)

type InvitationEmailMessage struct {
	EmailMessage        // DERIVED FROM
	code         string // [REQUIRED] Invitation Code
	byUser       string // [REQUIRED] User Name (orm.User.Name())
	atUser       string // [REQUIRED] User ID (orm.User.Email())
	message      string // [OPTIONAL] Message to Invitee
	objectName   string // [REQUIRED] Store Name (orm.Store.Name() or orm.Org.Name())
}

func (m *InvitationEmailMessage) IsValid() bool {
	return m.EmailMessage.IsValid() && (m.code != "") && (m.byUser != "") && (m.objectName != "")
}

func (m *InvitationEmailMessage) Code() string {
	return m.code
}

func (m *InvitationEmailMessage) SetCode(code string) (string, error) {
	// Is Software Activation Code Empty?
	code = strings.TrimSpace(code)
	if code == "" {
		return "", errors.New("Invitaion Code is Required")
	}

	// Code is always lower case
	code = strings.ToLower(code)

	// Current State
	current := m.code

	// New State
	m.code = code
	return current, nil
}

func (m *InvitationEmailMessage) ByUser() string {
	return m.byUser
}

func (m *InvitationEmailMessage) SetByUser(user string) (string, error) {
	// Is Name Empty?
	user = strings.TrimSpace(user)
	if user == "" {
		return "", errors.New("User Name is Required")
	}

	// Current State
	current := m.byUser

	// New State
	m.byUser = user
	return current, nil
}

func (m *InvitationEmailMessage) AtUser() string {
	return m.atUser
}

func (m *InvitationEmailMessage) SetAtUser(user string) (string, error) {
	// Is Email Empty?
	user = strings.TrimSpace(user)
	if user == "" {
		return "", errors.New("From User Email is Required")
	}

	// User ID is always lower case
	user = strings.ToLower(user)

	// Current State
	current := m.atUser

	// New State
	m.atUser = user
	return current, nil
}

func (m *InvitationEmailMessage) Message() string {
	return m.message
}

func (m *InvitationEmailMessage) SetMessage(msg string) (string, error) {
	// Is Template Name Empty?
	msg = strings.TrimSpace(msg)

	// Current State
	current := m.message

	// New State
	m.message = msg
	return current, nil
}

func (m *InvitationEmailMessage) ObjectName() string {
	return m.objectName
}

func (m *InvitationEmailMessage) SetObjectName(name string) (string, error) {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("Object Name is Required")
	}

	// Current State
	current := m.objectName

	// New State
	m.objectName = name
	return current, nil
}

// MarshalJSON implements json.Marshal
func (m InvitationEmailMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[InvitationEmailMessage] Message is Invalid")
	}

	// Is Message Creation Date Set?
	if m.created == "" { // NO: Use Current Time
		m.created = shared.UTCTimeStamp()
	}

	// QUEUE Counter and Settings //
	queue := &struct {
		RequeueCount int    `json:"count,omitempty"`
		ErrorCode    int    `json:"errorcode,omitempty"`
		ErrorTime    string `json:"errortime,omitempty"`
		ErrorMessage string `json:"errormsg,omitempty"`
	}{}

	// Has the Message been Requeued?
	if m.requeueCount > 0 { // YES
		queue.RequeueCount = m.requeueCount
	}

	// Is this an Error Message?
	if m.errorCode > 0 { // YES
		queue.ErrorCode = m.errorCode
		queue.ErrorTime = m.errorTime
		queue.ErrorMessage = m.errorMessage
	}

	// Email Settings //
	email := &struct {
		Template string             `json:"template"`
		Locale   string             `json:"locale"`
		To       string             `json:"to"`
		From     string             `json:"from,omitempty"`
		CC       string             `json:"cc,omitempty"`
		BCC      string             `json:"bcc,omitempty"`
		Headers  *map[string]string `json:"headers,omitempty"`
	}{
		Template: m.template,
		Locale:   m.Locale(),
		To:       m.to,
		From:     m.from,
		CC:       m.cc,
		BCC:      m.bcc,
		Headers:  m.headers,
	}

	// Email Invitation //
	invite := &struct {
		Code    string `json:"code"`
		Invitee string `json:"invitee"`
		Email   string `json:"invitee_email"`
		To      string `json:"to"`
	}{
		Code:    m.code,
		Invitee: m.byUser,
		Email:   m.atUser,
		To:      m.objectName,
	}
	// Complete JSON Message //
	output := &struct {
		Version int                     `json:"version"`
		ID      string                  `json:"id"`
		Type    string                  `json:"type"`
		SubType string                  `json:"subtype,omitempty"`
		Params  *map[string]interface{} `json:"data,omitempty"`
		Created string                  `json:"created"`
		Queue   interface{}             `json:"queue,omitempty"`
		Email   interface{}             `json:"email"`
		Invite  interface{}             `json:"invite"`
	}{
		Version: m.version,
		ID:      m.id,
		Type:    m.mtype,
		SubType: m.msubtype,
		Params:  m.params,
		Created: m.created,
		Queue:   queue,
		Email:   email,
		Invite:  invite,
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *InvitationEmailMessage) UnmarshalJSON(b []byte) error {
	in := &struct {
		Version int                     `json:"version"`
		ID      string                  `json:"id"`
		Type    string                  `json:"type"`
		SubType string                  `json:"subtype,omitempty"`
		Params  *map[string]interface{} `json:"params,omitempty"`
		Created string                  `json:"created"`
		Queue   *struct {
			RequeueCount int    `json:"count,omitempty"`
			ErrorCode    int    `json:"errorcode,omitempty"`
			ErrorTime    string `json:"errortime,omitempty"`
			ErrorMessage string `json:"errormsg,omitempty"`
		} `json:"errormsg,omitempty"`
		Email *struct {
			Template string             `json:"template"`
			Locale   string             `json:"locale"`
			To       string             `json:"to"`
			From     string             `json:"from,omitempty"`
			CC       string             `json:"cc,omitempty"`
			BCC      string             `json:"bcc,omitempty"`
			Headers  *map[string]string `json:"headers,omitempty"`
		} `json:"email"`
		Invite *struct {
			Code    string `json:"code"`
			Invitee string `json:"invitee"`
			Email   string `json:"invitee_email"`
			To      string `json:"to"`
		} `json:"invite"`
	}{}

	err := json.Unmarshal(b, &in)
	if err != nil {
		return err
	}

	// Basic Message Information
	m.version = in.Version
	m.id = in.ID
	m.mtype = in.Type
	m.msubtype = in.Type
	m.params = in.Params
	m.created = in.Created

	// QUEUE Message Control Information //
	if in.Queue != nil {
		m.requeueCount = in.Queue.RequeueCount

		// Has Error Message?
		if in.Queue.ErrorCode > 0 { // YES
			m.errorCode = in.Queue.ErrorCode
			m.errorTime = in.Queue.ErrorTime
			m.errorMessage = in.Queue.ErrorMessage
		}
	}

	// EMAIL Message Parameters //
	m.template = in.Email.Template
	m.locale = in.Email.Locale
	m.to = in.Email.To
	m.from = in.Email.From
	m.cc = in.Email.CC
	m.bcc = in.Email.BCC
	m.headers = in.Email.Headers

	// INVITATION Message Parameters //
	m.code = in.Invite.Code
	m.byUser = in.Invite.Invitee
	m.atUser = in.Invite.Email
	m.objectName = in.Invite.To
	return nil
}

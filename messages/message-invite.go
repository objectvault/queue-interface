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
	"time"

	"github.com/objectvault/queue-interface/shared"
)

type InviteMessage struct {
	EmailMessage            // DERIVED FROM
	code         string     // [REQUIRED] Invitation Code
	byUser       string     // [REQUIRED] User Name (orm.User.Name())
	atUser       string     // [REQUIRED] User ID (orm.User.Email())
	message      string     // [OPTIONAL] Message to Invitee
	objectName   string     // [REQUIRED] Store Name (orm.Store.Name() or orm.Org.Name())
	expiration   *time.Time // [REQUIRED] Expiration Time Stamp
}

func NewInviteMessage(template string, code string) (*InviteMessage, error) {
	// Create Invitation
	m := &InviteMessage{}

	// Initialize Email Message Base
	err := InitEmailMessage(&m.EmailMessage, "invite", template)
	if err != nil {
		return nil, err
	}

	// Save Code (Note: ALLOW code == "")
	code = strings.TrimSpace(code)
	m.code = strings.ToLower(code)

	return m, nil
}

func (m *InviteMessage) IsValid() bool {
	return m.EmailMessage.IsValid() && (m.code != "") && (m.byUser != "") && (m.objectName != "") && (m.expiration != nil)
}

func (m *InviteMessage) Code() string {
	return m.code
}

func (m *InviteMessage) SetCode(code string) (string, error) {
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

func (m *InviteMessage) ByUser() string {
	return m.byUser
}

func (m *InviteMessage) SetByUser(user string) (string, error) {
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

func (m *InviteMessage) AtUser() string {
	return m.atUser
}

func (m *InviteMessage) SetAtUser(user string) (string, error) {
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

func (m *InviteMessage) Message() string {
	return m.message
}

func (m *InviteMessage) SetMessage(msg string) (string, error) {
	// Is Template Name Empty?
	msg = strings.TrimSpace(msg)

	// Current State
	current := m.message

	// New State
	m.message = msg
	return current, nil
}

func (m *InviteMessage) ObjectName() string {
	return m.objectName
}

func (m *InviteMessage) SetObjectName(name string) (string, error) {
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

func (m *InviteMessage) Expiration() *time.Time {
	return m.expiration
}

func (m *InviteMessage) ExpirationUTC() string {
	if m.expiration != nil {
		utc := m.expiration.UTC()
		// RETURN ISO 8601 / RFC 3339 FORMAT in UTC
		return utc.Format(time.RFC3339)
	}

	return ""
}

func (m *InviteMessage) SetExpiration(t time.Time) (*time.Time, error) {
	// Current State
	current := m.expiration

	// New State
	m.expiration = &t
	return current, nil
}

// MarshalJSON implements json.Marshal
func (m InviteMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[InviteMessage] Message is Invalid")
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
		Code       string `json:"code"`
		Invitee    string `json:"invitee"`
		Email      string `json:"invitee_email"`
		To         string `json:"to"`
		Message    string `json:"message,omitempty"`
		Expiration string `json:"expiration"`
	}{
		Code:       m.code,
		Invitee:    m.byUser,
		Email:      m.atUser,
		To:         m.objectName,
		Message:    m.message,
		Expiration: m.ExpirationUTC(),
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
func (m *InviteMessage) UnmarshalJSON(b []byte) error {
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
			Code       string `json:"code"`
			Invitee    string `json:"invitee"`
			Email      string `json:"invitee_email"`
			To         string `json:"to"`
			Message    string `json:"message,omitempty"`
			Expiration string `json:"expiration"`
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
	m.message = in.Invite.Message

	// Extract Expiration Time (Expected ISO8601 / RFC3339 Format)
	ts := in.Invite.Expiration
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	m.expiration = &t

	return nil
}

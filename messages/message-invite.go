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
)

type InvitationEmailMessage struct {
	version    int    // [REQUIRED] Message Version
	template   string // [REQUIRED] Email Template to Use
	language   string // [OPTIONAL:DEFAULT en_US] System Default Language
	to         string // [REQUIRED] Invitee's Email Address
	from       string // [OPTIONAL] default no-reply@...
	code       string // [REQUIRED] Invitation Code
	byUser     string // [REQUIRED] User Name (orm.User.Name())
	atUser     string // [REQUIRED] User ID (orm.User.Email())
	message    string // [OPTIONAL] Message to Invitee
	objectName string // [REQUIRED] Store Name (orm.Store.Name() or orm.Org.Name())
}

func (m *InvitationEmailMessage) Version() int {
	if m.version > 0 {
		return m.version
	}

	return 1
}

func (m *InvitationEmailMessage) SetVersion(v int) (int, error) {
	// Valid Version?
	if v <= 0 { // NO
		return 0, errors.New("[InvitationEmailMessage] Invalid Message Version")
	}

	// Current State
	current := m.version

	// New State
	m.version = v
	return current, nil
}

func (m *InvitationEmailMessage) Template() string {
	return m.template
}

func (m *InvitationEmailMessage) SetTemplate(t string) (string, error) {
	// Is Template Name Empty?
	t = strings.TrimSpace(t)
	if t == "" {
		return "", errors.New("Invitation Template is Required")
	}

	// To is always lower case
	t = strings.ToLower(t)

	// Current State
	current := m.template

	// New State
	m.template = t
	return current, nil
}

func (m *InvitationEmailMessage) Language() string {
	if m.language == "" {
		return "en_US"
	}

	return m.language
}

func (m *InvitationEmailMessage) SetLanguage(l string) (string, error) {
	// Is Template Name Empty?
	l = strings.TrimSpace(l)

	// Language Tuplets are always lower case
	l = strings.ToLower(l)

	// Current State
	current := m.language

	// New State
	m.language = l
	return current, nil
}

func (m *InvitationEmailMessage) To() string {
	return m.to
}

func (m *InvitationEmailMessage) SetTo(to string) (string, error) {
	// Is Template Name Empty?
	to = strings.TrimSpace(to)
	if to == "" {
		return "", errors.New("Email Destination is Required")
	}

	// To is always lower case
	to = strings.ToLower(to)

	// Current State
	current := m.to

	// New State
	m.to = to
	return current, nil
}

func (m *InvitationEmailMessage) From(d string) string {
	if m.from == "" {
		return d
	}

	return m.from
}

func (m *InvitationEmailMessage) SetFrom(from string) (string, error) {
	// Is Template Name Empty?
	from = strings.TrimSpace(from)

	// From is always lower case
	from = strings.ToLower(from)

	// Current State
	current := m.from

	// New State
	m.from = from
	return current, nil
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

func (m *InvitationEmailMessage) IsValid() bool {
	return (m.template != "") && (m.to != "") && (m.code != "") && (m.byUser != "") && (m.atUser != "") && (m.objectName != "")
}

// MarshalJSON implements json.Marshal
func (m InvitationEmailMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[InvitationEmailMessage] Message is Invalid")
	}

	return json.Marshal(&struct {
		Version    int    `json:"version"`
		Template   string `json:"template"`
		Language   string `json:"locale"`
		To         string `json:"to"`
		From       string `json:"from,omitempty"`
		Code       string `json:"code"`
		ByUser     string `json:"invitee"`
		AtUser     string `json:"invitee-email"`
		ObjectName string `json:"object-name"`
	}{
		Version:    m.Version(),
		Template:   m.Template(),
		Language:   m.Language(),
		To:         m.to,
		From:       m.from,
		Code:       m.code,
		ByUser:     m.byUser,
		AtUser:     m.atUser,
		ObjectName: m.objectName,
	})
}

// UnmarshalJSON implements json.Unmarshal
func (m *InvitationEmailMessage) UnmarshalJSON(b []byte) error {
	iem := &struct {
		Version    int    `json:"version"`
		Template   string `json:"template"`
		Language   string `json:"locale"`
		To         string `json:"to"`
		From       string `json:"from,omitempty"`
		Code       string `json:"code"`
		ByUser     string `json:"invitee"`
		AtUser     string `json:"invitee-email"`
		ObjectName string `json:"object-name"`
	}{}

	err := json.Unmarshal(b, &iem)
	if err != nil {
		return err
	}

	m.version = iem.Version
	m.template = iem.Template
	m.language = iem.Language
	m.to = iem.To
	m.from = iem.From
	m.code = iem.Code
	m.byUser = iem.ByUser
	m.atUser = iem.AtUser
	m.objectName = iem.ObjectName
	return nil
}

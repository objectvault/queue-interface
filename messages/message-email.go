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
	"strconv"
	"strings"
)

type EmailMessage struct {
	version  int                     // [REQUIRED] Message Version
	template string                  // [REQUIRED] Email Template to Use
	language string                  // [OPTIONAL:DEFAULT en_US] User Preferred Language
	to       string                  // [REQUIRED] User Email (i.e. test@test.net)
	from     string                  // [OPTIONAL] default no-reply@test-to.com
	cc       string                  // [OPTIONAL] cc1@test-to.com;cc2@test-to.com;
	bcc      string                  // [OPTIONAL] bcc1@test-to.com;bcc2@test-to.com;
	params   *map[string]interface{} // [OPTIONAL] TEMPLATE Parameters first.name = 'Paulo', etc.
	headers  *map[string]string      // [OPTIONAL] Extra Headers
}

func (m *EmailMessage) IsValid() bool {
	return (m.template != "") && (m.to != "") && (m.params != nil)
}

func (m *EmailMessage) Version() int {
	if m.version > 0 {
		return m.version
	}

	return 1
}

func (m *EmailMessage) SetVersion(v int) (int, error) {
	// Valid Version?
	if v <= 0 { // NO
		return 0, errors.New("[LicenseCreateMessage] Invalid Message Version")
	}

	// Current State
	current := m.version

	// New State
	m.version = v
	return current, nil
}

func (m *EmailMessage) Template() string {
	return m.template
}

func (m *EmailMessage) SetTemplate(t string) (string, error) {
	// Is Template Name Empty?
	t = strings.TrimSpace(t)
	if t == "" {
		return "", errors.New("Email Template is Required")
	}

	// To is always lower case
	t = strings.ToLower(t)

	// Current State
	current := m.template

	// New State
	m.template = t
	return current, nil
}

func (m *EmailMessage) Language() string {
	if m.language == "" {
		return "en_US"
	}

	return m.language
}

func (m *EmailMessage) SetLanguage(l string) (string, error) {
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

func (m *EmailMessage) To() string {
	return m.to
}

func (m *EmailMessage) SetTo(to string) (string, error) {
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

func (m *EmailMessage) From(d string) string {
	if m.from == "" {
		return d
	}

	return m.from
}

func (m *EmailMessage) SetFrom(from string) (string, error) {
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

func (m *EmailMessage) CC() string {
	return m.cc
}

func (m *EmailMessage) SetCC(cc string) (string, error) {
	// Is Template Name Empty?
	cc = strings.TrimSpace(cc)

	// CCs are always lower case
	cc = strings.ToLower(cc)

	// Current State
	current := m.cc

	// New State
	m.cc = cc
	return current, nil
}

func (m *EmailMessage) BCC() string {
	return m.bcc
}

func (m *EmailMessage) SetBCC(bcc string) (string, error) {
	// Is Template Name Empty?
	bcc = strings.TrimSpace(bcc)

	// BCC's are always lower case
	bcc = strings.ToLower(bcc)

	// Current State
	current := m.bcc

	// New State
	m.bcc = bcc
	return current, nil
}

func (m *EmailMessage) GetParameters() *map[string]interface{} {
	return m.params
}

func (m *EmailMessage) HasParameter(n string) bool {
	_, ok := (*m.params)[n]
	return ok
}

func (m *EmailMessage) Parameter(n string) interface{} {
	if m.params == nil {
		return ""
	}

	value, ok := (*m.params)[n]
	if ok {
		return value
	}
	return ""
}

func (m *EmailMessage) SetParameter(n string, v string) error {
	// Do we already have a Parameters Map?
	if m.params == nil { // NO: Create One
		d := make(map[string]interface{})
		m.params = &d
	}

	// New State
	(*m.params)[n] = v
	return nil
}

func (m *EmailMessage) SetIntParameter(n string, v int) error {
	return m.SetParameter(n, strconv.Itoa(v))
}

func (m *EmailMessage) UnsetParameter(n string) (interface{}, error) {
	if m.params == nil {
		return "", nil
	}

	// Do we have a Current Value?
	current, ok := (*m.params)[n]
	if ok { // YES: Delete It
		delete(*m.params, n)
	}
	// ELSE: No - Just Ignore Request

	return current, nil
}

func (m *EmailMessage) GetHeaders() *map[string]string {
	return m.headers
}

func (m *EmailMessage) HasHeader(n string) bool {
	if m.headers == nil {
		return false
	}

	_, ok := (*m.headers)[n]
	return ok
}

func (m *EmailMessage) Header(n string) string {
	if m.headers == nil {
		return ""
	}

	value, ok := (*m.headers)[n]
	if ok {
		return value
	}
	return ""
}

func (m *EmailMessage) SetHeader(n string, v string) error {
	// Do we already have a Headers Map?
	if m.headers == nil { // NO: Create One
		m.headers = &map[string]string{}
	}

	// New State
	(*m.headers)[n] = v
	return nil
}

func (m *EmailMessage) UnsetHeader(n string) (string, error) {
	if m.headers == nil {
		return "", nil
	}

	// Do we have a Current Value?
	current, ok := (*m.headers)[n]
	if ok { // YES: Delete It
		delete(*m.headers, n)
	}
	// ELSE: No - Just Ignore Request

	return current, nil
}

// MarshalJSON implements json.Marshal
func (m EmailMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[EmailMessage] Message is Invalid")
	}

	return json.Marshal(&struct {
		Version  int                     `json:"version"`
		Template string                  `json:"template"`
		Language string                  `json:"locale"`
		To       string                  `json:"to"`
		From     string                  `json:"from,omitempty"`
		CC       string                  `json:"cc,omitempty"`
		BCC      string                  `json:"bcc,omitempty"`
		Params   *map[string]interface{} `json:"params,omitempty"`
		Headers  *map[string]string      `json:"headers,omitempty"`
	}{
		Version:  m.Version(),
		Template: m.template,
		Language: m.Language(),
		To:       m.to,
		From:     m.from,
		CC:       m.cc,
		BCC:      m.bcc,
		Params:   m.params,
		Headers:  m.headers,
	})
}

// UnmarshalJSON implements json.Unmarshal
func (m *EmailMessage) UnmarshalJSON(b []byte) error {
	me := &struct {
		Version  int                     `json:"version"`
		Template string                  `json:"template"`
		Language string                  `json:"locale"`
		To       string                  `json:"to"`
		From     string                  `json:"from,omitempty"`
		CC       string                  `json:"cc,omitempty"`
		BCC      string                  `json:"bcc,omitempty"`
		Params   *map[string]interface{} `json:"params,omitempty"`
		Headers  *map[string]string      `json:"headers,omitempty"`
	}{}

	err := json.Unmarshal(b, &me)
	if err != nil {
		return err
	}

	m.version = me.Version
	m.template = me.Template
	m.language = me.Language
	m.to = me.To
	m.from = me.From
	m.cc = me.CC
	m.bcc = me.BCC
	m.params = me.Params
	m.headers = me.Headers
	return nil
}

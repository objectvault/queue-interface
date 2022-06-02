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

// cSpell:ignore mtype, msubtype

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/objectvault/queue-interface/shared"
)

type EmailMessage struct {
	QueueMessage                    // DERIVED FROM
	template     string             // [REQUIRED] Email Template to Use
	locale       string             // [OPTIONAL:DEFAULT en_US] User Preferred Language
	to           string             // [REQUIRED] User Email (i.e. test@test.net)
	from         string             // [OPTIONAL] default no-reply@test-to.com
	cc           string             // [OPTIONAL] cc1@test-to.com;cc2@test-to.com;
	bcc          string             // [OPTIONAL] bcc1@test-to.com;bcc2@test-to.com;
	headers      *map[string]string // [OPTIONAL] Extra Headers
}

func (m *EmailMessage) IsValid() bool {
	return m.QueueMessage.IsValid() && (m.template != "") && (m.to != "")
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

func (m *EmailMessage) Locale() string {
	if m.locale == "" {
		return "en_US"
	}

	return m.locale
}

func (m *EmailMessage) SetLocale(l string) (string, error) {
	// Is Template Name Empty?
	l = strings.TrimSpace(l)

	// Language Tuplets are always lower case
	l = strings.ToLower(l)

	// Current State
	current := m.locale

	// New State
	m.locale = l
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
	}{
		Version: m.version,
		ID:      m.id,
		Type:    m.mtype,
		SubType: m.msubtype,
		Params:  m.params,
		Created: m.created,
		Queue:   queue,
		Email:   email,
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *EmailMessage) UnmarshalJSON(b []byte) error {
	me := &struct {
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
	}{}

	err := json.Unmarshal(b, &me)
	if err != nil {
		return err
	}

	// Basic Message Information
	m.version = me.Version
	m.id = me.ID
	m.mtype = me.Type
	m.msubtype = me.Type
	m.params = me.Params
	m.created = me.Created

	// QUEUE Message Control Information //
	if me.Queue != nil {
		m.requeueCount = me.Queue.RequeueCount

		// Has Error Message?
		if me.Queue.ErrorCode > 0 { // YES
			m.errorCode = me.Queue.ErrorCode
			m.errorTime = me.Queue.ErrorTime
			m.errorMessage = me.Queue.ErrorMessage
		}
	}

	// EMAIL Message Parameters //
	m.template = me.Email.Template
	m.locale = me.Email.Locale
	m.to = me.Email.To
	m.from = me.Email.From
	m.cc = me.Email.CC
	m.bcc = me.Email.BCC
	m.headers = me.Email.Headers

	return nil
}

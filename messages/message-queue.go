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
	"strconv"
	"strings"
	"time"

	"github.com/objectvault/queue-interface/shared"
)

type QueueMessage struct {
	version      int                     // [REQUIRED] Message Version
	id           string                  // [REQUIRED] Message ID
	mtype        string                  // [REQUIRED] Message Type
	msubtype     string                  // [OPTIONAL] Message Sub Type
	params       *map[string]interface{} // [OPTIONAL] Optional Context Parameters
	created      string                  // [REQUIRED] Original Message Creation TimeStamp
	requeueCount int                     // Number of Times Message Requeued
	errorCode    int                     // Error Code : 0 OK
	errorTime    string                  // Error Time Stamp
	errorMessage string                  // Error Message
}

func NewQueueMessage(t string, st string) (*QueueMessage, error) {
	m := &QueueMessage{}
	err := InitQueueMessage(m, t, st)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func InitQueueMessage(m *QueueMessage, t string, st string) error {
	m.version = 1

	// TODO: Need Some Sort of Prefix to Avoid Collisions if in a Cluster
	// Current TimeStamp
	now := time.Now().UTC()
	m.id = now.Format(time.RFC3339)

	// Validate Message Type
	t = strings.TrimSpace(t)
	if t == "" {
		return errors.New("[QueueMessage] Missing Message Type")
	}
	m.mtype = strings.ToLower(t)

	// Set Sub Type
	st = strings.TrimSpace(st)
	m.msubtype = strings.ToLower(st)

	return nil
}

func (m *QueueMessage) IsValid() bool {
	return (m.id != "") && (m.mtype != "")
}

func (m *QueueMessage) Version() int {
	if m.version > 0 {
		return m.version
	}

	return 1
}

func (m *QueueMessage) SetVersion(v int) (int, error) {
	// Valid Version?
	if v <= 0 { // NO
		return 0, errors.New("[QueueMessage] Invalid Message Version")
	}

	// Current State
	current := m.version

	// New State
	m.version = v
	return current, nil
}

func (m *QueueMessage) ID() string {
	return m.id
}

func (m *QueueMessage) SetID(id string) (string, error) {
	// Is ID Empty?
	id = strings.TrimSpace(id)
	if id == "" { // YES: Error
		return "", errors.New("[QueueMessage] Request ID is Required")
	}

	// IDs are always lower case
	id = strings.ToLower(id)

	// Current State
	current := m.id

	// New State
	m.id = id
	return current, nil
}

func (m *QueueMessage) Type() string {
	return m.mtype
}

func (m *QueueMessage) SetType(t string) (string, error) {
	// Is ID Empty?
	t = strings.TrimSpace(t)
	if t == "" { // YES: Error
		return "", errors.New("[QueueMessage] Request Type is Required")
	}

	// Message Types are always lower case
	t = strings.ToLower(t)

	// Current State
	current := m.mtype

	// New State
	m.mtype = t
	return current, nil
}

func (m *QueueMessage) SubType() string {
	return m.mtype
}

func (m *QueueMessage) SetSubType(s string) (string, error) {
	// Is Sub Type Empty?
	s = strings.TrimSpace(s)
	if s != "" { // NO: Message Sub Types are always lower case
		s = strings.ToLower(s)
	}

	// Current State
	current := m.msubtype

	// New State
	m.msubtype = s
	return current, nil
}

func (m *QueueMessage) GetParameters() *map[string]interface{} {
	return m.params
}

func (m *QueueMessage) HasParameter(n string) bool {
	_, ok := (*m.params)[n]
	return ok
}

func (m *QueueMessage) Parameter(n string) interface{} {
	if m.params == nil {
		return ""
	}

	value, ok := (*m.params)[n]
	if ok {
		return value
	}
	return ""
}

func (m *QueueMessage) SetParameter(n string, v interface{}) error {
	// Do we already have a Parameters Map?
	if m.params == nil { // NO: Create One
		d := make(map[string]interface{})
		m.params = &d
	}

	// New State
	(*m.params)[n] = v
	return nil
}

func (m *QueueMessage) SetIntParameter(n string, v int) error {
	return m.SetParameter(n, strconv.Itoa(v))
}

func (m *QueueMessage) UnsetParameter(n string) (interface{}, error) {
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

func (m *QueueMessage) Created() *time.Time {
	if m.created != "" {
		created, _ := time.Parse(time.RFC3339, m.created)
		return &created
	}

	return nil
}

func (m *QueueMessage) RequeueCount() int {
	return m.requeueCount
}

func (m *QueueMessage) ResetCount() int {
	current := m.requeueCount
	m.requeueCount = 0
	return current
}

func (m *QueueMessage) Requeue() int {
	m.requeueCount++
	return m.requeueCount
}

func (m *QueueMessage) SetError(c int, msg string) error {
	// Valid Error Code?
	if c > 0 { // NO
		return errors.New("[QueueMessage] Invalid Error Code")
	}

	// Valid Error Message?
	if msg == "" { // NO
		return errors.New("[QueueMessage] Missing Error Message")
	}

	m.errorCode = c
	m.errorMessage = msg
	m.errorTime = shared.UTCTimeStamp()
	return nil
}

func (m *QueueMessage) ErrorCode() int {
	return m.errorCode
}

func (m *QueueMessage) ErrorMessage() string {
	return m.errorMessage
}

func (m *QueueMessage) ErrorTime() *time.Time {
	if m.errorTime != "" {
		timestamp, _ := time.Parse(time.RFC3339, m.errorTime)
		return &timestamp
	}

	return nil
}

func (m *QueueMessage) IsError() bool {
	return (m.errorCode > 0)
}

// MarshalJSON implements json.Marshal
func (m QueueMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[QueueMessage] Message is Invalid")
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

	// Complete JSON Message //
	output := &struct {
		Version int                     `json:"version"`
		ID      string                  `json:"id"`
		Type    string                  `json:"type"`
		SubType string                  `json:"subtype,omitempty"`
		Params  *map[string]interface{} `json:"params,omitempty"`
		Created string                  `json:"created"`
		Queue   interface{}             `json:"queue,omitempty"`
	}{
		Version: m.version,
		ID:      m.id,
		Type:    m.mtype,
		SubType: m.msubtype,
		Params:  m.params,
		Created: m.created,
		Queue:   queue,
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *QueueMessage) UnmarshalJSON(b []byte) error {
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

	return nil
}

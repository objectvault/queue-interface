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

type QueueMessage struct {
	version      int          // [REQUIRED] Message Version
	id           string       // [REQUIRED] Message ID
	message      *interface{} // [REQUIRED] Actual Message Data
	created      string       // [REQUIRED] Original Message Creation TimeStamp
	requeueCount int          // Number of Times Message Requeued
	errorCode    int          // Error Code : 0 OK
	errorTime    string       // Error Time Stamp
	errorMessage string       // Error Message
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

func (m *QueueMessage) Message() *interface{} {
	return m.message
}

func (m *QueueMessage) SetMessage(msg interface{}) (interface{}, error) {
	// Is Template Empty?
	if msg == nil {
		return "", errors.New("[QueueMessage] Message Data is Required")
	}

	// Current State
	current := m.message

	// New State
	m.message = &msg
	m.created = shared.UTCTimeStamp()

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

func (m *QueueMessage) IsValid() bool {
	return (m.id != "") && (m.message != nil)
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

	output := &struct {
		Version      int          `json:"version"`
		ID           string       `json:"id"`
		Message      *interface{} `json:"message"`
		Created      string       `json:"created"`
		RequeueCount int          `json:"count,omitempty"`
		ErrorCode    int          `json:"errorcode,omitempty"`
		ErrorTime    string       `json:"errortime,omitempty"`
		ErrorMessage string       `json:"errormsg,omitempty"`
	}{
		Version: m.version,
		ID:      m.id,
		Message: m.message,
		Created: m.created,
	}

	// Has the Message been Requeued?
	if m.requeueCount > 0 { // YES
		output.RequeueCount = m.requeueCount
	}

	// Is this an Error Message?
	if m.errorCode > 0 { // YES
		output.ErrorCode = m.errorCode
		output.ErrorTime = m.errorTime
		output.ErrorMessage = m.errorMessage
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *QueueMessage) UnmarshalJSON(b []byte) error {
	me := &struct {
		Version      int          `json:"version"`
		ID           string       `json:"id"`
		Message      *interface{} `json:"message"`
		Created      string       `json:"created"`
		RequeueCount int          `json:"count,omitempty"`
		ErrorCode    int          `json:"errorcode,omitempty"`
		ErrorTime    string       `json:"errortime,omitempty"`
		ErrorMessage string       `json:"errormsg,omitempty"`
	}{}

	err := json.Unmarshal(b, &me)
	if err != nil {
		return err
	}

	// Basic Message Information
	m.version = me.Version
	m.id = me.ID
	m.message = me.Message
	m.created = me.Created
	m.requeueCount = me.RequeueCount

	// Is Error Message?
	if me.ErrorCode > 0 { // YES
		m.errorCode = me.ErrorCode
		m.errorTime = me.ErrorTime
		m.errorMessage = me.ErrorMessage
	}
	return nil
}

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

// cSpell:ignore atype, gofrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/objectvault/queue-interface/shared"
)

type QueueAction struct {
	version      int                     // [REQUIRED] Action Format Version
	guid         string                  // [REQUIRED] Action GUID
	atype        string                  // [REQUIRED] Action Type
	params       *map[string]interface{} // [OPTIONAL] Action Control Parameters
	props        *map[string]interface{} // [OPTIONAL] Action Context Properties
	created      string                  // [REQUIRED] Original Message Creation TimeStamp
	requeueCount int                     // Number of Times Message Requeued
	errorCode    int                     // Error Code : 0 OK
	errorTime    string                  // Error Time Stamp
	errorMessage string                  // Error Message
}

func NewQueueActionWithGUID(guid string, t string) (*QueueAction, error) {
	// Validate Message Type
	guid = strings.TrimSpace(t)
	if guid == "" {
		return nil, errors.New("[QueueAction] Missing Message GUID")
	}

	// Validate Message Type
	t = strings.TrimSpace(t)
	if t == "" {
		return nil, errors.New("[QueueAction] Missing Message Type")
	}

	m := &QueueAction{
		guid:  strings.ToLower(guid),
		atype: strings.ToLower(t),
	}

	return m, nil
}

func NewQueueAction(t string) (*QueueAction, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to Generate Action Message ID [%v]", err))
	}

	return NewQueueActionWithGUID(uid.String(), t)
}

func (m *QueueAction) IsValid() bool {
	return (m.guid != "") && (m.atype != "")
}

func (m *QueueAction) Version() int {
	if m.version > 0 {
		return m.version
	}

	return 1
}

func (m *QueueAction) ID() string {
	return m.guid
}

func (m *QueueAction) Type() string {
	return m.atype
}

func (m *QueueAction) SetType(t string) (string, error) {
	// Is ID Empty?
	t = strings.TrimSpace(t)
	if t == "" { // YES: Error
		return "", errors.New("[QueueAction] Request Type is Required")
	}

	// Message Types are always lower case
	t = strings.ToLower(t)

	// Current State
	current := m.atype

	// New State
	m.atype = t
	return current, nil
}

func (m *QueueAction) GetParameters() *map[string]interface{} {
	return m.params
}

func (m *QueueAction) SetParameters(p *map[string]interface{}) (*map[string]interface{}, error) {
	// Current State
	current := m.params

	// New State
	m.params = p
	return current, nil
}

func (m *QueueAction) GetProperties() *map[string]interface{} {
	return m.props
}

func (m *QueueAction) SetProperties(p *map[string]interface{}) (*map[string]interface{}, error) {
	// Current State
	current := m.props

	// New State
	m.props = p
	return current, nil
}

func (m *QueueAction) Created() *time.Time {
	if m.created != "" {
		created, _ := time.Parse(time.RFC3339, m.created)
		return &created
	}

	return nil
}

func (m *QueueAction) RequeueCount() int {
	return m.requeueCount
}

func (m *QueueAction) ResetCount() int {
	current := m.requeueCount
	m.requeueCount = 0
	return current
}

func (m *QueueAction) Requeue() int {
	m.requeueCount++
	return m.requeueCount
}

func (m *QueueAction) SetError(c int, msg string) error {
	// Valid Error Code?
	if c > 0 { // NO
		return errors.New("[QueueAction] Invalid Error Code")
	}

	// Valid Error Message?
	if msg == "" { // NO
		return errors.New("[QueueAction] Missing Error Message")
	}

	m.errorCode = c
	m.errorMessage = msg
	m.errorTime = shared.UTCTimeStamp()
	return nil
}

func (m *QueueAction) ErrorCode() int {
	return m.errorCode
}

func (m *QueueAction) ErrorMessage() string {
	return m.errorMessage
}

func (m *QueueAction) ErrorTime() *time.Time {
	if m.errorTime != "" {
		timestamp, _ := time.Parse(time.RFC3339, m.errorTime)
		return &timestamp
	}

	return nil
}

func (m *QueueAction) IsError() bool {
	return (m.errorCode > 0)
}

// MarshalJSON implements json.Marshal
func (m QueueAction) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[QueueAction] Message is Invalid")
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
		Params  *map[string]interface{} `json:"params,omitempty"`
		Props   *map[string]interface{} `json:"props,omitempty"`
		Created string                  `json:"created"`
		Queue   interface{}             `json:"queue,omitempty"`
	}{
		Version: m.version,
		ID:      m.guid,
		Type:    m.atype,
		Params:  m.params,
		Props:   m.props,
		Created: m.created,
		Queue:   queue,
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *QueueAction) UnmarshalJSON(b []byte) error {
	me := &struct {
		Version int                     `json:"version"`
		ID      string                  `json:"id"`
		Type    string                  `json:"type"`
		Params  *map[string]interface{} `json:"params,omitempty"`
		Props   *map[string]interface{} `json:"props,omitempty"`
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
	m.guid = me.ID
	m.atype = me.Type
	m.params = me.Params
	m.props = me.Props
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

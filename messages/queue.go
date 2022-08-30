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

// cSpell:ignore mtype
import (
	"errors"
	"strings"
	"time"
)

type QueueMessage struct {
	id           string     // [REQUIRED] Message ID (Preferably a GUID)
	mtype        string     // [REQUIRED] Message Type
	created      *time.Time // [REQUIRED] Original Message Creation TimeStamp
	requeueCount int        // Number of Times Message Requeued
	errorCode    int        // Error Code : 0 OK
	errorTime    *time.Time // Error Time Stamp
	errorMessage string     // Error Message
}

func InitQueueMessage(m QueueMessage, id string, t string) error {
	// Validate Message Type
	id = strings.TrimSpace(t)
	if id == "" {
		return errors.New("[QueueMessage] Missing Message ID")
	}

	// Validate Message Type
	t = strings.TrimSpace(t)
	if t == "" {
		return errors.New("[QueueMessage] Missing Message Type")
	}

	m.id = strings.ToLower(id)
	m.mtype = strings.ToLower(t)
	// m.created Should be Set Closer to When Message is Queued
	return nil
}

func (m *QueueMessage) IsValid() bool {
	return (m.id != "") && (m.mtype != "")
}

func (m *QueueMessage) ID() string {
	return m.id
}

func (m *QueueMessage) SetID(id string) error {
	// Is ID Empty?
	id = strings.TrimSpace(id)
	if id == "" { // YES: Error
		return errors.New("[QueueMessage] Request ID is Required")
	}

	// IDs are always lower case
	id = strings.ToLower(id)

	// New State
	m.id = id
	return nil
}

func (m *QueueMessage) Type() string {
	return m.mtype
}

func (m *QueueMessage) SetType(t string) error {
	// Is ID Empty?
	t = strings.TrimSpace(t)
	if t == "" { // YES: Error
		return errors.New("[QueueAction] Request Type is Required")
	}

	// Message Types are always lower case
	t = strings.ToLower(t)

	// New State
	m.mtype = t
	return nil
}

func (m *QueueMessage) Created() *time.Time {
	return m.created
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
	t := time.Now().UTC()
	m.errorTime = &t
	return nil
}

func (m *QueueMessage) ErrorCode() int {
	return m.errorCode
}

func (m *QueueMessage) ErrorMessage() string {
	return m.errorMessage
}

func (m *QueueMessage) ErrorTime() *time.Time {
	return m.errorTime
}

func (m *QueueMessage) IsError() bool {
	return (m.errorCode > 0)
}

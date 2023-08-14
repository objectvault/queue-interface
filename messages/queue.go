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
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/objectvault/common/maps"
)

// Current Message Processing Status
type QueueMessageStatus struct {
	errorCode        int             // [REQUIRED] Error Code (0 = OK)
	errorMessage     string          // [OPTIONAL] Error Message Text
	errorMessageI18N string          // [OPTIONAL] Error Message I18N Code
	extras           maps.MapWrapper // [OPTIONAL] Optional Information
}

// Constructor
func NewQueueMessageStatus() *QueueMessageStatus {
	o := &QueueMessageStatus{
		errorCode: 0,
	}

	return o
}

func (o *QueueMessageStatus) InError() bool {
	return (o.errorCode != 0)
}

func (o *QueueMessageStatus) ErrorCode() int {
	return o.errorCode
}

func (o *QueueMessageStatus) ErrorMessage() string {
	return o.errorMessage
}

func (o *QueueMessageStatus) SetError(code int, en string, i18n string) {
	o.errorCode = code
	o.errorMessage = strings.TrimSpace(en)
	o.errorMessageI18N = strings.TrimSpace(i18n)
}

func (o *QueueMessageStatus) Extras() map[string]interface{} {
	return o.extras.Map()
}

func (o *QueueMessageStatus) MarshalJSON() ([]byte, error) {
	// Convert to JSON
	return json.Marshal(&struct {
		ErrorCode        int         `json:"error_code"`
		ErrorMessage     string      `json:"error_message,omitempty"`
		ErrorMessageI18N string      `json:"error_message_i18n,omitempty"`
		Extras           interface{} `json:"extras,omitempty"`
	}{
		ErrorCode:        o.errorCode,
		ErrorMessage:     o.errorMessage,
		ErrorMessageI18N: o.errorMessageI18N,
		Extras:           o.extras.Map(),
	})
}

type QueueMessageHeader struct {
	version int                 // [REQUIRED] Message Version
	id      string              // [REQUIRED] Message ID (Preferably a GUID)
	parent  string              // [OPTIONAL] Associated Parent Message ID
	props   maps.MapWrapper     // [OPTIONAL] Message Processing Properties
	status  *QueueMessageStatus // [OPTIONAL] Message Processing Status
	created *time.Time          // [OPTIONAL] Message Creation Date
}

// Constructor
func NewQueueMessageHeader(id string, parent string) *QueueMessageHeader {
	o := &QueueMessageHeader{
		version: 1,
	}

	o.SetID(id)
	o.SetParent(parent)

	return o
}

func (o *QueueMessageHeader) IsValid() bool {
	return (o.version > 0) && (o.id != "")
}

func (o *QueueMessageHeader) Version() int {
	return o.version
}

func (o *QueueMessageHeader) ID() string {
	return o.id
}

func (o *QueueMessageHeader) SetID(id string) {
	o.id = strings.TrimSpace(id)
	if o.id != "" {
		o.id = strings.ToLower(o.id)
	}
}

func (o *QueueMessageHeader) Parent() string {
	return o.parent
}

func (o *QueueMessageHeader) SetParent(id string) {
	o.parent = strings.TrimSpace(id)
	if o.parent != "" {
		o.parent = strings.ToLower(o.parent)
	}
}

func (o *QueueMessageHeader) SetProperties(m map[string]interface{}) {
	o.props = *maps.NewMapWrapper(m)
}

func (o *QueueMessageHeader) Status() *QueueMessageStatus {
	return o.status
}

func (o *QueueMessageHeader) Created() time.Time {
	if o.created == nil {
		now := time.Now().UTC()
		o.created = &now
	}

	return *o.created
}

func (o *QueueMessageHeader) MarshalJSON() ([]byte, error) {
	if !o.IsValid() {
		return nil, errors.New("[QueueMessageHeader] Is not valid")
	}

	// Convert to JSON
	j := &struct {
		Version int         `json:"version"`
		ID      string      `json:"id"`
		Parent  string      `json:"parent,omitempty"`
		Props   interface{} `json:"props,omitempty"`
		Status  interface{} `json:"status,omitempty"`
		Created time.Time   `json:"created"`
	}{
		Version: o.version,
		ID:      o.id,
		Parent:  o.parent,
		Created: o.Created(),
	}

	// Properties Set?
	if !o.props.IsEmpty() {
		j.Props = o.props.Map()
	}

	// Status Set?
	if o.status != nil {
		j.Status = o.status
	}

	// Convert Structure to JSON
	return json.Marshal(j)
}

type QueueMessage struct {
	header *QueueMessageHeader // [REQUIRED] Message Header
	body   interface{}         // [REQUIRED] Message Content
}

func NewQueueMessage(id string, message interface{}) *QueueMessage {
	o := &QueueMessage{
		header: NewQueueMessageHeader(id, ""),
		body:   message,
	}

	return o
}

func (o *QueueMessage) IsValid() bool {
	return (o.header != nil) && o.header.IsValid() && (o.body != nil)
}

func (o *QueueMessage) Header() *QueueMessageHeader {
	if o.header == nil {
		o.header = &QueueMessageHeader{}
	}

	return o.header
}

func (o *QueueMessage) Message() interface{} {
	return o.body
}

func (o *QueueMessage) SetMessage(message interface{}) {
	o.body = message
}

func (o *QueueMessage) MarshalJSON() ([]byte, error) {
	if !o.IsValid() {
		return nil, errors.New("[QueueMessage] Is not valid")
	}

	// Convert to JSON
	return json.Marshal(&struct {
		Header  interface{} `json:"header"`
		Message interface{} `json:"body"`
	}{
		Header:  o.header,
		Message: o.body,
	})
}

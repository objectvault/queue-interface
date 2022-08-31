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

// cSpell:ignore gofrs, mtype
import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/objectvault/common/maps"
	"github.com/objectvault/queue-interface/shared"
)

type ActionMessage struct {
	QueueMessage
	version int                    // [REQUIRED] Action Format Version
	params  map[string]interface{} // [OPTIONAL] Action Control Parameters
	props   map[string]interface{} // [OPTIONAL] Action Context Properties
}

func NewQueueActionWithGUID(guid string, t string) (*ActionMessage, error) {
	a := &ActionMessage{}

	// Initialize Message
	err := InitQueueAction(a, guid, t)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func NewQueueAction(t string) (*ActionMessage, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[ActionMessage] Failed to Generate Action Message ID [%v]", err)
	}

	return NewQueueActionWithGUID(uid.String(), t)
}

func InitQueueAction(a *ActionMessage, guid string, t string) error {
	// Validate Message Type
	t = strings.TrimSpace(t)
	if t == "" {
		return errors.New("[ActionMessage] Missing Message Type")
	}

	// Set Default Version
	a.version = 1
	return InitQueueMessage(&(a.QueueMessage), guid, "action:"+t)
}

func (m *ActionMessage) Version() int {
	if m.version > 0 {
		return m.version
	}

	return 1
}

func (m *ActionMessage) SetVersion(v int) error {
	// Valid Version?
	if v <= 0 { // NO
		return errors.New("[ActionMessage] Invalid Message Version")
	}

	// New State
	m.version = v
	return nil
}

func (m *ActionMessage) HasParameter(path string) bool {
	return maps.Has(m.params, path)
}

func (m *ActionMessage) GetParameter(path string) (interface{}, error) {
	return maps.Get(m.params, path)
}

func (o *ActionMessage) SetParameter(path string, v interface{}, force bool) error {
	m, e := maps.Set(o.params, path, v, force)
	if e == nil {
		o.params = m
	}
	return e
}

func (o *ActionMessage) ClearParameter(path string) error {
	m, e := maps.Clear(o.params, path)
	if e == nil {
		o.params = m
	}
	return e
}

func (m *ActionMessage) GetParameters() map[string]interface{} {
	return m.params
}

func (m *ActionMessage) SetParameters(p map[string]interface{}) error {
	m.params = p
	return nil
}

func (m *ActionMessage) HasProperty(path string) bool {
	return maps.Has(m.props, path)
}

func (m *ActionMessage) GetProperty(path string) (interface{}, error) {
	return maps.Get(m.props, path)
}

func (o *ActionMessage) SetProperty(path string, v interface{}, force bool) error {
	m, e := maps.Set(o.props, path, v, force)
	if e == nil {
		o.props = m
	}
	return e
}

func (o *ActionMessage) ClearProperty(path string) error {
	m, e := maps.Clear(o.props, path)
	if e == nil {
		o.props = m
	}
	return e
}

func (m *ActionMessage) GetProperties() map[string]interface{} {
	return m.props
}

func (m *ActionMessage) SetProperties(p map[string]interface{}) error {
	m.props = p
	return nil
}

// MarshalJSON implements json.Marshal
func (m ActionMessage) MarshalJSON() ([]byte, error) {
	// Is Message Valid?
	if !m.IsValid() { // NO
		return nil, errors.New("[ActionMessage] Message is Invalid")
	}

	// Is Message Creation Date Set?
	if m.created == nil { // NO: Use Current Time
		t := time.Now().UTC()
		m.created = &t
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
		queue.ErrorTime = shared.ToJSONTimeStamp(m.errorTime)
		queue.ErrorMessage = m.errorMessage
	}

	// Complete JSON Message //
	output := &struct {
		Version int                    `json:"version"`
		ID      string                 `json:"id"`
		Type    string                 `json:"type"`
		Params  map[string]interface{} `json:"params,omitempty"`
		Props   map[string]interface{} `json:"props,omitempty"`
		Created string                 `json:"created"`
		Queue   interface{}            `json:"queue,omitempty"`
	}{
		Version: m.version,
		ID:      m.id,
		Type:    m.mtype,
		Params:  m.params,
		Props:   m.props,
		Created: shared.ToJSONTimeStamp(m.created),
		Queue:   queue,
	}

	return json.Marshal(output)
}

// UnmarshalJSON implements json.Unmarshal
func (m *ActionMessage) UnmarshalJSON(b []byte) error {
	me := &struct {
		Version int                    `json:"version"`
		ID      string                 `json:"id"`
		Type    string                 `json:"type"`
		Params  map[string]interface{} `json:"params,omitempty"`
		Props   map[string]interface{} `json:"props,omitempty"`
		Created string                 `json:"created"`
		Queue   *struct {
			RequeueCount int    `json:"count,omitempty"`
			ErrorCode    int    `json:"errorcode,omitempty"`
			ErrorTime    string `json:"errortime,omitempty"`
			ErrorMessage string `json:"errormsg,omitempty"`
		} `json:"queue,omitempty"`
	}{}

	err := json.Unmarshal(b, &me)
	if err != nil {
		return err
	}

	// Basic Message Information
	m.version = me.Version
	m.id = me.ID
	m.mtype = me.Type
	m.params = me.Params
	m.props = me.Props
	m.created = shared.FromJSONTimeStamp(me.Created)

	// QUEUE Message Control Information //
	if me.Queue != nil {
		m.requeueCount = me.Queue.RequeueCount

		// Has Error Message?
		if me.Queue.ErrorCode > 0 { // YES
			m.errorCode = me.Queue.ErrorCode
			m.errorTime = shared.FromJSONTimeStamp(me.Queue.ErrorTime)
			m.errorMessage = me.Queue.ErrorMessage
		}
	}

	return nil
}

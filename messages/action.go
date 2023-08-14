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

// cSpell:ignore gofrs, atype
import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/objectvault/common/maps"
)

type ActionMessageContent struct {
	atype  string          // [REQUIRED] Action Type
	params maps.MapWrapper // [OPTIONAL] Action Control Parameters
	props  maps.MapWrapper // [OPTIONAL] Action Context Properties
}

func NewActionMessageContent(t string) *ActionMessageContent {
	o := &ActionMessageContent{}
	o.SetType(t)
	return o
}

func (o *ActionMessageContent) IsValid() bool {
	return (o.atype != "")
}

func (o *ActionMessageContent) Type() string {
	return o.atype
}

func (o *ActionMessageContent) SetType(t string) {
	o.atype = strings.TrimSpace(t)
	if o.atype != "" {
		o.atype = strings.ToLower(o.atype)
	}
}

func (o *ActionMessageContent) SetParameters(m map[string]interface{}) {
	o.params = *maps.NewMapWrapper(m)
}

func (o *ActionMessageContent) SetProperties(m map[string]interface{}) {
	o.props = *maps.NewMapWrapper(m)
}

func (o *ActionMessageContent) MarshalJSON() ([]byte, error) {
	if !o.IsValid() {
		return nil, errors.New("[ActionMessageContent] Is not valid")
	}

	// JSON Structure
	j := &struct {
		Type   string      `json:"type"`
		Params interface{} `json:"params,omitempty"`
		Props  interface{} `json:"props,omitempty"`
	}{
		Type: o.atype,
	}

	// Parameters Set?
	if !o.params.IsEmpty() {
		j.Params = o.params.Map()
	}

	// Properties Set?
	if !o.props.IsEmpty() {
		j.Props = o.props.Map()
	}

	// Convert Structure to JSON
	return json.Marshal(j)
}

type ActionMessage struct {
	QueueMessage
}

func NewQueueActionMessage(t string) (*ActionMessage, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[ActionMessage] Failed to Generate Action Message ID [%v]", err)
	}

	return NewQueueActionWithGUID(uid.String(), t)
}

func NewQueueActionWithGUID(guid string, t string) (*ActionMessage, error) {
	o := &ActionMessage{}

	// Initialize Action Message
	err := InitQueueAction(o, guid, t)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func InitQueueAction(o *ActionMessage, guid string, t string) error {
	// Initialize Header
	o.header = NewQueueMessageHeader(guid, "")

	// Initialize Content
	t = strings.TrimSpace(t)
	if t != "" {
		t = "action:" + t
	}

	o.QueueMessage.SetMessage(NewActionMessageContent(t))

	return nil
}

func GetActionMessageContent(o *ActionMessage) *ActionMessageContent {
	m := o.QueueMessage.Message()
	if m != nil {
		c, ok := m.(*ActionMessageContent)
		if ok {
			return c
		}
	}

	return nil
}

func (o *ActionMessage) IsValid() bool {
	if (o.header != nil) && o.header.IsValid() {
		c := GetActionMessageContent(o)
		if c != nil {
			return c.IsValid()
		}
	}

	return false
}

func (o *ActionMessage) Params() *maps.MapWrapper {
	c := GetActionMessageContent(o)
	if c != nil {
		return &c.params
	}

	return nil
}

func (o *ActionMessage) SetParameters(m map[string]interface{}) error {
	c := GetActionMessageContent(o)
	if c != nil {
		c.SetParameters(m)
		return nil
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

func (o *ActionMessage) Props() *maps.MapWrapper {
	c := GetActionMessageContent(o)
	if c != nil {
		return &c.props
	}

	return nil
}

func (o *ActionMessage) SetProperties(m map[string]interface{}) error {
	c := GetActionMessageContent(o)
	if c != nil {
		c.SetProperties(m)
		return nil
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

func (o *ActionMessage) SetParameter(path string, v interface{}) error {
	p := o.Params()
	if p != nil {
		return p.Set(path, v, true)
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

func (o *ActionMessage) SetStringParameter(path string, s string, clear bool) error {
	p := o.Params()
	if p != nil {
		if s == "" && clear {
			return p.Clear(p)
		}

		return p.Set(path, s, true)
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

func (o *ActionMessage) SetProperty(path string, v interface{}) error {
	p := o.Props()
	if p != nil {
		return p.Set(path, v, true)
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

func (o *ActionMessage) SetStringProperty(path string, s string, clear bool) error {
	// Set Parameter
	p := o.Props()
	if p != nil {
		if s == "" && clear {
			return p.Clear(p)
		}

		return p.Set(path, s, true)
	}

	return errors.New("[ActionMessage] Initialize Message before using")
}

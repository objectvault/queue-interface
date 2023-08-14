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

// cSpell:ignore gofrs, objectname, storename
import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/objectvault/queue-interface/shared"
)

type InviteMessage struct {
	EmailMessage // DERIVED FROM
}

func NewInviteMessageWithGUID(guid string, ot string, code string) (*InviteMessage, error) {
	m := &InviteMessage{}
	err := InitInviteMessage(m, guid, ot, code)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func NewInviteMessage(ot string, code string) (*InviteMessage, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[InviteMessage] Failed to Generate Action Message ID [%v]", err)
	}

	return NewInviteMessageWithGUID(uid.String(), ot, code)
}

func InitInviteMessage(m *InviteMessage, guid string, ot string, code string) error {
	ot = strings.TrimSpace(ot)
	if ot == "" {
		return errors.New("[InviteMessage] Invitation Object Type Required")
	}

	// Initialize Email Message
	err := InitEmailMessage(&(m.EmailMessage), guid, "invite:"+strings.ToLower(ot), "")
	if err != nil {
		return err
	}

	// Set Message Code
	err = m.SetCode(code)
	if err != nil {
		return err
	}

	return nil
}

func (m *InviteMessage) IsValid() bool {
	return m.EmailMessage.IsValid() && (m.Code() != "") && (m.ByUser() != "") && (m.ObjectName() != "") && (m.Expiration() != nil)
}

func (m *InviteMessage) Code() string {
	p := m.Params()
	if p != nil {
		code, e := p.GetDefault("code", "")
		if e == nil {
			return code.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetCode(code string) error {
	// Is Invitation Activation Code Empty?
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.New("[InviteMessage] Invitation Code is Required")
	}

	return m.SetProperty("code", strings.ToLower(code))
}

func (m *InviteMessage) ByUser() string {
	p := m.Params()
	if p != nil {
		name, e := p.GetDefault("by-name", "")
		if e == nil {
			return name.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetByUser(name string) error {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("[InviteMessage] From User Name is Required")
	}

	return m.SetProperty("by-name", name)
}

func (m *InviteMessage) ByEmail() string {
	p := m.Params()
	if p != nil {
		email, e := p.GetDefault("by-email", "")
		if e == nil {
			return email.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetByEmail(email string) error {
	// Is Email Empty?
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("[InviteMessage] From User Email is Required")
	}

	return m.SetProperty("by-email", strings.ToLower(email))
}

func (m *InviteMessage) Message() string {
	p := m.Params()
	if p != nil {
		email, e := p.GetDefault("message", "")
		if e == nil {
			return email.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetMessage(msg string) error {
	return m.SetStringProperty("message", msg, true)
}

func (m *InviteMessage) ObjectName() string {
	p := m.Params()
	if p != nil {
		name, e := p.GetDefault("objectname", "")
		if e == nil {
			return name.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetObjectName(name string) error {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("[InviteMessage] Object Name is Required")
	}

	return m.SetProperty("objectname", name)
}

func (m *InviteMessage) StoreName() string {
	p := m.Params()
	if p != nil {
		name, e := p.GetDefault("storename", "")
		if e == nil {
			return name.(string)
		}
	}

	return ""
}

func (m *InviteMessage) SetStoreName(name string) error {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("[InviteMessage] Object Name is Required")
	}

	return m.SetProperty("storename", name)
}

func (m *InviteMessage) Expiration() *time.Time {
	p := m.Params()
	if p != nil {
		t, e := p.Get("expiration")
		if e == nil || t != nil {
			return shared.FromJSONTimeStamp(t.(string))
		}
	}

	return nil
}

func (m *InviteMessage) SetExpiration(t time.Time) error {
	return m.SetProperty("expiration", shared.ToJSONTimeStamp(&t))
}

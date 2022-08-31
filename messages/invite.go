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

// cSpell:ignore gofrs, objectname
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
	code, e := m.GetProperty("code")
	if e != nil || code == nil {
		return ""
	}

	return code.(string)
}

func (m *InviteMessage) SetCode(code string) error {
	// Is Invitation Activation Code Empty?
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.New("[InviteMessage] Invitation Code is Required")
	}

	// Set Property
	return m.SetProperty("code", strings.ToLower(code), true)
}

func (m *InviteMessage) ByUser() string {
	name, e := m.GetProperty("by-name")
	if e != nil || name == nil {
		return ""
	}

	return name.(string)
}

func (m *InviteMessage) SetByUser(name string) error {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("[InviteMessage] From User Name is Required")
	}

	// Set Property
	return m.SetProperty("by-name", name, true)
}

func (m *InviteMessage) ByEmail() string {
	email, e := m.GetProperty("by-email")
	if e != nil || email == nil {
		return ""
	}

	return email.(string)
}

func (m *InviteMessage) SetByEmail(email string) error {
	// Is Email Empty?
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("[InviteMessage] From User Email is Required")
	}

	// Set Property
	return m.SetProperty("by-email", strings.ToLower(email), true)
}

func (m *InviteMessage) Message() string {
	msg, e := m.GetProperty("message")
	if e != nil || msg == nil {
		return ""
	}

	return msg.(string)
}

func (m *InviteMessage) SetMessage(msg string) error {
	// Is Message Empty?
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return m.ClearProperty("message")
	}

	// Set Property
	return m.SetProperty("message", msg, true)
}

func (m *InviteMessage) ObjectName() string {
	n, e := m.GetProperty("objectname")
	if e != nil || n == nil {
		return ""
	}

	return n.(string)
}

func (m *InviteMessage) SetObjectName(name string) error {
	// Is Name Empty?
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("[InviteMessage] Object Name is Required")
	}

	e := m.SetProperty("objectname", name, true)
	if e != nil {
		return e
	}
	return nil
}

func (m *InviteMessage) Expiration() *time.Time {
	t, e := m.GetProperty("expiration")
	if e != nil || t == nil {
		return nil
	}

	return shared.FromJSONTimeStamp(t.(string))
}

func (m *InviteMessage) SetExpiration(t time.Time) error {
	e := m.SetProperty("expiration", shared.ToJSONTimeStamp(&t), true)
	if e != nil {
		return e
	}

	return nil
}

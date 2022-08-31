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
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

type EmailMessage struct {
	ActionMessage // DERIVED FROM
}

func NewEmailMessageWithGUID(guid string, st string, template string) (*EmailMessage, error) {
	m := EmailMessage{}
	err := InitEmailMessage(m, guid, st, template)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func NewEmailMessage(st string, template string) (*EmailMessage, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[EmailMessage] Failed to Generate Action Message ID [%v]", err)
	}

	return NewEmailMessageWithGUID(uid.String(), st, template)
}

func InitEmailMessage(m EmailMessage, guid string, et string, template string) error {
	et = strings.TrimSpace(et)
	template = strings.TrimSpace(template)

	if et == "" {
		if template == "" {
			return errors.New("[EmailMessage] Untyped email requires template")
		}

		et = "email"
	} else {
		et = strings.ToLower(et)
		et = "email:" + et
	}

	// Initialize Action Message
	err := InitQueueAction(m.ActionMessage, guid, et)
	if err != nil {
		return err
	}

	// Save Template (Note: ALLOW template == "")
	template = strings.TrimSpace(template)
	if template != "" {
		m.SetParameter("template", strings.ToLower(template), true)
	}

	return nil
}

func (m *EmailMessage) IsValid() bool {
	return m.IsValid() && m.HasParameter("template") && m.HasParameter("to")
}

func (m *EmailMessage) Template() string {
	t, e := m.GetParameter("template")
	if e != nil {
		return ""
	}
	return t.(string)
}

func (m *EmailMessage) SetTemplate(t string) error {
	// Is Template Name Empty?
	t = strings.TrimSpace(t)
	if t == "" {
		return errors.New("Email Template is Required")
	}

	// Set Parameter
	return m.SetParameter("template", strings.ToLower(t), true)
}

func (m *EmailMessage) Locale() string {
	l, e := m.GetParameter("locale")
	if e != nil || l == nil {
		return "en_us"
	}

	return l.(string)
}

func (m *EmailMessage) SetLocale(l string) error {
	// Is Locale Empty?
	l = strings.TrimSpace(l)
	if l == "" {
		return m.ClearParameter("locale")
	}

	// Set Parameter
	return m.SetParameter("locale", strings.ToLower(l), true)
}

func (m *EmailMessage) To() string {
	to, e := m.GetParameter("to")
	if e != nil {
		return ""
	}

	return to.(string)
}

func (m *EmailMessage) SetTo(to string) error {
	// Is Template Name Empty?
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("Email Destination is Required")
	}

	// Set Parameter
	return m.SetParameter("to", strings.ToLower(to), true)
}

func (m *EmailMessage) From(d string) string {
	f, e := m.GetParameter("from")
	if e != nil || f == nil {
		return d
	}

	return f.(string)
}

func (m *EmailMessage) SetFrom(from string) error {
	// Is FROM Empty?
	from = strings.TrimSpace(from)
	if from == "" {
		return m.ClearParameter("from")
	}

	// Set Parameter
	return m.SetParameter("from", strings.ToLower(from), true)
}

func (m *EmailMessage) CC() string {
	cc, e := m.GetParameter("cc")
	if e != nil || cc == nil {
		return ""
	}

	return cc.(string)
}

func (m *EmailMessage) SetCC(cc string) error {
	// Is CC Empty?
	cc = strings.TrimSpace(cc)
	if cc == "" {
		return m.ClearParameter("cc")
	}

	// Set Parameter
	return m.SetParameter("cc", strings.ToLower(cc), true)
}

func (m *EmailMessage) BCC() string {
	bcc, e := m.GetParameter("bcc")
	if e != nil || bcc == nil {
		return ""
	}

	return bcc.(string)
}

func (m *EmailMessage) SetBCC(bcc string) error {
	// Is BCC Empty?
	bcc = strings.TrimSpace(bcc)
	if bcc == "" {
		return m.ClearParameter("bcc")
	}

	// Set Parameter
	return m.SetParameter("bcc", strings.ToLower(bcc), true)
}

func (m *EmailMessage) GetHeaders() map[string]interface{} {
	h, e := m.GetParameter("headers")
	if e != nil || h == nil {
		return nil
	}

	return h.(map[string]interface{})
}

func (m *EmailMessage) HasHeader(n string) bool {
	return m.HasParameter("headers." + strings.ToLower(n))
}

func (m *EmailMessage) Header(n string) string {
	h, e := m.GetParameter("headers." + strings.ToLower(n))
	if e != nil || h == nil {
		return ""
	}

	return h.(string)
}

func (m *EmailMessage) SetHeader(n string, v string) error {
	return m.SetParameter("headers."+strings.ToLower(n), strings.TrimSpace(v), true)
}

func (m *EmailMessage) ClearHeader(n string) error {
	e := m.ClearParameter("headers." + strings.ToLower(n))
	if e != nil {
		return e
	}

	h := m.GetHeaders()
	if h != nil && len(h) == 0 {
		return m.ClearHeaders()
	}

	return nil
}

func (m *EmailMessage) ClearHeaders() error {
	e := m.ClearParameter("headers")
	if e != nil {
		return e
	}

	return nil
}

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

func NewEmailMessage(st string, template string) (*EmailMessage, error) {
	// Create GUID (V4 see https://www.sohamkamani.com/uuid-versions-explained/)
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[EmailMessage] Failed to Generate Action Message ID [%v]", err)
	}

	return NewEmailMessageWithGUID(uid.String(), st, template)
}

func NewEmailMessageWithGUID(guid string, st string, template string) (*EmailMessage, error) {
	m := &EmailMessage{}
	err := InitEmailMessage(m, guid, st, template)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func InitEmailMessage(m *EmailMessage, guid string, et string, template string) error {
	et = strings.TrimSpace(et)
	template = strings.TrimSpace(template)

	if et == "" {
		if template == "" {
			return errors.New("[EmailMessage] Untyped email requires template")
		}

		et = "email"
	} else {
		et = "email:" + et
	}

	// Initialize Action Message
	err := InitQueueAction(&(m.ActionMessage), guid, et)
	if err != nil {
		return err
	}

	// Save Template (Note: ALLOW template == "")
	template = strings.TrimSpace(template)
	if template != "" {
		m.SetTemplate(strings.ToLower(template))
	}

	return nil
}

func (m *EmailMessage) IsValid() bool {
	return m.IsValid() && (m.Template() != "") && (m.To() != "")
}

func (m *EmailMessage) Template() string {
	p := m.Params()
	if p != nil {
		t, e := p.GetDefault("template", "")
		if e == nil {
			return t.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetTemplate(t string) error {
	// Is Template Name Empty?
	t = strings.TrimSpace(t)
	if t == "" {
		return errors.New("Email Template is Required")
	}

	return m.SetParameter("template", strings.ToLower(t))
}

func (m *EmailMessage) Locale() string {
	p := m.Params()
	if p != nil {
		l, e := p.GetDefault("locale", "en_us")
		if e == nil {
			return l.(string)
		}
	}

	return "en_us"
}

func (m *EmailMessage) SetLocale(l string) error {
	return m.SetStringParameter("template", strings.ToLower(l), true)
}

func (m *EmailMessage) To() string {
	p := m.Params()
	if p != nil {
		to, e := p.GetDefault("to", "")
		if e == nil {
			return to.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetTo(to string) error {
	// Is Template Name Empty?
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("Email Destination is Required")
	}

	return m.SetParameter("to", strings.ToLower(to))
}

func (m *EmailMessage) From(d string) string {
	p := m.Params()
	if p != nil {
		from, e := p.GetDefault("from", "")
		if e == nil {
			return from.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetFrom(from string) error {
	return m.SetStringParameter("from", strings.ToLower(from), true)
}

func (m *EmailMessage) CC() string {
	p := m.Params()
	if p != nil {
		cc, e := p.GetDefault("cc", "")
		if e == nil {
			return cc.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetCC(cc string) error {
	return m.SetStringParameter("cc", strings.ToLower(cc), true)
}

func (m *EmailMessage) BCC() string {
	p := m.Params()
	if p != nil {
		bcc, e := p.GetDefault("bcc", "")
		if e == nil {
			return bcc.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetBCC(bcc string) error {
	return m.SetStringParameter("bcc", strings.ToLower(bcc), true)
}

func (m *EmailMessage) GetHeaders() map[string]interface{} {
	p := m.Params()
	if p != nil {
		h, e := p.Get("headers")
		if e == nil && h != nil {
			return h.(map[string]interface{})
		}
	}

	return nil
}

func (m *EmailMessage) HasHeader(n string) bool {
	p := m.Params()
	if p != nil {
		return p.Has("headers." + strings.ToLower(n))
	}

	return false
}

func (m *EmailMessage) Header(n string) string {
	p := m.Params()
	if p != nil {
		h, e := p.GetDefault("headers."+strings.ToLower(n), "")
		if e == nil {
			return h.(string)
		}
	}

	return ""
}

func (m *EmailMessage) SetHeader(n string, v string) error {
	return m.SetStringParameter("headers."+strings.ToLower(n), strings.TrimSpace(v), true)
}

func (m *EmailMessage) ClearHeader(n string) error {
	p := m.Params()
	if p != nil {
		return p.Clear("headers." + strings.ToLower(n))
	}

	return nil
}

func (m *EmailMessage) ClearHeaders() error {
	p := m.Params()
	if p != nil {
		return p.Clear("headers")
	}

	return nil
}

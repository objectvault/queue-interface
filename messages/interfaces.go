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

import "time"

type IMessage interface {
	IsValid() bool
	ID() string
	Type() string
	Created() *time.Time

	Requeue() int
	RequeueCount() int
	ResetCount() int

	ErrorCode() int
	ErrorMessage() string
	ErrorTime() *time.Time
	IsError() bool
}

type IActionMessage interface {
	IMessage

	Version() int

	HasParameter(path string) bool
	GetParameter(path string) (interface{}, error)
	SetParameter(path string, v interface{}, force bool) error
	ClearParameter(path string) error

	GetParameters() map[string]interface{}
	SetParameters(p map[string]interface{}) error

	HasProperty(path string) bool
	GetProperty(path string) (interface{}, error)
	SetProperty(path string, v interface{}, force bool) error
	ClearProperty(path string) error

	GetProperties() map[string]interface{}
	SetProperties(p map[string]interface{}) error
}

type IEmailMessage interface {
	IActionMessage

	Template() string
	SetTemplate(t string) error
	Locale() string
	SetLocale(l string) error

	To() string
	SetTo(to string) error
	From(d string) string
	SetFrom(from string) error
	CC() string
	SetCC(cc string) error
	BCC() string
	SetBCC(bcc string) error
	HasHeader(n string) bool
	Header(n string) string
	SetHeader(n string, v string) error
	ClearHeader(n string) error
	ClearHeaders() error
}

type IInviteEmailMessage interface {
	IEmailMessage

	Code() string
	SetCode(code string) error
	ByUser() string
	SetByUser(name string) error
	ByEmail() string
	SetByEmail(email string) error
	Message() string
	SetMessage(msg string) error
	ObjectName() string
	SetObjectName(name string) error
	Expiration() *time.Time
	SetExpiration(t time.Time) error
}

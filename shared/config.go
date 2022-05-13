package shared

/*
 * This file is part of the ObjectVault Project.
 * Copyright (C) 2020-2022 Paulo Ferreira <vault at sourcenotes.org>
 *
 * This work is published under the GNU AGPLv3.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"encoding/json"
	"errors"
	"time"
)

// TYPE DEFINTION FOR CONFIG FILE //
type Server struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type AMQPConnection struct {
	User     string                 `json:"user,omitempty"`
	Password string                 `json:"password,omitempty"`
	Server   *Server                `json:"server,omitempty"`
	VHost    string                 `json:"vhost,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type Queue struct {
	Servers     []AMQPConnection `json:"servers,omitempty"` // List of AMQP Servers
	QueuePrefix string           `json:"prefix,omitempty"`  // [REQUIRED] Prefix to Queue Name
}

type Queues struct {
	Activation *Queue `json:"activation,omitempty"` // Message Queue Configuration: Activation
	Mail       *Queue `json:"mail,omitempty"`       // Message Queue Configuration: Email
}

// UTCTimeStamp Return UTC Time Stamp String in RFC 3339
func UTCTimeStamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Helpers
func ToQueue(source interface{}) (*Queue, error) {
	// Do we have Queue Configuration?
	if source == nil { // NO
		return nil, errors.New("No Queue Configuration")
	}

	// Create
	q := &Queue{}

	// Convert Source to JSON
	bo, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	// JSON Back to Queue
	err = json.Unmarshal(bo, q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

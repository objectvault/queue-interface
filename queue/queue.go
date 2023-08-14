// cSpell:ignore vhost
package queue

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
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/objectvault/queue-interface/shared"
)

type AMQPServerConnection struct {
	connection *amqp.Connection          // Server Connection
	channels   *map[string]*amqp.Channel // Channels to Server
	servers    []shared.AMQPConnection   // Connection Settings for Multiple Servers
	prefix     string                    // Queue Name Prefix
	queue      string                    // Default Queue Name
}

func (c *AMQPServerConnection) queueName(name string) (string, error) {
	if name == "" {
		name = c.queue
	}

	if name == "" {
		return "", errors.New("[queueName] Missing Queue Name")
	}

	if c.prefix == "" {
		return name, nil
	}
	return c.prefix + "-" + name, nil
}

func (c *AMQPServerConnection) getChannel(name string) *amqp.Channel {
	// Do we have any Open Channels?
	if c.channels != nil { // YES: Is the Required Channel Opened?
		ch, ok := (*c.channels)[name]
		if ok { // YES: Return Reference to Channel
			return ch
		}
	}

	return nil
}

func (c *AMQPServerConnection) queueURI(con *shared.AMQPConnection) (string, error) {
	// Do we have a User Defined?
	user := con.User
	if user == "" { // NO
		log.Println("[queueURI] Server Configuration Missing User [DEFAULT=guest]")
		user = "guest"
	}

	// Do we have a Password Defined?
	password := con.Password
	if password == "" { // NO
		log.Println("[queueURI] Server Configuration Missing Password [DEFAULT=guest]")
		password = "guest"
	}

	// String Builder
	var builder strings.Builder

	// URI BUILT ACCORDING TO https://www.rabbitmq.com/uri-spec.html
	// AUTHENTICATION COMPONENT of URI //

	// Is a User Set?
	var auth string
	if user != "" { // YES
		// Is Password Set?
		if password != "" { // YES
			fmt.Fprintf(&builder, "%s:%s", user, password)
			auth = builder.String()
			builder.Reset()
		} else { // NO: Just User
			auth = user
		}
	}

	// CONNECTION COMPONENT of URI //

	// Do we have a Servers Definition?
	server := con.Server
	if server == nil { // NO
		return "", errors.New("[queueURI] Server Configuration Missing Server Definition")
	}

	// Get Server Host (IP or ADDRESS)
	host := server.Host
	if host == "" {
		return "", errors.New("[queueURI] Server Configuration Contains Invalid Server Definitions")
	}

	// [OPTIONAL] Get Server Port
	port := server.Port

	// Does Server Have Specific Port?
	var connection string
	if port != 0 { // YES: Build Server Address
		fmt.Fprintf(&builder, "%s:%d", host, port)
		connection = builder.String()
		builder.Reset()
	} else { // NO: Just Add the Server
		connection = host
	}

	// [OPTIONAL] Virtual Host
	vhost := con.VHost

	// BUILD URI //
	if auth != "" {
		if vhost != "" {
			fmt.Fprintf(&builder, "amqp://%s@%s/%s", auth, connection, vhost)
		} else {
			fmt.Fprintf(&builder, "amqp://%s@%s", auth, connection)
		}
	} else {
		if vhost != "" {
			fmt.Fprintf(&builder, "amqp://%s/%s", connection, vhost)
		} else {
			fmt.Fprintf(&builder, "amqp://%s", connection)
		}
	}

	// TODO Handle Server Options (Convert to URI Query Options)
	return builder.String(), nil
}

func (c *AMQPServerConnection) openConnection() (*amqp.Connection, error) {
	limit := len(c.servers)
	// Do we have a Connection Set?
	if limit == 0 { // NO: Abort
		return nil, errors.New("[AMQPServerConnection] No Connection Settings")
	}

	for i := 0; i < limit; i++ {
		server := &c.servers[0]
		// Can we Create a URI from the Information?
		uri, err := c.queueURI(server)
		if err != nil { // NO
			continue
		}

		// Can we Create a Connection from the URI?
		newConnection, err := amqp.Dial(uri)
		if err == nil { // NO
			return newConnection, nil
		}
	}

	return nil, errors.New("[openConnection] Unable to Connect to any Servers")
}

func (c *AMQPServerConnection) SetConnection(s []shared.AMQPConnection) error {
	// Do we already have a connection open?
	if c.connection != nil { // YES: Close it
		c.CloseConnection()
	}

	c.servers = s
	return nil
}

func (c *AMQPServerConnection) Prefix() string {
	return c.prefix
}

func (c *AMQPServerConnection) SetPrefix(p string) error {
	// Do we already have a connection open?
	if c.connection != nil { // YES: Close it
		c.CloseConnection()
	}

	c.prefix = p
	return nil
}

func (c *AMQPServerConnection) DefaultQueue() string {
	return c.queue
}

func (c *AMQPServerConnection) SetDefaultQueue(name string) error {
	c.queue = name
	return nil
}

func (c *AMQPServerConnection) HasConnection() bool {
	return c.connection != nil
}

func (c *AMQPServerConnection) OpenConnection() (*amqp.Connection, error) {
	// Do we already have a connection open?
	if c.connection != nil { // YES: Return it
		return c.connection, nil
	}

	// Open a New Connection
	newConnection, err := c.openConnection()
	if err != nil {
		return nil, err
	}

	c.connection = newConnection
	return c.connection, nil
}

func (c *AMQPServerConnection) ResetConnection() (*amqp.Connection, error) {
	// Do we already have a connection open?
	if c.connection != nil { // YES: Close it
		c.CloseConnection()
	}

	return c.OpenConnection()
}

func (c *AMQPServerConnection) CloseConnection() error {
	// Do we have an open connection?
	if c.connection != nil { // YES: Close it
		// Do we have Open Channels
		if c.channels != nil { // YES: Close any Open Channels
			var err error
			for _, ch := range *c.channels {
				err = ch.Close()
				if err != nil {
					log.Println("[CloseConnection] Error Closing Channel")
				}
			}
		}
		// Clear Channels
		c.channels = nil

		// Close the Connection
		err := c.connection.Close()
		if err != nil {
			log.Println("[CloseConnection] Error Closing Connections")
		}
		// Clear Connections
		c.connection = nil
		return err
	}

	return nil
}

func (c *AMQPServerConnection) IsChannelOpen(name string) bool {
	ch := c.getChannel(name)
	if ch != nil {
		return true
	}

	return false
}

func (c *AMQPServerConnection) OpenChannel(name string) (*amqp.Channel, error) {
	// Do we have a Server Connection?
	if c.connection == nil { // NO: Abort
		return nil, errors.New("[OpenChannel] NO Connection Established")
	}

	// Do we have any Open Channels?
	ch := c.getChannel(name)
	if ch != nil { // YES
		return ch, nil
	}

	// Do we have a Channels Cache?
	if c.channels == nil { // NO: Create it
		c.channels = &map[string]*amqp.Channel{}
	}

	// Open a Channel to the Server
	ch, err := c.connection.Channel()
	if err != nil {
		log.Println("[OpenChannel] Failed to Open Channel [" + name + "]")
		return nil, err
	}

	// Cache Channel
	(*c.channels)[name] = ch
	return ch, nil
}

func (c *AMQPServerConnection) OpenQueueChannel(name string, queue string, create bool) (*amqp.Channel, error) {
	// Get Queue Name
	queue, err := c.queueName(queue)
	if err != nil {
		return nil, err
	}

	// Create Channel Name
	chq := name + "." + queue

	// Does the Queue Channel Exist?
	ch := c.getChannel(chq)
	if ch != nil { // YES
		return ch, nil
	}
	// ELSE: No - Create Channel and/or Queue

	// Can we open the Channel?
	ch, err = c.OpenChannel(chq)
	if err != nil { // NO
		log.Println("[OpenQueueChannel] Unable to Open Channel")
		return nil, err
	}

	// Should we Try to Create the Queue?
	if create { // YES
		// Make Sure Queue is Created
		_, err = ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)

		// Was Queue Created?
		if err != nil { // NO: Abort
			log.Println("[OpenQueueChannel] Failed to Open a Channel to Queue [" + queue + "]")
			return nil, err
		}
	}

	// Cache Queue Channel (ALIAS)
	(*c.channels)[chq] = ch
	return ch, nil
}

func (c *AMQPServerConnection) QueuePublishString(channel string, queue string, msg string) error {
	ch, err := c.OpenQueueChannel(channel, queue, false)
	if err != nil {
		return err
	}

	qName, _ := c.queueName(queue)
	err = ch.Publish(
		"",    // exchange : Queue Default Exchange
		qName, // routing key : Queue Name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

	if err != nil {
		log.Println("[QueuePublishString] Failed Publishing Message to Queue [" + queue + "]")
	}

	return err
}

func (c *AMQPServerConnection) DefaultQueuePublishJSON(channel string, msg interface{}) error {
	return c.QueuePublishJSON(channel, "", msg)
}

func (c *AMQPServerConnection) QueuePublishJSON(channel string, queue string, msg interface{}) error {
	ch, err := c.OpenQueueChannel(channel, queue, false)
	if err != nil {
		return err
	}

	// Marshall Message to JSON Object
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	log.Printf("publishing %dB body (%s)", len(body), body)

	qName, _ := c.queueName(queue)
	err = ch.Publish(
		"",    // exchange : Queue Default Exchange
		qName, // routing key : Queue Name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	if err != nil {
		log.Println("[QueuePublishJSON] Failed Publishing Message to Queue [" + queue + "]")
	}

	return err
}

func (c *AMQPServerConnection) DefaultQueueRetrieve(channel string) (*amqp.Delivery, error) {
	return c.QueueRetrieve(channel, "")
}

func (c *AMQPServerConnection) QueueRetrieve(channel string, queue string) (*amqp.Delivery, error) {
	ch, err := c.OpenQueueChannel(channel, queue, false)
	if err != nil {
		return nil, err
	}

	// Get Next Message on Queue
	qName, _ := c.queueName(queue)
	delivery, ok, err := ch.Get(qName, false)

	// Did we receive an error?
	if err != nil { // YES: Abort
		return nil, err
	}

	// Is Queue Empty?
	if !ok { // YES: Exit
		return nil, nil
	}

	// Return Message
	return &delivery, nil
}

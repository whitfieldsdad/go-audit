package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/deptofdefense/safelock"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/afero"
)

type EventSink interface {
	Write(context.Context, Event) error
}

type JSONLFileEventSink struct {
	Path string `json:"path"`
}

func NewJSONLFileEventSink(path string) (*JSONLFileEventSink, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create directory")
	}
	s := &JSONLFileEventSink{Path: path}
	return s, nil
}

func (s *JSONLFileEventSink) GetLockFilePath() string {
	return fmt.Sprintf("%s.lock", s.Path)
}

func (s *JSONLFileEventSink) Write(ctx context.Context, event Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}
	fs := afero.NewOsFs()
	lock := safelock.NewFileLock(0, s.GetLockFilePath(), fs)
	err = lock.Lock()
	if err != nil {
		return errors.Wrap(err, "failed to lock file")
	}
	defer lock.Unlock()

	f, err := fs.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	return nil
}

type DirectoryEventSink struct {
	Path string `json:"path"`
}

func NewDirectoryEventSink(path string) (*DirectoryEventSink, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create directory")
	}
	s := &DirectoryEventSink{Path: path}
	return s, nil
}

func (s *DirectoryEventSink) GetFilePath(event Event) string {
	return fmt.Sprintf("%s/%s.json", s.Path, event.Header.Id)
}

func (s *DirectoryEventSink) Write(ctx context.Context, event Event) error {
	path := s.GetFilePath(event)
	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	return nil
}

type StdoutEventSink struct{}

func (s *StdoutEventSink) Write(ctx context.Context, event Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}
	fmt.Println(string(b))
	return nil
}

type AMQPEventSink struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	QueueName string `json:"queue_name"`
	conn      *amqp.Connection
}

func NewAMQPEventSink(host string, port int, username, password, queueName string) *AMQPEventSink {
	return &AMQPEventSink{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		QueueName: queueName,
	}
}

func (s *AMQPEventSink) GetURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", s.Username, s.Password, s.Host, s.Port)
}

func (s *AMQPEventSink) GetMaskedURI() string {
	return fmt.Sprintf("amqp://%s:***@%s:%d/", s.Username, s.Host, s.Port)
}

func (s *AMQPEventSink) Connect() error {
	var err error
	maskedURI := s.GetMaskedURI()
	log.Infof("Connecting to %s", maskedURI)
	s.conn, err = amqp.Dial(s.GetURI())
	if err != nil {
		return errors.Wrap(err, "failed to dial amqp")
	}
	log.Infof("Connected to %s", maskedURI)
	return nil
}

func (s *AMQPEventSink) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *AMQPEventSink) Write(ctx context.Context, event Event) error {
	log.Debugf("Writing event to %s", s.GetMaskedURI())
	if s.conn == nil {
		err := s.Connect()
		if err != nil {
			return errors.Wrap(err, "failed to connect")
		}
	}
	ch, err := s.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to create channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(s.QueueName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to declare exchange")
	}
	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	log.Debugf("Wrote event to %s", s.GetMaskedURI())
	return nil
}

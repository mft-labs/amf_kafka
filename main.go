package amf_kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/alecthomas/log4go"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"strconv"
	"time"
)

type KafkaConnector struct {
	host string
	port int
	topic string
	partition int
	log      log4go.Logger
	writer *kafka.Writer
	reader *kafka.Reader
	conn *kafka.Conn
	username string
	password  string
}

func (kfk *KafkaConnector) Init(host string, port int, topic string, partition int, log log4go.Logger) {
	kfk.host = host
	kfk.port = port
	kfk.topic = topic
	kfk.partition = partition
	kfk.log = log
	kfk.username = ""
	kfk.password = ""
	kfk.writer = nil
	kfk.reader = nil
	kfk.conn = nil
	//kfk.InitWriter()
	//kfk.InitReader()
}

func (kfk *KafkaConnector)  UpdateCreds(username, password string) error {
	kfk.username = username
	kfk.password = password
	return nil
}

func (kfk *KafkaConnector) InitConnection() (*kafka.Conn, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	//fmt.Printf("Password security:%v\n",kfk.GetSASL())
	d := kafka.Dialer{
		SASLMechanism: kfk.GetSASL(),
	}
	var err error
	//conn, err := d.DialLeader(context.Background(), "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	kfk.conn, err = d.DialLeader(ctx, "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	if err!=nil {
		return nil, err
	}
	//kfk.conn.SetWriteDeadline(time.Now().Add(60*time.Second))
	fmt.Printf("Established conection to "+kfk.host+":"+strconv.Itoa(kfk.port))
	return kfk.conn,err
}

func (kfk *KafkaConnector) InitConnectionWithSSL(clientcert, clientkey string) (*kafka.Conn, error) {
	keypair, err := tls.LoadX509KeyPair(clientcert, clientkey)
	if err != nil {
		return nil, fmt.Errorf("failt to load Access Key and/or Access Certificate: %s", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	d := kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{keypair},
			//RootCAs:      caCertPool,
		},
		SASLMechanism: kfk.GetSASL(),
	}
	//conn, err := d.DialLeader(context.Background(), "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	kfk.conn, err = d.DialLeader(ctx, "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	if err!=nil {
		fmt.Printf("Error occurred while connecting to KAFKA:%v\n",err)
		return nil, err
	}
	//kfk.conn.SetWriteDeadline(time.Now().Add(60*time.Second))
	fmt.Printf("Established conection to "+kfk.host+":"+strconv.Itoa(kfk.port))
	return kfk.conn,err
}

func (kfk *KafkaConnector) InitConnectionWithScram512(clientcert, clientkey string) (*kafka.Conn, error) {
	keypair, err := tls.LoadX509KeyPair(clientcert, clientkey)
	if err != nil {
		return nil, fmt.Errorf("faild to load Access Key and/or Access Certificate: %s", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	d := kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{keypair},
			//RootCAs:      caCertPool,
		},
		SASLMechanism: kfk.GetSASLSHA512(),
	}
	//conn, err := d.DialLeader(context.Background(), "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	kfk.conn, err = d.DialLeader(ctx, "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	if err!=nil {
		fmt.Printf("Error occurred while connecting with KAFKA SCRAM_SHA_512:%v",err)
		return nil, err
	}
	//kfk.conn.SetWriteDeadline(time.Now().Add(60*time.Second))
	fmt.Printf("Established conection to "+kfk.host+":"+strconv.Itoa(kfk.port))
	return kfk.conn,err
}

func (kfk *KafkaConnector) InitConnectionWithScram256(clientcert, clientkey string) (*kafka.Conn, error) {
	keypair, err := tls.LoadX509KeyPair(clientcert, clientkey)
	if err != nil {
		return nil, fmt.Errorf("failt to load Access Key and/or Access Certificate: %s", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	d := kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{keypair},
			//RootCAs:      caCertPool,
		},
		SASLMechanism: kfk.GetSASLSHA256(),
	}
	//conn, err := d.DialLeader(context.Background(), "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	kfk.conn, err = d.DialLeader(ctx, "tcp", kfk.host+":"+strconv.Itoa(kfk.port), kfk.topic, kfk.partition)
	if err!=nil {
		return nil, err
	}
	//kfk.conn.SetWriteDeadline(time.Now().Add(60*time.Second))
	fmt.Printf("Established conection to "+kfk.host+":"+strconv.Itoa(kfk.port))
	return kfk.conn,err
}

func(kfk *KafkaConnector) GetSASL() sasl.Mechanism {
	return plain.Mechanism{
		Username: kfk.username,
		Password: kfk.password,
	}
}

func (kfk*KafkaConnector) GetSASLSHA512() sasl.Mechanism {
	mech, _ := scram.Mechanism(scram.SHA512, kfk.username, kfk.password)
	return mech
}

func (kfk*KafkaConnector) GetSASLSHA256() sasl.Mechanism {
	mech, _ := scram.Mechanism(scram.SHA256, kfk.username, kfk.password)
	return mech
}

func (kfk *KafkaConnector) InitWriter() error {

	kfk.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kfk.host+":"+strconv.Itoa(kfk.port)},
		Topic:   kfk.topic,
		//Partition: kfk.partition,
		Balancer: &kafka.LeastBytes{},
	})

	return nil
}

func (kfk *KafkaConnector) InitReader() error {
	kfk.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kfk.host+":"+strconv.Itoa(kfk.port)},
		Topic:   kfk.topic,
		//Partition: kfk.partition,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		CommitInterval: time.Second,
	})
	return nil
}

func (kfk *KafkaConnector) GetMessage(msg []byte) kafka.Message {
	msg2 := &kafka.Message{Value:msg}
	return *msg2
}

func (kfk *KafkaConnector) ReadMessage() (kafka.Message, error) {
	m, err := kfk.reader.ReadMessage(context.Background())
	return m, err
}

func (kfk *KafkaConnector) WriteMessage2(msg kafka.Message) error {
	if kfk.conn != nil {
		_, err:= kfk.conn.WriteMessages(msg)
		return err
	}
	return fmt.Errorf("Connection to KAFKA not established")

}

func (kfk *KafkaConnector) WriteMessage(msg kafka.Message) error {
	if kfk.writer != nil {
		fmt.Printf("Using writer\n")
		return kfk.writer.WriteMessages(context.Background(),msg)
	}
	return fmt.Errorf("Invalid connection to kafka")
}

func (kfk *KafkaConnector) Close() {
	kfk.writer.Close()
	kfk.reader.Close()
}

func (kfk *KafkaConnector) CloseWriter() {
	kfk.writer.Close()
}

func (kfk *KafkaConnector) CloseReader() {
	kfk.reader.Close()
}

package ednaevents
//
//import (
//"fmt"
//"os"
//
//"github.com/confluentinc/confluent-kafka-go/kafka"
//)
//
//func main() {
//	fmt.Println("start")
//
//	broker := "BROKER"
//	topic := "demo1"
//	_ = topic
//	//topic := "dp.vault.auditlog"
//	username := "USERNAME"
//	password := "PASSWORD"
//
//	//
//
//	config := &kafka.ConfigMap{
//		"bootstrap.servers": broker,
//		"sasl.username": username,
//		"sasl.password": password,
//		"sasl.mechanism": "PLAIN",
//		"security.protocol": "SASL_SSL",
//	}
//
//	p, err := kafka.NewProducer(config)
//	defer p.Close()
//
//	if err != nil {
//		fmt.Printf("Failed to create producer: %s\n", err)
//		os.Exit(1)
//	}
//
//	fmt.Printf("Created Producer %v\n", p)
//
//	doneChan := make(chan struct{})
//
//	go func() {
//		defer close(doneChan)
//		for e := range p.Events() {
//			switch ev := e.(type) {
//			case *kafka.Message:
//				m := ev
//				if m.TopicPartition.Error != nil {
//					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
//				} else {
//					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
//						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
//				}
//				return
//
//			default:
//				fmt.Printf("Ignored event: %s\n", ev)
//			}
//		}
//	}()
//
//	value := "Hello Go!"
//	p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
//
//	<-doneChan
//
//	fmt.Println("DONE!")
//
//}
//

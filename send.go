package main

import (
	"log"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Message struct{
	Text string `json:"text"`
	Sender string `json:"sender"`
}

func failOnErrore(err error,msg string) {
	if err!=nil{
		log.Fatalf("%s:%s",msg,err)
	}
}

//var body Message

func bodyFrom(args []string) string {
    var s string
    if (len(args) < 2) || os.Args[1] == "" {
        s = "hello"
    } else {
        s = strings.Join(args[1:], " ")
    }
    return s
}

func GetMessage(c *gin.Context){
//	body.Text = c.PostForm("text")
//	body.Sender = c.PostForm("sender")
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnErrore(err,"Failed to connect to RabbitMQ")
	defer conn.Close()

	ch,err := conn.Channel()
	failOnErrore(err,"Failed to open a Channel")
	queue,err := ch.QueueDeclare(
		"hello",//name
		false,//durable
		false,//delete when unsend
		false,//exclusive
		false,//no-wait
		nil,//arguments
	)
	failOnErrore(err,"Failed to create a queue")
//	ss:= []byte(fmt.Sprintf(`{"text":"%s","sender":"%s"}`,body.Text,body.Sender))
	ss:= bodyFrom(os.Args)

	err = ch.Publish(
		"",//exchange
		queue.Name,//routingkey
		false,//mandatory
		false,//immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:"text/plain",
			Body:  []byte(ss),

		})
		failOnErrore(err,"Failed to Publish a message")
		log.Printf("[x] send %s",ss)

}
func main() {
	r := gin.Default()
	r.POST("/message",GetMessage)
		r.Run()
}

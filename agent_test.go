package emailagent

import (
	"io/ioutil"
	"log"
	"testing"
)

// TestSendEmail TODO
// 2019/09/23 22:29:58
func TestSendEmail(t *testing.T) {
	agent, err := NewEmailAgent("smtp.qq.com", 465, "sunliang711@qq.com", "ydsgtgyzupoofcba")
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	agent.SendEmail([]string{"vimisbug001@163.com", "vimisbug002@163.com"}, "subject is not empty", "this is from a email agent ", true)
	// agent.SendEmail("some@163.com", "subject is not empty2", "this is from a email agent 2")
}

func TestSendNew(t *testing.T) {
	agent, err := NewEmailAgent("smtp.qq.com", 465, "sunliang711@qq.com", "ydsgtgyzupoof")
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()

	msg := NewMessage("subject 11", "<b>hihi body</b>hh")
	msg.To = []string{"vimisbug001@163.com", "vimisbug002@163.com"}
	msg.CC = []string{"sunliang711@163.com"}
	pic, _ := ioutil.ReadFile("gateio.jpeg")
	msg.Attach("gateio.jpeg", pic)
	txt,_ := ioutil.ReadFile("haha.txt")
	msg.Attach("haha.txt",txt)

	err = agent.Send(msg, true)
	if err != nil {
		t.Fatal(err)
	}
}

package emailagent

import (
	"log"
	"testing"
)

// TestSendEmail TODO
// 2019/09/23 22:29:58
func TestSendEmail(t *testing.T) {
	agent, err := NewEmailAgent("smtp.qq.com", 465, "some@qq.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	agent.SendEmail("some@163.com", "subject is not empty", "this is from a email agent ")
	agent.SendEmail("some@163.com", "subject is not empty2", "this is from a email agent 2")
}

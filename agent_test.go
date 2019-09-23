package emailagent

import (
	"log"
	"testing"
)

// TestSendEmail TODO
// 2019/09/23 22:29:58
func TestSendEmail(t *testing.T) {
	agent := NewEmailAgent("smtp.qq.com", 465, "sunliang711@qq.com", "password")
	if err := agent.Init(); err != nil {
		log.Fatal(err)
	}
	agent.SendEmail("someone@163.com", "subject is not empty", "this is from a email agent ")
	agent.Close()
}

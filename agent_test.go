package emailagent

import (
	"log"
	"testing"
)

// TestSendEmail TODO
// 2019/09/23 22:29:58
func TestSendEmail(t *testing.T) {
	agent, err := NewEmailAgent("smtp.qq.com", 465, "sunliang711@qq.com", "ydsgtgyzupoof")
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	agent.SendEmail([]string{"vimisbug001@163.com","vimisbug002@163.com"}, "subject is not empty", "this is from a email agent ",true)
	// agent.SendEmail("some@163.com", "subject is not empty2", "this is from a email agent 2")
}

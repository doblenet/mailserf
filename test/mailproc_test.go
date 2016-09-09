package test_mailserf

import (
	mailserf "../"
	"fmt"
	"testing"
	"strings"
)

const (
	
	test1_to = "somebody#somedomain.tld@engine.doblenet.internal"
	test1_from = "someone@bratos.city" 
	test1 = `From: One <somebody@somedomain.tld>
To: One Girl <aria.something@temple.bratos.city>
Subject: today's duties
Importance: High
X-Loop-Prevention: 1234abcdef
Received: from mta2-69.azuresend.com (mta2-69.azuresend.com [213.206.105.69])
    by hermes.doblenet.net (Postfix) with ESMTP id 3sG4QG5QwWz3X
    for <jltallon@doblenet.com>; Fri, 19 Aug 2016 16:01:41 +0200 (CEST)
Date: Fri, 19 Aug 2016 15:41:38 +0200
Received: from web044 (127.0.0.1) by mta2-68.azuresend.com id hms8ka1iod83 for <jltallon@doblenet.com>; Fri, 19 Aug 2016 15:41:38 +0200 (envelope-from <bounce-jltallon=doblenet.com@eu.azuresend.com>)
X-Original-To: jltallon@doblenet.com
Delivered-To: jltallon@doblenet.com
X-Greylist: delayed 1201 seconds by postgrey-1.35 at hermes.doblenet.net; Fri, 19 Aug 2016 16:01:42 CEST
`
)

func TestProcessor(t *testing.T) {

	e := mailserf.Envelope{From: test1_from, To: test1_to}
	fmt.Println("envelope:",e)
	
	
	p := mailserf.New(e)
// 	if nil != err {
// 		t.Fail()
// 	}
	
	p.Procs["subject"] = func(x string, c mailserf.ProcessorCtx) bool {
		fmt.Println("$Subject:", x)
		return true
	}
	
	
	var c mailserf.ProcessorCtx
	err := p.Run(strings.NewReader(test1), &c)
	if nil!=err {
		fmt.Println("ERROR: ", err.Error())
		t.Fail()
	}
	
	fmt.Println("\n** Received:")
	//fmt.Println(p.ReceivedHeaders())
	dumpList(p.ReceivedHeaders())
	
	fmt.Println("\n** Unmatched:")
// 	fmt.Println(p.UnmatchedHeaders())
	dumpMap(p.UnmatchedHeaders())
	
}


func dumpMap(m map[string]string) {
	for k,v := range m {
		fmt.Printf("%s:\t%s\n", k,v)
	}
}
func dumpList(l []string) {
	for i,v := range l {
		fmt.Printf("[%d] %s\n", i,v)
	}
}

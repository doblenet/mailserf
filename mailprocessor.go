// Copyright 2016 Jose Luis Tallon
// License: APL-2.0
package mailserf

// e-mail (RFC2821-compliant) processor using IOC for message parsing
// 

import (
	"errors"
	"bufio"
	"io"
	"strings"
	"fmt"
)

type ProcessorCtx interface {
	
}

type HeaderProc func(x string, ctx ProcessorCtx) bool

type MailProcessor struct {
	
	envelope	*Envelope
	Procs		map[string]HeaderProc
	headers		map[string]string
	received	[]string
	do_body		bool
	bodylines	[]string
}


func New(e Envelope) *MailProcessor {
	
	o := MailProcessor{envelope: &e}
	
	o.Procs = make(map[string]HeaderProc, 3)
	o.headers = make(map[string]string, 5)
	o.received = make([]string, 0, 3)

	// Initialize HeaderProcs with a suitable function; 
	// o gets bound to the newly-created object, so the reference will work
// 	o.Procs["received"] = func(x string, ctx ProcessorCtx) bool {
// 		o.received = append(o.received, x)
// 		return true
// 	}
	
	return &o
}


func (o *MailProcessor) SetHook(header_name string, f HeaderProc) bool {
	
	if _,ok := o.Procs[header_name]; ok {
		// Already registered? better don't overwrite ...
		return false
	}
	o.Procs[header_name] = f
	return true
}


func (o *MailProcessor) UnmatchedHeaders() map[string]string {
	
	return o.headers
}

func (o *MailProcessor) ReceivedHeaders() []string {
	return o.received
}

func (o *MailProcessor) EnableBodyProcessing(y bool) {
	o.do_body = y
}


func (o *MailProcessor) Run(input io.Reader, ctx ProcessorCtx) error {

	if len(o.Procs) < 1 {
		return errors.New("No processing functions supplied!?")
	}
	
	
	var prevLine string
	var prevHeader string
	
	// Instantiate the scanner
	// the default "splitFunc", ScanLines() already removes the trailing '\r', so we're good to go
	scanner := bufio.NewScanner(input)
	
	for scanner.Scan() {
		
		l := scanner.Text()		// overwritten on each loop...
		if ""==l && ""!=prevLine {
			fmt.Println("#### EMPTY ####")
			// end of headers; Exit loop, since we already have what we wanted :O
			break
		} else {
			prevLine = l		// store current line (raw)
		}
		
		// Continuation?
		if ' '==l[0] || '\t'==l[0] {	// continuation....
			if "received" == prevHeader {
				o.received = append(o.received,l)	// full line
			} else {
				v := strings.TrimSpace(l)
				o.headers[prevHeader] = (o.headers[prevHeader]+v)
			}
			// Do NOT update the previous header (continuation...)
			continue
		}
		
		// Regular header: do process it
		i := strings.Index(l,":")
		if i > 2 {
			
			h := strings.ToLower(l[:i])
			r := l[i+2:]
			if f,ok := o.Procs[h]; ok {
				if ! f(r, ctx) {
					o.headers[h] = r
				}
			} else if "received"==h {
				o.received  = append(o.received,l)	// append the full line !
			} else {
				o.headers[h] = r
			}
			prevHeader = h
		}
		
	}	// end scanning
	
	if err:=scanner.Err(); nil!=err {
		fmt.Println("ERROR")
		return err
	} else {
		// Non-errored... so some lines left
		
		for scanner.Scan() {
			l := scanner.Text()
			if o.do_body {
				o.bodylines = append(o.bodylines,l)
			}
			// else just swallow it up :O
		}
	}
	
	return nil
}

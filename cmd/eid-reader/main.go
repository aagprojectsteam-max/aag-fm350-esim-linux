package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	euiccmbim "github.com/damonto/euicc-go/driver/mbim"
	"github.com/damonto/euicc-go/lpa"
	wwanmbim "github.com/damonto/wwan-go/mbim"
)

const (
	device   = "/dev/wwan0mbim0"
	esimSlot = 2 // euicc-go uses 1-based slot numbering
)

func masked(value string) string {
	if len(value) <= 12 { return value }
	return value[:6] + strings.Repeat("*", len(value)-12) + value[len(value)-6:]
}

func run() (err error) {
	if len(os.Args) != 2 { return fmt.Errorf("usage: %s OUTPUT_FILE", os.Args[0]) }
	outputFile := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second); defer cancel()
	rawClient, err := wwanmbim.Open(ctx, wwanmbim.WithProxy(device), wwanmbim.WithSlot(esimSlot))
	if err != nil { return fmt.Errorf("opening MBIM proxy for eSIM slot 2: %w", err) }
	channel, err := euiccmbim.NewWithClient(rawClient)
	if err != nil { _ = rawClient.Close(); return fmt.Errorf("creating eUICC MBIM channel: %w", err) }
	client, err := lpa.New(&lpa.Options{Channel: channel, Timeout: 45*time.Second})
	if err != nil { _ = channel.Disconnect(); return fmt.Errorf("opening the eUICC ISD-R application: %w", err) }
	defer func(){ if closeErr:=client.Close(); err==nil && closeErr!=nil { err=fmt.Errorf("closing eUICC channel: %w", closeErr) } }()
	eid, err := client.EID(); if err != nil { return fmt.Errorf("reading EID: %w", err) }
	eidText := fmt.Sprintf("%X", eid); if eidText=="" { return fmt.Errorf("the eUICC returned an empty EID") }
	if err:=os.WriteFile(outputFile, []byte(eidText+"\n"), 0o600); err!=nil { return fmt.Errorf("saving EID: %w", err) }
	fmt.Printf("EID_MASKED=%s\n", masked(eidText)); fmt.Printf("EID_SAVED=%s\n", outputFile)
	profiles, listErr := client.ListProfile(nil,nil)
	if listErr != nil { fmt.Printf("PROFILE_COUNT=unavailable\n"); fmt.Printf("PROFILE_LIST_WARNING=%v\n", listErr) } else { fmt.Printf("PROFILE_COUNT=%d\n", len(profiles)) }
	return nil
}
func main(){ if err:=run(); err!=nil { fmt.Fprintf(os.Stderr,"ERROR: %v\n",err); os.Exit(1) } }

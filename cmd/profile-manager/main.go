package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	euiccmbim "github.com/damonto/euicc-go/driver/mbim"
	"github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	wwanmbim "github.com/damonto/wwan-go/mbim"
)

const (
	device = "/dev/wwan0mbim0"
	esimSlot = 2
)

func mask(v string) string {
	v = strings.TrimSpace(v)
	if len(v) <= 10 { return strings.Repeat("*", len(v)) }
	return v[:6] + strings.Repeat("*", len(v)-10) + v[len(v)-4:]
}

func openClient(ctx context.Context) (*lpa.Client, error) {
	raw, err := wwanmbim.Open(ctx, wwanmbim.WithProxy(device), wwanmbim.WithSlot(esimSlot))
	if err != nil { return nil, err }
	ch, err := euiccmbim.NewWithClient(raw)
	if err != nil { _ = raw.Close(); return nil, err }
	c, err := lpa.New(&lpa.Options{Channel: ch, Logger: slog.New(slog.NewTextHandler(io.Discard,nil)), Timeout: 90*time.Second})
	if err != nil { _ = ch.Disconnect(); return nil, err }
	return c,nil
}

func list(c *lpa.Client) ([]*sgp22.ProfileInfo,error) {
	for i:=0;i<5;i++ {
		p,err:=c.ListProfile(nil,nil)
		if err==nil { return p,nil }
		if i==4 { return nil,err }
		time.Sleep(2*time.Second)
	}
	return nil,errors.New("unreachable")
}

func show(ps []*sgp22.ProfileInfo) {
	fmt.Printf("PROFILE_COUNT=%d\n",len(ps))
	for i,p:=range ps {
		fmt.Printf("PROFILE_%d_NAME=%s\n",i+1,strings.TrimSpace(p.ProfileName))
		fmt.Printf("PROFILE_%d_PROVIDER=%s\n",i+1,strings.TrimSpace(p.ServiceProviderName))
		fmt.Printf("PROFILE_%d_ICCID_MASKED=%s\n",i+1,mask(p.ICCID.String()))
		fmt.Printf("PROFILE_%d_CLASS=%s\n",i+1,p.ProfileClass.String())
		fmt.Printf("PROFILE_%d_STATE=%s\n",i+1,p.ProfileState.String())
	}
}

func choose(ps []*sgp22.ProfileInfo,suffix,provider,name string)(*sgp22.ProfileInfo,error) {
	var hit []*sgp22.ProfileInfo
	for _,p:=range ps {
		if suffix!=""&&!strings.HasSuffix(p.ICCID.String(),suffix){continue}
		if provider!=""&&!strings.EqualFold(strings.TrimSpace(p.ServiceProviderName),provider){continue}
		if name!=""&&!strings.EqualFold(strings.TrimSpace(p.ProfileName),name){continue}
		if suffix==""&&provider==""&&name==""{continue}
		hit=append(hit,p)
	}
	if len(hit)!=1{return nil,fmt.Errorf("selector matched %d profiles",len(hit))}
	return hit[0],nil
}

func main() {
	action:=flag.String("action","list","list, enable, disable")
	suffix:=flag.String("suffix","","ICCID suffix selector")
	provider:=flag.String("provider","","provider selector")
	name:=flag.String("name","","profile name selector")
	flag.Parse()

	ctx,cancel:=context.WithTimeout(context.Background(),4*time.Minute)
	defer cancel()
	c,err:=openClient(ctx)
	if err!=nil{fmt.Fprintln(os.Stderr,"ERROR:",err);os.Exit(1)}
	defer c.Close()

	ps,err:=list(c)
	if err!=nil{fmt.Fprintln(os.Stderr,"ERROR:",err);os.Exit(1)}
	show(ps)
	if *action=="list"{return}

	p,err:=choose(ps,*suffix,*provider,*name)
	if err!=nil{fmt.Fprintln(os.Stderr,"ERROR:",err);os.Exit(1)}

	var opErr error
	switch *action {
	case "enable":
		if p.ProfileState.String()=="enabled"{fmt.Println("TARGET_ALREADY_ENABLED=yes");return}
		fmt.Println("PROFILE_ENABLE_REQUEST=starting")
		opErr=c.EnableProfile(p.ICCID,true)
	case "disable":
		if p.ProfileState.String()=="disabled"{fmt.Println("TARGET_ALREADY_DISABLED=yes");return}
		fmt.Println("PROFILE_DISABLE_REQUEST=starting")
		opErr=c.DisableProfile(p.ICCID,true)
	default:
		fmt.Fprintln(os.Stderr,"ERROR: unknown action");os.Exit(2)
	}

	time.Sleep(3*time.Second)
	after,readErr:=list(c)
	if readErr==nil {
		show(after)
		want:="enabled"; if *action=="disable"{want="disabled"}
		for _,x:=range after {
			if x.ICCID.String()==p.ICCID.String()&&x.ProfileState.String()==want {
				if opErr!=nil{fmt.Printf("PROFILE_%s_TRANSPORT_WARNING=%v\n",strings.ToUpper(*action),opErr)}
				fmt.Printf("PROFILE_%s_VERIFIED=yes\n",strings.ToUpper(*action))
				return
			}
		}
	}
	if opErr!=nil{fmt.Fprintln(os.Stderr,"ERROR:",opErr)}
	if readErr!=nil{fmt.Fprintln(os.Stderr,"READBACK_ERROR:",readErr)}
	os.Exit(1)
}

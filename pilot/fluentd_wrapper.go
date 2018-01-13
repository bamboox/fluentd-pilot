package pilot

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"strings"
	"context"
	log "github.com/Sirupsen/logrus"
)

var fluentd *exec.Cmd
var cancel context.CancelFunc
var ctx context.Context

func StartFluentd(debug bool) error {
	if fluentd != nil {
		return fmt.Errorf("fluentd already started")
	}
	log.Warn("start fluentd")
	cmdArgs := []string{"-c", "/etc/fluentd/fluentd.conf", "-p", "/etc/fluentd/plugins"}
	if debug {
		cmdArgs = append(cmdArgs, "-v")
	}
	ctx, cancel = context.WithCancel(context.Background())
	fluentd = exec.CommandContext(ctx, "/usr/bin/fluentd", cmdArgs...)
	fluentd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	fluentd.Stderr = os.Stderr
	fluentd.Stdout = os.Stdout
	err := fluentd.Start()
	if err != nil {
		go func() {
			fluentd.Wait()
		}()
	}
	return err
}

func shell(command string) string {
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error %v", err)
	}
	return strings.Replace(string(out), "\n", "", -1)
}

func ReloadFluentd() error {
	if fluentd == nil {
		return fmt.Errorf("fluentd have not started")
	}
	log.Warn("reload fluentd")
	ch := make(chan struct{})
	go func(pid int) {
		command := fmt.Sprintf("pgrep -P %d", pid)
		childId := shell(command)
		if (childId == "") {
			//restart: always
			os.Exit(1)
			close(ch)
			return
		}

		log.Infof("before reload childId : %s", childId)
		fluentd.Process.Signal(syscall.SIGHUP)
		time.Sleep(5 * time.Second)
		afterChildId := shell(command)
		log.Infof("after reload childId : %s", afterChildId)
		if childId == afterChildId {
			log.Infof("kill childId : %s", childId)
			shell("kill -9 " + childId)
		}
		close(ch)
	}(fluentd.Process.Pid)
	<-ch
	return nil
}

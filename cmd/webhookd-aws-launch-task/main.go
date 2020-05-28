package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aaronland/go-aws-ecs"
	"github.com/sfomuseum/go-flags/flagset"	
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	fs := flagset.NewFlagSet("webhookd-aws-launch-task")
	
	var ecs_dsn = fs.String("ecs-dsn", "", "A valid (go-whosonfirst-aws) ECS DSN.")

	var ecs_container = fs.String("ecs-container", "", "The name of your AWS ECS container.")
	var ecs_cluster = fs.String("ecs-cluster", "", "The name of your AWS ECS cluster.")
	var ecs_task = fs.String("ecs-task", "", "The name of your AWS ECS task (inclusive of its version number),")

	var ecs_launch_type = fs.String("ecs-launch-type", "FARGATE", "...")
	var ecs_public_ip = fs.String("ecs-public-ip", "ENABLED", "...")

	var ecs_subnets multi.MultiString
	fs.Var(&ecs_subnets, "ecs-subnet", "One or more AWS subnets in which your task will run.")

	var ecs_security_groups multi.MultiString
	fs.Var(&ecs_security_groups, "ecs-security-group", "One of more AWS security groups your task will assume.")

	var monitor = fs.Bool("monitor", false, "...")

	var logs = fs.Bool("logs", false, "...")
	var logs_dsn = fs.String("logs-dsn", "", "A valid (go-whosonfirst-aws) CloudWatchLogs DSN.")

	var mode = fs.String("mode", "cli", "...")
	var command = fs.String("command", "", "...")
	var command_insecure = fs.Bool("command-insecure", false, "...")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "WEBHOOKD", true)

	if err != nil {
		log.Fatal(err)
	}

	if *command == "" {
		log.Fatal("Missing command")
	}

	if *logs == true {
		*monitor = true
	}

	if *logs_dsn == "" {
		*logs_dsn = *ecs_dsn
	}

	if *mode == "lambda" {

		expand := func(candidates []string, sep string) []string {

			expanded := make([]string, 0)

			for _, c := range candidates {

				for _, v := range strings.Split(c, sep) {
					expanded = append(expanded, v)
				}
			}

			return expanded
		}

		ecs_subnets = expand(ecs_subnets, ",")
		ecs_security_groups = expand(ecs_security_groups, ",")
	}

	task_opts := &ecs.TaskOptions{
		DSN:            *ecs_dsn,
		Task:           *ecs_task,
		Container:      *ecs_container,
		Cluster:        *ecs_cluster,
		Subnets:        ecs_subnets,
		SecurityGroups: ecs_security_groups,
		LaunchType:     *ecs_launch_type,
		PublicIP:       *ecs_public_ip,
	}

	launchTask := func(command string, args ...interface{}) (interface{}, error) {

		str_cmd := fmt.Sprintf(command, args...)
		cmd := strings.Split(str_cmd, ",")	// match the weird Docker/ECS stuff

		// log.Println(strings.Join(cmd, " "))
		// log.Println(cmd)
		
		task_rsp, err := ecs.LaunchTask(task_opts, cmd...)

		if err != nil {
			return nil, err
		}

		log.Println("TASKS", task_rsp.Tasks)

		if !*monitor {
			return task_rsp.Tasks, nil
		}

		monitor_opts := &ecs.MonitorTaskOptions{
			DSN:       *ecs_dsn,
			Container: *ecs_container,
			Cluster:   *ecs_cluster,
			WithLogs:  *logs,
			LogsDSN:   *logs_dsn,
		}

		return ecs.MonitorTasks(monitor_opts, task_rsp.Tasks...)
	}

	switch *mode {

	case "cli":

		for _, repo := range flag.Args() {

			rsp, err := launchTask(*command, repo)

			if err != nil {
				log.Fatal(err)
			}

			log.Println(rsp)
		}

	case "lambda":

		re, err := regexp.Compile(`^[a-zA-Z0-9\-_]+$`)

		if err != nil {
			log.Fatal(err)
		}

		string_handler := func(ctx context.Context, payload string) (interface{}, error) {

			if !*command_insecure {

				if !re.MatchString(payload) {
					return nil, errors.New("Invalid payload")
				}
			}

			return launchTask(*command, payload)
		}

		lambda.Start(string_handler)

	default:
		log.Fatal("Unknown mode")
	}

	os.Exit(0)
}

package main

import (
	"context"
	_ "encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/whosonfirst/go-whosonfirst-aws/ecs"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"os"
	"strings"
)

func main() {

	var ecs_dsn = flag.String("ecs-dsn", "", "A valid (go-whosonfirst-aws) ECS DSN.")

	var ecs_container = flag.String("ecs-container", "", "The name of your AWS ECS container.")
	var ecs_cluster = flag.String("ecs-cluster", "", "The name of your AWS ECS cluster.")
	var ecs_task = flag.String("ecs-task", "", "The name of your AWS ECS task (inclusive of its version number),")

	var ecs_launch_type = flag.String("ecs-launch-type", "FARGATE", "...")
	var ecs_public_ip = flag.String("ecs-public-ip", "ENABLED", "...")

	var ecs_subnets flags.MultiString
	flag.Var(&ecs_subnets, "ecs-subnet", "One or more AWS subnets in which your task will run.")

	var ecs_security_groups flags.MultiString
	flag.Var(&ecs_security_groups, "ecs-security-group", "One of more AWS security groups your task will assume.")

	var monitor = flag.Bool("monitor", false, "...")

	var logs = flag.Bool("logs", false, "...")
	var logs_dsn = flag.String("logs-dsn", "", "A valid (go-whosonfirst-aws) CloudWatchLogs DSN.")

	var mode = flag.String("mode", "cli", "...")
	var command = flag.String("command", "", "...")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("WEBHOOKD")

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
		cmd := strings.Split(str_cmd, " ")

		log.Println("LAUNCH TASK", cmd)

		task_rsp, err := ecs.LaunchTask(task_opts, cmd...)

		if err != nil {
			log.Println("OH NO", err)
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

		string_handler := func(ctx context.Context, payload string) (interface{}, error) {
			
			log.Println("RECEIVED", payload)
			
			return launchTask(*command, payload)
		}

		lambda.Start(string_handler)

	default:
		log.Fatal("Unknown mode")
	}

	os.Exit(0)
}

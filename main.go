package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/mackerelio/checkers"
)

type Option struct {
	target   string
	strict   bool
	gt       float64
	lt       float64
	warning  bool
	critical bool
	name     string
}

func main() {
	opt := new(Option)
	flag.StringVar(&opt.target, "target", "", "target metric key")
	flag.BoolVar(&opt.strict, "strict", false, "create UNKNOWN alert when target metric is not found")
	flag.Float64Var(&opt.gt, "gt", math.NaN(), "set condition: [target value] > [option value]")
	flag.Float64Var(&opt.lt, "lt", math.NaN(), "set condition: [target value] < [option value]")
	flag.BoolVar(&opt.warning, "warning", false, "create CRITICAL alert when target value meets condition")
	flag.BoolVar(&opt.critical, "critical", true, "create WARNING alert when target value meets condition")
	flag.StringVar(&opt.name, "name", "check-value-from-metrics-plugin", "checker name for report")
	flag.Parse()

	if opt.target == "" {
		log.Fatalln("-target is required")
	}

	threshold := math.NaN()
	gt := false
	if !math.IsNaN(opt.gt) {
		threshold = opt.gt
		gt = true
	}
	if !math.IsNaN(opt.lt) {
		if gt {
			log.Fatalln("specify only one of either -gt or -lt")
		}
		threshold = opt.lt
	}
	if math.IsNaN(threshold) {
		log.Fatalln("-gt or -lt is required")
	}

	if opt.warning && opt.critical {
		log.Fatalln("specify only one of either -warning or -critical")
	}

	runner := &Runner{
		target:    opt.target,
		strict:    opt.strict,
		gt:        gt,
		threshold: threshold,
		warning:   opt.warning,
	}
	status, message := runner.Run()
	checker := checkers.NewChecker(status, message)
	checker.Name = opt.name
	checker.Exit()
}

type Runner struct {
	target    string
	strict    bool
	gt        bool
	threshold float64
	warning   bool
}

func (r *Runner) Run() (status checkers.Status, message string) {
	status = checkers.OK

	value, meetsCond, err := r.meetsCond()
	if err != nil {
		if !r.strict && errors.Is(err, ErrTargetValueNotFound) {
			return
		}
		status = checkers.UNKNOWN
		message = err.Error()
		return
	}

	if meetsCond {
		if r.warning {
			status = checkers.WARNING
		} else {
			status = checkers.CRITICAL
		}
		op := "<"
		if r.gt {
			op = ">"
		}
		message = fmt.Sprintf("%f %s %f", value, op, r.threshold)
		return
	}
	return
}

var ErrTargetValueNotFound = fmt.Errorf("target value not found")

func (r *Runner) meetsCond() (value float64, meetsCond bool, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	var valueStr string
	for scanner.Scan() {
		items := strings.Fields(scanner.Text())
		if len(items) < 3 {
			continue
		}
		key := items[0]
		if r.target != key {
			continue
		}
		valueStr = items[1]
		break
	}

	if valueStr == "" {
		err = ErrTargetValueNotFound
		return
	}

	value, err = strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return
	}
	if r.gt && value > r.threshold {
		meetsCond = true
		return
	}
	if !r.gt && value < r.threshold {
		meetsCond = true
		return
	}
	return
}

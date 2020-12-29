package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func viperGetIP(name string) string {
	value := viper.GetString(name)
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	ip := net.ParseIP(value)
	if ip == nil {
		fmt.Fprintf(os.Stderr, "Invalid IP address: %v\n", value)
		os.Exit(1)
	}
	return ip.String()
}

func viperGetRegExp(name string) *regexp.Regexp {
	value := viper.GetString(name)
	re, err := regexp.Compile(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid regexp: %v: %v\n", value, err)
		os.Exit(1)
	}
	return re
}

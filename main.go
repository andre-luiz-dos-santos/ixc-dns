package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

func boot() error {
	file := flag.String("config", "", "configuration file")
	ixcDSN := flag.String("ixc-dsn", "/ixcprovedor", "IXC database DSN")
	dnsBind := flag.String("dns-bind", "127.0.0.1:53", "DNS bind address")
	flag.Parse()

	viper.SetConfigType("toml")
	if *file != "" {
		viper.SetConfigFile(*file)
	} else {
		viper.SetConfigName("ixc-dns")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/ixc-dns")
	}

	viper.SetDefault("ixc.dsn", ixcDSN)
	viper.SetDefault("dns.bind", dnsBind)
	viper.SetDefault("dns.ttl", 1)
	viper.SetDefault("dns.login.regexp", `^login-([^.]+)\.`)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	ixc := &IXC{}

	err = ixc.OpenDB(viper.GetString("ixc.dsn"))
	if err != nil {
		return err
	}

	srv := &dns.Server{
		Net:  "udp",
		Addr: viper.GetString("dns.bind"),
	}
	srv.Handler = &dnsHandler{
		ixc:        ixc,
		loginRE:    viperGetRegExp("dns.login.regexp"),
		answerTTL:  viper.GetInt("dns.ttl"),
		filteredIP: viperGetIP("dns.ip.filtered"),
		notFoundIP: viperGetIP("dns.ip.not_found"),
		errorIP:    viperGetIP("dns.ip.error"),
	}
	return srv.ListenAndServe()
}

func main() {
	err := boot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

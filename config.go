package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Config struct {
	Files struct {
		Pubring  string
		Mlist2   string
		Pubkey   string
		Secring  string
		Adminkey string
		Help     string
		Pooldir  string
		Maildir  string
		IDlog    string
		ChunkDB  string
	}
	Urls struct {
		Fetch   bool
		Pubring string
		Mlist2  string
	}
	Mail struct {
		Sendmail       bool
		Pipe           string
		Outfile        bool
		UseTLS         bool
		SMTPRelay      string
		SMTPPort       int
		MXRelay        bool
		EnvelopeSender string
		Username       string
		Password       string
		OutboundName   string
		OutboundAddy   string
		CustomFrom     bool
	}
	Stats struct {
		Minlat    int
		Maxlat    int
		Minrel    float32
		Relfinal  float32
		Chain     string
		Numcopies int
		Distance  int
		StaleHrs  int
	}
	Pool struct {
		Size int
		Rate int
		Loop int
	}
	Remailer struct {
		Name        string
		Address     string
		Exit        bool
		MaxSize     int
		IDexp       int
		ChunkExpire int
		Keylife     int
		Keygrace    int
		Loglevel    string
		Daemon      bool
	}
}

func init() {
	var err error
	// Function as a client
	flag.BoolVar(&flag_client, "mail", false, "Function as a client")
	flag.BoolVar(&flag_client, "m", false, "Function as a client")
	// Send (from pool)
	flag.BoolVar(&flag_send, "send", false, "Force pool send")
	flag.BoolVar(&flag_send, "S", false, "Force pool send")
	// Perform remailer actions
	flag.BoolVar(&flag_remailer, "remailer", false,
		"Perform routine remailer actions")
	flag.BoolVar(&flag_remailer, "M", false,
		"Perform routine remailer actions")
	// Start remailer as a daemon
	flag.BoolVar(&flag_daemon, "daemon", false,
		"Start remailer as a daemon. (Requires -M")
	flag.BoolVar(&flag_daemon, "D", false,
		"Start remailer as a daemon. (Requires -M")
	// Remailer chain
	flag.StringVar(&flag_chain, "chain", "", "Remailer chain")
	flag.StringVar(&flag_chain, "l", "", "Remailer chain")
	// Recipient address
	flag.StringVar(&flag_to, "to", "", "Recipient email address")
	flag.StringVar(&flag_to, "t", "", "Recipient email address")
	// Subject header
	flag.StringVar(&flag_subject, "subject", "", "Subject header")
	flag.StringVar(&flag_subject, "s", "", "Subject header")
	// Number of copies
	flag.IntVar(&flag_copies, "copies", 0, "Number of copies")
	flag.IntVar(&flag_copies, "c", 0, "Number of copies")
	// Config file
	flag.StringVar(&flag_config, "config", "", "Config file")
	// Read STDIN
	flag.BoolVar(&flag_stdin, "read-mail", false, "Read a message from stdin")
	flag.BoolVar(&flag_stdin, "R", false, "Read a message from stdin")
	// Write to STDOUT
	flag.BoolVar(&flag_stdout, "stdout", false, "Write message to stdout")
	// Inject dummy
	flag.BoolVar(&flag_dummy, "dummy", false, "Inject a dummy message")
	flag.BoolVar(&flag_dummy, "d", false, "Inject a dummy message")
	// Disable dummy messaging
	flag.BoolVar(&flag_nodummy, "nodummy", false, "Don't send dummies")
	// Print Version
	flag.BoolVar(&flag_version, "version", false, "Print version string")
	flag.BoolVar(&flag_version, "V", false, "Print version string")
	// Memory usage
	flag.BoolVar(&flag_meminfo, "meminfo", false, "Print memory info")

	// Figure out the dir of the yamn binary
	var dir string
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	flag.StringVar(&flag_basedir, "dir", dir, "Base directory")

	// Set defaults and read config file
	cfg.Files.Pubkey = path.Join(flag_basedir, "key.txt")
	cfg.Files.Pubring = path.Join(flag_basedir, "pubring.mix")
	cfg.Files.Secring = path.Join(flag_basedir, "secring.mix")
	cfg.Files.Mlist2 = path.Join(flag_basedir, "mlist2.txt")
	cfg.Files.Adminkey = path.Join(flag_basedir, "adminkey.txt")
	cfg.Files.Help = path.Join(flag_basedir, "help.txt")
	cfg.Files.Pooldir = path.Join(flag_basedir, "pool")
	cfg.Files.Maildir = path.Join(flag_basedir, "Maildir")
	cfg.Files.IDlog = path.Join(flag_basedir, "idlog")
	cfg.Files.ChunkDB = path.Join(flag_basedir, "chunkdb")
	cfg.Urls.Fetch = true
	cfg.Urls.Pubring = "http://www.mixmin.net/yamn/pubring.mix"
	cfg.Urls.Mlist2 = "http://www.mixmin.net/yamn/mlist2.txt"
	cfg.Mail.Sendmail = false
	cfg.Mail.Outfile = false
	cfg.Mail.SMTPRelay = "snorky.mixmin.net"
	cfg.Mail.SMTPPort = 587
	cfg.Mail.UseTLS = true
	cfg.Mail.MXRelay = true
	cfg.Mail.EnvelopeSender = "nobody@nowhere.invalid"
	cfg.Mail.Username = ""
	cfg.Mail.Password = ""
	cfg.Mail.OutboundName = "Anonymous Remailer"
	cfg.Mail.OutboundAddy = "remailer@domain.invalid"
	cfg.Mail.CustomFrom = false
	cfg.Stats.Minrel = 98.0
	cfg.Stats.Relfinal = 99.0
	cfg.Stats.Minlat = 2
	cfg.Stats.Maxlat = 60
	cfg.Stats.Chain = "yamn4,*,*"
	cfg.Stats.Numcopies = 1
	cfg.Stats.Distance = 2
	cfg.Stats.StaleHrs = 24
	cfg.Pool.Size = 45
	cfg.Pool.Rate = 65
	cfg.Pool.Loop = 300
	cfg.Remailer.Name = "anon"
	cfg.Remailer.Address = "mix@nowhere.invalid"
	cfg.Remailer.Exit = false
	cfg.Remailer.MaxSize = 12
	cfg.Remailer.IDexp = 14
	cfg.Remailer.ChunkExpire = 60
	cfg.Remailer.Keylife = 60
	cfg.Remailer.Keygrace = 28
	cfg.Remailer.Loglevel = "info"
	cfg.Remailer.Daemon = false
}

func flags() {
	var err error
	flag.Parse()
	flag_args = flag.Args()
	if flag_version {
		fmt.Println(version)
		os.Exit(0)
	} else if flag_config != "" {
		err = gcfg.ReadFileInto(&cfg, flag_config)
		if err != nil {
			fmt.Fprintf(
				os.Stderr, "Unable to read %s\n", flag_config)
			os.Exit(1)
		}
	} else if os.Getenv("YAMNCFG") != "" {
		err = gcfg.ReadFileInto(&cfg, os.Getenv("YAMNCFG"))
		if err != nil {
			fmt.Fprintf(
				os.Stderr, "Unable to read %s\n", flag_config)
			os.Exit(1)
		}
	} else {
		fn := path.Join(flag_basedir, "yamn.cfg")
		err = gcfg.ReadFileInto(&cfg, fn)
		if err != nil {
			if !flag_client {
				fmt.Println(err)
			}
			fmt.Println("Using internal, default config.")
		}
	}
}

var flag_basedir string
var flag_client bool
var flag_send bool
var flag_remailer bool
var flag_daemon bool
var flag_chain string
var flag_to string
var flag_subject string
var flag_args []string
var flag_config string
var flag_copies int
var flag_stdin bool
var flag_stdout bool
var flag_dummy bool
var flag_nodummy bool
var flag_version bool
var flag_meminfo bool
var cfg Config

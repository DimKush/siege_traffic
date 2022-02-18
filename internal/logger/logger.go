package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/spf13/viper"
)

const (
	ERROR = iota
	WARNING
	INFO
	DEBUG
	TRACE
)

type log_struct struct {
	path_to_logs   string
	today_log_file string
	zLog           *zerolog.Logger
}

func (data *log_struct) createLogNewDate() error {
	var strb strings.Builder

	log_path_dir_foo := func() string {
		if runtime.GOOS == "windows" {
			return "\\"
		}

		return "/"
	}
	var separator string = log_path_dir_foo()

	strb.WriteString(data.path_to_logs)
	strb.WriteString(separator)
	strb.WriteString(data.today_log_file)

	log_file, err := os.Create(strb.String())

	if err != nil {
		err := fmt.Errorf("Cannot create the file %s, Reason: %s", strb.String(), err.Error())
		return err
	}

	log_file.Close()

	return nil
}

func (data *log_struct) checkLogDateFile() (bool, error) {
	files, err := ioutil.ReadDir(data.path_to_logs)
	if err != nil {
		err := fmt.Errorf("Cannot open the directory with logs %s. Reason : %s", data.path_to_logs, err.Error())
		return false, err
	}

	for _, file := range files {
		if file.Name() == data.today_log_file {
			return true, nil
		}
	}

	return false, nil
}

func (data *log_struct) init() error {
	strLevel := strings.ToUpper(viper.GetString("log_level"))

	var complete_path = func(data string) string {
		var strb strings.Builder
		if runtime.GOOS == "windows" {
			strb.WriteString(data)
			strb.WriteString("\\")
			return strb.String()
		} else {
			strb.WriteString(data)
			strb.WriteString("/")
			return strb.String()
		}
	}

	if viper.GetString("log_path") == "" {
		log_path_dir_foo := func() string {
			if runtime.GOOS == "windows" {
				return "c:\\dimkush_guestbook\\log\\"
			}

			return "/opt/dimkush_guestbook/log/"
		}
		data.path_to_logs = log_path_dir_foo()
	} else {
		data.path_to_logs = complete_path(viper.GetString("log_path"))
	}

	current_dt := time.Now()

	var strb strings.Builder
	strb.WriteString("guestbook_")

	current_date_str := current_dt.Format("2006-Jan-02")

	strb.WriteString(current_date_str)

	strb.WriteString(".log")

	data.today_log_file = strb.String()

	if needNewfile, err := data.checkLogDateFile(); err != nil {
		return err
	} else {
		if !needNewfile {
			if err := data.createLogNewDate(); err != nil {
				return err
			}
		}
	}

	strb.Reset()

	strb.WriteString(data.path_to_logs)
	strb.WriteString(data.today_log_file)

	file, err := os.OpenFile(strb.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		err := fmt.Errorf("Cannot open the log file %s, Reason : %s", strb.String(), err.Error())
		return err
	}

	writer := diode.NewWriter(file, 10000, 10*time.Microsecond, func(missed int) {
		fmt.Printf("Logger dropped %d messages", missed)
	})

	zlogger := zerolog.New(writer).With().Caller().Timestamp().Logger().Output(file)
	data.zLog = &zlogger

	switch strLevel {
	case "ERROR":
		{
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}
	case "WARNING":
		{
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		}
	case "INFO":
		{
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	case "DEBUG":
		{
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	case "TRACE":
		{
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
	default:
		{
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}
	}

	return nil
}

func NewLogger() (zerolog.Logger, error) {
	var logger_obj log_struct

	if err := logger_obj.init(); err != nil {
		return zerolog.Logger{}, err
	}

	return *logger_obj.zLog, nil
}

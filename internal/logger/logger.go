package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	ERROR = iota
	WARNING
	INFO
	DEBUG
	TRACE
)

func createLogNewDate(pathToLogs, todayLogFile string) error {
	var strb strings.Builder

	log_path_dir_foo := func() string {
		if runtime.GOOS == "windows" {
			return "\\"
		}

		return "/"
	}
	var separator string = log_path_dir_foo()

	strb.WriteString(pathToLogs)
	strb.WriteString(separator)
	strb.WriteString(todayLogFile)

	logFile, err := os.Create(strb.String())

	if err != nil {
		err := fmt.Errorf("Cannot create the file %s, Reason: %s", strb.String(), err.Error())
		return err
	}

	fmt.Printf("file create : %s\n", logFile.Name())
	err = os.Chmod(strb.String(), 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = logFile.Close()
	if err != nil {
		fmt.Println("Cannot close log file : %s", err)
		os.Exit(1)
	}
	return nil
}

func checkLogDateFile(pathToLogs, todayLogFile string) (bool, error) {
	files, err := ioutil.ReadDir(pathToLogs)
	if err != nil {
		err := fmt.Errorf("Cannot open the directory with logs %s. Reason : %s", pathToLogs, err.Error())
		return false, err
	}

	for _, file := range files {
		if file.Name() == todayLogFile {
			return true, nil
		}
	}

	return false, nil
}

func InitLogger() error {
	strLevel := strings.ToUpper(viper.GetString("log_level"))
	pathToLogs := ""
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
		logPathDirFoo := func() string {
			if runtime.GOOS == "windows" {
				return "c:\\traffic\\"
			}

			return "/opt/traffic/"
		}
		pathToLogs = logPathDirFoo()
	} else {
		pathToLogs = complete_path(viper.GetString("log_path"))
	}

	if _, err := os.Stat(pathToLogs); os.IsNotExist(err) {
		errMkDir := os.Mkdir(pathToLogs, 0777)

		err = multierr.Append(err, errMkDir)
	}

	currentDt := time.Now()

	var strb strings.Builder
	strb.WriteString("traffic_")

	currentDateStr := currentDt.Format("2006-Jan-02")

	strb.WriteString(currentDateStr)

	strb.WriteString(".log")

	todayLogFile := strb.String()

	if needNewFile, err := checkLogDateFile(pathToLogs, todayLogFile); err != nil {
		return err
	} else {
		if !needNewFile {
			if err := createLogNewDate(pathToLogs, todayLogFile); err != nil {
				return err
			}
		}
	}

	strb.Reset()

	strb.WriteString(pathToLogs)
	strb.WriteString(todayLogFile)

	fmt.Printf("file open : %s\n", strb.String())
	file, err := os.OpenFile(strb.String(), syscall.O_RDWR, 0666)
	if err != nil {
		err := fmt.Errorf("Cannot open the log file %s, Reason : %s", strb.String(), err.Error())
		return err
	}

	writer := diode.NewWriter(file, 10000, 10*time.Microsecond, func(missed int) {
		fmt.Printf("Logger dropped %d messages", missed)
	})

	log.Logger = zerolog.New(writer).With().Caller().Timestamp().Logger().Output(file)

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
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}

	return nil
}

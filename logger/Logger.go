package logger

import (
	logger "github.com/sirupsen/logrus"
)

const LOG_LEVEL_INFO = 1
const LOG_LEVEL_WARN = 2
const LOG_LEVEL_ERROR = 3

func Log(trace string, log string, logLevel int) {

	//if configs.CanDebug() {
		switch logLevel {
		case LOG_LEVEL_INFO:
			{
				logger.Info("["+trace+"]", " "+log)
			}

		case LOG_LEVEL_WARN:
			{
				logger.Warn("["+trace+"]", " "+log)
			}

		case LOG_LEVEL_ERROR:
			{
				logger.Error("["+trace+"]", " "+log)
			}
		}
	//}
}

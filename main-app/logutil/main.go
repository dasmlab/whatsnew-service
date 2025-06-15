// Package logutil provides a reusable structured logger using Logrus.
package logutil

import (
        "os"

        "github.com/sirupsen/logrus"
)

// InitLogger initializes a logrus.Logger with standard formatting and returns it.
// The provided component name is included in all log entries.
func InitLogger(component string) *logrus.Entry {
        logger := logrus.New()

        // Output to stdout instead of the default stderr
        logger.Out = os.Stdout

        // Set the log level (customize as needed)
        logger.SetLevel(logrus.DebugLevel)

        // Use text formatter with timestamp
        //logger.SetFormatter(&logrus.TextFormatter{
        //      FullTimestamp: true,
        //})

        // Set Log Formatter (JSON for CKS-style/Loki style things)
        logger.SetFormatter(&logrus.JSONFormatter{})

        // Return an Entry pre-tagged with component field
        return logger.WithField("component", component)
}



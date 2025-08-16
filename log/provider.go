package log

import (
	"strings"

	"github.com/sirupsen/logrus"
)

const maskedFilter = "MASKED"

type ProviderLogger struct {
	filters  []string
	replacer *strings.Replacer
	logger   *logrus.Logger
}

//nolint:gomnd,mnd
func NewProviderLogger(baseLogger *logrus.Logger, filters ...string) *ProviderLogger {
	// replacer uses pairs for replacement
	// we have all the first word of pairs in filters array.
	// so in replace words. we must have: [filters[0], masked, filters[1], masked, ...]
	// indexes 0, 2, 4, 6 are actual words and indexes 1,3,5,7 are replacements.
	// in each iteration we fill 2 indexes of replaceWords
	// replaceWords[even] = actualWord
	// replaceWords[odd] = replacement
	replaceWords := make([]string, 2*len(filters))
	for i, j := 0, 0; i < len(filters); i++ {
		replaceWords[j] = filters[i]
		replaceWords[j+1] = maskedFilter

		j += 2
	}

	r := strings.NewReplacer(replaceWords...)

	return &ProviderLogger{
		replacer: r,
		filters:  filters,
		logger:   baseLogger,
	}
}

func (p *ProviderLogger) LogHTTPCall(req, resp string) {
	filteredRequest, filteredResponse := p.filterRequestResponse(req, resp)
	fields := logrus.Fields{
		"request":  filteredRequest,
		"response": filteredResponse,
	}

	logrus.WithFields(fields).Infof("call captured")
}

func (p *ProviderLogger) filterRequestResponse(req, resp string) (string, string) {
	return p.replacer.Replace(req), p.replacer.Replace(resp)
}

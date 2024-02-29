package regex

import "regexp"

func RegexCompile(regexs []string) (regexCls []*regexp.Regexp, err error) {
	regexCls = make([]*regexp.Regexp, 0, len(regexs))
	if len(regexs) > 0 {
		for _, reg := range regexs {
			if regcomplie, err := regexp.Compile(reg); err != nil {
				return nil, err
			} else {
				regexCls = append(regexCls, regcomplie)
			}
		}
	}
	return regexCls, nil
}

func RegexMatchString(regexCls []*regexp.Regexp, value string) bool {
	for _, reg := range regexCls {
		if reg.MatchString(value) {
			return true
		}
	}
	return false
}

func RegexNonMatchString(regexCls []*regexp.Regexp, value string) bool {
	return !RegexMatchString(regexCls, value)
}

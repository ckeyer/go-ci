package checks

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gojp/goreportcard/check"
)

type Grade string

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
	Error         string              `json:"error"`
}

type checksResp struct {
	Checks               []score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	LastRefresh          time.Time `json:"last_refresh"`
	HumanizedLastRefresh string    `json:"humanized_last_refresh"`
}

func RunChecks(dir string, filenames []string) (checksResp, error) {
	if len(filenames) == 0 {
		return checksResp{}, fmt.Errorf("no .go files found")
	}

	checks := []check.Check{
		check.GoFmt{Dir: dir, Filenames: filenames},
		// check.GoVet{Dir: dir, Filenames: filenames},
		check.GoLint{Dir: dir, Filenames: filenames},
		// check.GoCyclo{Dir: dir, Filenames: filenames},
		// check.License{Dir: dir, Filenames: []string{}},
		// check.Misspell{Dir: dir, Filenames: filenames},
		// check.IneffAssign{Dir: dir, Filenames: filenames},
		// check.ErrCheck{Dir: dir, Filenames: filenames}, // disable errcheck for now, too slow and not finalized
	}

	ch := make(chan score)
	for _, c := range checks {
		go func(c check.Check) {
			p, summaries, err := c.Percentage()
			errMsg := ""
			if err != nil {
				log.Errorf("ERROR: (%s) %v", c.Name(), err)
				errMsg = err.Error()
			}
			s := score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Weight:        c.Weight(),
				Percentage:    p,
				Error:         errMsg,
			}
			ch <- s
		}(c)
	}

	resp := checksResp{
		Files:                len(filenames),
		LastRefresh:          time.Now().UTC(),
		HumanizedLastRefresh: time.Now().String(),
	}

	var total float64
	var totalWeight float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		log.WithFields(log.Fields{
			"desc":       s.Description,
			"name":       s.Name,
			"error":      s.Error,
			"weight":     s.Weight,
			"percentage": s.Percentage,
		}).Infof("%+v", s.FileSummaries)
		resp.Checks = append(resp.Checks, s)
		total += s.Percentage * s.Weight
		totalWeight += s.Weight
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}
	total /= totalWeight

	// sort.Sort(ByWeight(resp.Checks))
	resp.Average = total
	resp.Issues = len(issues)

	return resp, nil
}

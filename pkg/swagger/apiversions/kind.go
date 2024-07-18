package apiversions

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ParseKindResults() {
	previousTable := &ResourcesTable{
		Version: "",
		Kinds:   map[string][]string{},
		Headers: nil,
		Rows:    nil,
	}

	for _, version := range []string{
		"1.18.19",
		"1.19.11",
		"1.20.7",
		"1.21.2",
		"1.22.4",
		"1.23.0",
	} {
		// TODO spin up cluster, collect data, kill cluster
		//   instead of reading from static file

		headers, rows, err := ReadCSV(fmt.Sprintf("../kube/data/v%s-api-resources.txt", version))
		utils.Die(err)
		rsTable, err := NewResourcesTable(version, headers, rows)
		utils.Die(err)

		fmt.Printf("%s\n", rsTable.KindResourcesTable())

		//for kind, apiVersions := range rsTable.Kinds {
		//	fmt.Printf("%s, %+v\n", kind, apiVersions)
		//}

		fmt.Printf("comparing %s to %s\n", previousTable.Version, rsTable.Version)
		resourceDiff := previousTable.Diff(rsTable)
		fmt.Printf("changed:\n%s\n", resourceDiff.Table(map[string]bool{}, map[string]bool{}))
		//for kind, change := range resourceDiff.Changed {
		//	if len(change.Added) != 0 || len(change.Removed) != 0 {
		//		fmt.Printf("kind %s; added: %+v, removed: %+v, same: %+v\n", kind, change.Added, change.Removed, change.Same)
		//	}
		//}

		previousTable = rsTable
	}
}

func ReadCSV(path string) ([]string, [][]string, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to read file %s", path)
	}
	lines := strings.Split(string(in), "\n")
	headings, startIndexes := FindFieldStarts(lines[0])
	logrus.Debugf("headings: %+v\n%+v\n", headings, startIndexes)
	var rows [][]string
	for ix, line := range lines[1:] {
		logrus.Debugf("line: <%s>\n", line)
		if ix == len(lines)-2 && line == "" {
			break
		}
		var fields []string
		for i, start := range startIndexes {
			var stop int
			if i == len(startIndexes)-1 {
				stop = len(line)
			} else {
				stop = startIndexes[i+1]
			}
			trimmed := strings.TrimRight(line[start:stop], " ")
			fields = append(fields, trimmed)
			logrus.Debugf("trimmed? %+v\n", trimmed)
		}
		rows = append(rows, fields)
	}
	return headings, rows, nil
}

func FindFieldStarts(line string) ([]string, []int) {
	regex := regexp.MustCompile(`\S+`)
	nums := regex.FindAllStringIndex(line, -1)
	var headings []string
	var startIndexes []int
	for _, ns := range nums {
		logrus.Debugf("%+v\n", ns)
		start, stop := ns[0], ns[1]
		headings = append(headings, line[start:stop])
		startIndexes = append(startIndexes, start)
	}
	return headings, startIndexes
}

package core

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func FindInSlice[K comparable](arr []K, comparer func(K) bool) (K, bool) {
	var result K
	found := false
	for _, cur := range arr {
		if comparer(cur) {
			found = true
			result = cur
			break
		}
	}
	return result, found
}

func FindAllInSlice[K comparable](arr []K, comparer func(K) bool) ([]K, bool) {
	var results []K
	for _, cur := range arr {
		if comparer(cur) {
			results = append(results, cur)
		}
	}
	return results, len(results) > 0
}

func removeFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func changedDirs(filesChanged []string, modulesDir string, modulesPath string) []string {
	dirschanged := make([]string, 0)
	for _, fc := range filesChanged {
		if strings.HasPrefix(fc, modulesDir) {
			a := strings.Split(fc, "/")
			if len(a) > 2 { // make sure the changed file is of the form [azure resource-group main.tf]
				dirExists := true
				_, err := os.Stat(path.Join(modulesPath, a[1]))
				if os.IsNotExist(err) {
					dirExists = false
				}
				inDirschanged := false
				for _, dir := range dirschanged {
					if dir == a[1] {
						inDirschanged = true
					}
				}
				if inDirschanged == false && dirExists == true {
					dirschanged = append(dirschanged, a[1])
				}
			}
		}
	}

	return dirschanged
}

func getVersion(dir string) (string, error) {
	file, err := os.Open(path.Join(dir, "VERSION"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	retval := strings.TrimSuffix(string(bytes), "\n")

	return retval, nil
}

func FormatCurrency(value float64, currency Currency) string {
	return FormatCurrencyWithPrecision(value, currency, 2)
}

func FormatCurrencyWithPrecision(value float64, currency Currency, decimals int) string {
	symbol := currencySymbols[currency]
	precisionFormat := fmt.Sprintf("%%.%df", decimals)
	formattedValue := fmt.Sprintf(precisionFormat, value)
	if symbol == "" && currency != "" {
		return fmt.Sprintf("%s %s", formattedValue, currency)
	} else {
		return fmt.Sprintf("%s%s", symbol, formattedValue)
	}
}

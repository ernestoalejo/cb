package v0

import (
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("ngmin", 0, ngmin)
}

func ngmin(c config.Config, q *registry.Queue) error {
	scripts := filepath.Join("client", "temp", "scripts")
	if err := filepath.Walk(scripts, walkFn); err != nil {
		return errors.New(err)
	}
	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.New(err)
	}
	if info.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".js" {
		return nil
	}

	lines, err := utils.ReadLines(path)
	if err != nil {
		return err
	}

	newlines := []string{}
	for _, line := range lines {
		// TODO:
		//GlobalCtrl.$inject = ['$rootScope', '$location', 'Selector'];
		//function GlobalCtrl($rootScope, $location, Selector) {
		//m.factory('GlobalMsg', function($timeout) {
		//m.directive('match', function() {
		//m.config(function($routeProvider, $locationProvider) {
		newlines = append(newlines, line)
	}

	return nil
}

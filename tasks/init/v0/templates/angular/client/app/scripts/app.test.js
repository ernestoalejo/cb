'use strict';


var m = angular.module('test', ['{{.AppName}}', 'ngMockE2E']);

m.run(function($httpBackend) {
  $httpBackend.whenGET(/.*/).passThrough();
});


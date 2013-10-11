'use strict';


var m = angular.module('mocks.modals', []);


// This mock only serves testing purposes, to allow importing httpInterceptor
// without causing a cycle in the injector.
m.factory('Modal', function() {
  return {};
});

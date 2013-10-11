'use strict';


var m = angular.module('mocks.scroll', []);


m.factory('ScrollTo', function() {
  var dests_ = [];

  var f = function(dest) {
    dests_.push(dest);
  };

  f.formError = function() {
    dests_.push('formError');
  };

  f.getDests = function() {
    return dests_;
  };

  return f;
});

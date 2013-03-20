'use strict';


var m = angular.module('httpInterceptor', ['services.global']);

m.config(function($httpProvider) {
  $httpProvider.responseInterceptors.push('httpInterceptor');
});

m.factory('httpInterceptor', function($q, GlobalMsg) {
  var total = 0;

  return function(promise) {
    total++;
    GlobalMsg.set('loading');

    return promise.then(function(response) {
      total--;
      if (total == 0 && GlobalMsg.get() == 'loading')
        GlobalMsg.set('');
      return response;
    }, function(response) {
      total--;
      GlobalMsg.set('');
      $('#http-error').modal();
      return $q.reject(response);
    });
  }
});

'use strict';


var m = angular.module('httpInterceptor', ['services.global']);

m.config(function($httpProvider) {
  $httpProvider.interceptors.push('httpInterceptor');
});


m.factory('httpInterceptor', function($q, GlobalMsg) {
  var total_ = 0;

  function error_() {
    total_--;
    GlobalMsg.set('');
    ErrorRegister.set('http-error');
  }

  return {
    'request': function(config) {
      total_++;
      GlobalMsg.set('loading');
      return config || $q.when(config);
    },

    'requestError': function(rejection) {
      error_();
      return $q.reject(rejection);
    },

    'response': function(response) {
      total_--;
      if (total_ == 0 && GlobalMsg.get() == 'loading') {
        GlobalMsg.set('');
      }
      return response || $q.when(response);
    },

    'responseError': function(rejection) {
      error_();
      return $q.reject(rejection);
    }
  };
});

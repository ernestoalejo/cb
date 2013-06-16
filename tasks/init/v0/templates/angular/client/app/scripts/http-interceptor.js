'use strict';


var m = angular.module('httpInterceptor', ['services.global']);

m.config(function($httpProvider) {
  $httpProvider.interceptors.push('httpInterceptor');
});


m.factory('httpInterceptor', function($q, GlobalMsg, ErrorRegister) {
  var total_ = 0;

  function error_() {
    total_--;
    GlobalMsg.hideLoading();
    ErrorRegister.set('http-error');
  }

  return {
    'request': function(config) {
      total_++;
      GlobalMsg.showLoading();
      return config || $q.when(config);
    },

    'requestError': function(rejection) {
      error_();
      return $q.reject(rejection);
    },

    'response': function(response) {
      total_--;
      if (total_ <= 0 && GlobalMsg.isLoading()) {
        GlobalMsg.hideLoading();
      }
      return response || $q.when(response);
    },

    'responseError': function(rejection) {
      error_();
      return $q.reject(rejection);
    }
  };
});

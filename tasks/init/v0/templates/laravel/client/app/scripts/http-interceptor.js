'use strict';


var m = angular.module('httpInterceptor', []);

m.config(function($httpProvider) {
  $httpProvider.interceptors.push('httpInterceptor');
});


m.factory('httpInterceptor', function($q, GlobalMsg, ErrorRegister) {
  function error_() {
    GlobalMsg.forceHideLoading();
    ErrorRegister.set('http-error');
  }

  return {
    'request': function(config) {
      GlobalMsg.showLoading();
      return config || $q.when(config);
    },

    'requestError': function(rejection) {
      error_();
      return $q.reject(rejection);
    },

    'response': function(response) {
      GlobalMsg.hideLoading();

      if (response.config.url.substring(0, 3) != '/_/') {
        return response || $q.when(response);
      }

      if (response.headers('X-Response-Processor') != 'json') {
        error_();
        return $q.reject('no response processor header');
      }
      return response.data || $q.when(response.data);
    },

    'responseError': function(rejection) {
      error_();
      return $q.reject(rejection);
    }
  };
});

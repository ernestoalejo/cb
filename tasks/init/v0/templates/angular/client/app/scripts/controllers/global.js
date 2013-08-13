'use strict';


var m = angular.module('controllers.global', [
]);


m.controller('AppCtrl', function($rootScope, $location) {
  $rootScope.$on('$routeChangeSuccess', function(e) {
    if (window._gaq)
      window._gaq.push(['_trackPageview', $location.url()]);
  });

  $rootScope.$on('$routeChangeError', function(e, cur, prev, msg) {
    if (msg == 'notlogged') {
      $location.path('/');
    } else if (msg == 'logged') {
      $location.path('/accounts/login');
    } else if (msg == 'admin') {
      $location.path('/');
    } else {
      throw new Error('unkwnown route error: ' + msg);
    }
  });
});


m.controller('NotFoundCtrl', function($http, $location) {
  $http.post('/_/not-found', {path: $location.path()});
});


m.controller('ErrorCtrl', function($scope, ErrorRegister) {
  $scope.ErrorRegister = ErrorRegister;
  $scope.close = function() {
    ErrorRegister.clean();
  };
});

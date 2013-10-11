'use strict';


var m = angular.module('controllers.global', []);


m.controller('AppCtrl', function($scope, $rootScope, $location, Modal,
    ErrorRegister, GlobalMsg) {
  $rootScope.$on('$routeChangeSuccess', function() {
    if (window._gaq) {
      window._gaq.push(['_trackPageview', $location.url()]);
    }
  });

  $scope.ErrorRegister = ErrorRegister;
  $scope.$watch('ErrorRegister.isNull()', function() {
    if (!ErrorRegister.isNull()) {
      ErrorRegister.clean();
      Modal.errorDialog();
    }
  });

  $scope.GlobalMsg = GlobalMsg;
});


m.controller('NotFoundCtrl', function($http, $location) {
  $http.post('/_/not-found', {path: $location.path()});
});

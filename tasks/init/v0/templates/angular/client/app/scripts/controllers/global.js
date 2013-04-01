'use strict';


var m = angular.module('controllers.global', [
  'services.global'
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


m.controller('NotFoundCtrl', function() {
  // empty
});


m.controller('GlobalMsgCtrl', function($scope, GlobalMsg) {
  $scope.gm = GlobalMsg;

  $scope.close = function() {
    GlobalMsg.set('');
  };
});


m.controller('FeedbackCtrl', function($scope, $http, GlobalMsg) {
  $scope.open = function() {
    $scope.dlgOpened = true;
  };

  $scope.close = function() {
    $scope.dlgOpened = false;
  };

  $scope.dlgOpened = false;
  $scope.opts = {
    backdropFade: true,
    dialogFade: true
  };

  $scope.send = function() {
    $scope.dlgOpened = false;

    var msg = $scope.message;
    $scope.message = '';

    $http.post('/_/feedback', {message: msg}).success(function() {
      GlobalMsg.setTemp('Hemos recibido tu mensaje correctamente', 'success');
    }).error(function() {
      $scope.message = msg;
    });
  };
});


m.controller('ErrorCtrl', function($scope, ErrorRegister) {
  $scope.ErrorRegister = ErrorRegister;
  $scope.close = function() {
    ErrorRegister.clean();
  };
});

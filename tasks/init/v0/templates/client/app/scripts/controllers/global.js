'use strict';


var m = angular.module('controllers.global', [
  'services.global'
]);


m.controller('GlobalCtrl', function($rootScope, $location, Selector) {
  // Change the sidebar and navbar when navigating
  $rootScope.$on('$routeChangeStart', function() {
    Selector.setDirty();
  });
  $rootScope.$on('$routeChangeSuccess', function(e) {
    Selector.clearDirty();

    // Google Analytics (if present)
    if (window._gaq)
      window._gaq.push(['_trackPageview', $location.url()]);
  });

  $rootScope.$on('$routeChangeError', function(e, cur, prev, msg) {
    if (msg == 'notlogged') {
      $location.path('/');
    } else if (msg == 'logged') {
      PagesCache.add($location.url());
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


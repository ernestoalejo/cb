'use strict';


/**
 * Control some global services needed for the page.
 */
GlobalCtrl.$inject = ['$rootScope', '$location', 'Selector'];
function GlobalCtrl($rootScope, $location, Selector) {
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
}


/**
 * Show a page not-found error for the client routes.
 */
NotFoundCtrl.$inject = [];
function NotFoundCtrl() { }


/**
 * Controller for the global message showed on success/error/warning/...
 */
GlobalMsgCtrl.$inject = ['$scope', 'GlobalMsg'];
function GlobalMsgCtrl($scope, GlobalMsg) {
  $scope.gm = GlobalMsg;

  $scope.close = function() {
    GlobalMsg.set('');
  };
}


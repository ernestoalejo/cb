'use strict';


describe('Controller: GlobalMsgCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, GlobalMsg;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');

    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller('GlobalMsgCtrl', {$scope: scope});
  }));

  it('should save the service in the scope', function() {
    expect(scope.gm).toBe(GlobalMsg);
  });

  it('should clean the message on close', function() {
    GlobalMsg.set('testing');
    scope.close();
    expect(GlobalMsg.get()).toBe('');
  });
});


describe('Controller: GlobalMsgCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, Selector, $rootScope, $location;
  beforeEach(inject(function($injector) {
    Selector = $injector.get('Selector');
    $rootScope = $injector.get('$rootScope');
    $location = $injector.get('$location');

    var $controller = $injector.get('$controller');

    scope = $rootScope.$new();
    $controller('GlobalCtrl', {$scope: scope});
  }));

  it('should dirty & clean the selector on navigation', function() {
    Selector.clearDirty();
    $rootScope.$broadcast('$routeChangeStart');
    expect(Selector.isDirty()).toBeTruthy();
    $rootScope.$broadcast('$routeChangeSuccess');
    expect(Selector.isDirty()).toBeFalsy();
  });

  it('should notify analytics of page changes if present', function() {
    $rootScope.$broadcast('$routeChangeSuccess');
    window._gaq = [];
    $location.url('/testing');
    $rootScope.$broadcast('$routeChangeSuccess');
    expect(window._gaq.length).toBe(1);
    expect(window._gaq[0].length).toBe(2);
    expect(window._gaq[0][0]).toBe('_trackPageview');
    expect(window._gaq[0][1]).toBe('/testing');
  });

  it('should react on errors', (function() {
      $location.path('/notlogged');
      $rootScope.$broadcast('$routeChangeError', '/notlogged', '/prev', 'notlogged');
      expect($location.path()).toBe('/');

      $location.path('/logged');
      $rootScope.$broadcast('$routeChangeError', '/logged', '/prev', 'logged');
      expect($location.path()).toBe('/accounts/login');

      $location.path('/admin');
      $rootScope.$broadcast('$routeChangeError', '/admin', '/prev', 'admin');
      expect($location.path()).toBe('/');

      $location.path('/unknown');
      expect(function() {
        $rootScope.$broadcast('$routeChangeError', '/unknown', '/prev', 'unknown');
      }).toThrow(new Error("unkwnown route error: unknown"));
    }));
});


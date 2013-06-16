'use strict';


describe('Controller: AppCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, $rootScope, $location;
  beforeEach(inject(function($injector) {
    $rootScope = $injector.get('$rootScope');
    $location = $injector.get('$location');

    var $controller = $injector.get('$controller');

    scope = $rootScope.$new();
    $controller('AppCtrl', {$scope: scope});
  }));

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
    $rootScope.$broadcast('$routeChangeError', '/notlogged', '/prev',
        'notlogged');
    expect($location.path()).toBe('/');

    $location.path('/logged');
    $rootScope.$broadcast('$routeChangeError', '/logged', '/prev',
        'logged');
    expect($location.path()).toBe('/accounts/login');

    $location.path('/admin');
    $rootScope.$broadcast('$routeChangeError', '/admin', '/prev', 'admin');
    expect($location.path()).toBe('/');

    $location.path('/unknown');
    expect(function() {
      $rootScope.$broadcast('$routeChangeError', '/unknown', '/prev',
          'unknown');
    }).toThrow(new Error('unkwnown route error: unknown'));
  }));
});


describe('Controller: NotFoundCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, $httpBackend;
  beforeEach(inject(function($injector) {
    $httpBackend = $injector.get('$httpBackend');

    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller('NotFoundCtrl', {$scope: scope});
  }));

  it('should reset the form & show a message on success', function() {
    $httpBackend.expectPOST('/_/not-found').respond({});
    $httpBackend.flush();
  });
});


describe('Controller: ErrorCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, GlobalMsg, ErrorRegister;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');
    ErrorRegister = $injector.get('ErrorRegister');
    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller('ErrorCtrl', {$scope: scope});
  }));

  it('should scope the register', function() {
    expect(scope.ErrorRegister).toBe(ErrorRegister);
  });

  it('should clean the register on close', function() {
    ErrorRegister.set('foo');
    scope.close();
    expect(ErrorRegister.isNull()).toBeTruthy();
  });
});

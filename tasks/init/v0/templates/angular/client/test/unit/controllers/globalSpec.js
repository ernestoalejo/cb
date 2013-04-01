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


describe('Controller: FeedbackCtrl', function() {
  beforeEach(module('controllers.global'));

  var scope, GlobalMsg, $httpBackend;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');
    $httpBackend = $injector.get('$httpBackend');

    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller('FeedbackCtrl', {$scope: scope});
  }));

  it('should open/close the dialog', function() {
    expect(scope.dlgOpened).toBeFalsy();
    scope.open();
    expect(scope.dlgOpened).toBeTruthy();
    scope.close();
    expect(scope.dlgOpened).toBeFalsy();
  });

  it('should have the correct options', function() {
    expect(scope.opts).toEqualData({
      backdropFade: true,
      dialogFade: true
    });
  });

  it('should reset the form & show a message on success', function() {
    $httpBackend.expectPOST('/_/feedback').respond({});

    scope.message = 'testing';
    scope.send();

    expect(scope.message).toBe('');
    $httpBackend.flush();

    expect(GlobalMsg.get()).toBe('Hemos recibido tu mensaje correctamente');
    expect(GlobalMsg.getClass()).toBe('label-success');
  });

  it('should put the message again in the textarea on fail', function() {
    $httpBackend.expectPOST('/_/feedback').respond(function() {
      return [403, {}];
    });

    scope.message = 'testing';
    scope.send();

    expect(scope.message).toBe('');
    $httpBackend.flush();
    expect(scope.message).toBe('testing');
  });
});

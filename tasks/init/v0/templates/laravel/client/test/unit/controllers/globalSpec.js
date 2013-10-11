'use strict';


describe('Controller: AppCtrl', function() {
  beforeEach(module('controllers.global'));
  beforeEach(module('services.modals'));

  var scope;
  beforeEach(inject(function($controller, $rootScope) {
    scope = $rootScope.$new();
    $controller('AppCtrl', {$scope: scope});
  }));

  it('should notify Analytics', inject(function($rootScope, $location) {
    scope.bodyData = {};

    $rootScope.$broadcast('$routeChangeSuccess');
    window._gaq = [];
    $location.url('/testing');
    $rootScope.$broadcast('$routeChangeSuccess');
    expect(window._gaq.length).toBe(1);
    expect(window._gaq[0].length).toBe(2);
    expect(window._gaq[0][0]).toBe('_trackPageview');
    expect(window._gaq[0][1]).toBe('/testing');
  }));

  it('should scope some services', inject(function(ErrorRegister, GlobalMsg) {
    expect(scope.ErrorRegister).toBe(ErrorRegister);
    expect(scope.GlobalMsg).toBe(GlobalMsg);
  }));

  it('should show a modal if there is an error', inject(function(Modal,
      ErrorRegister, $rootScope) {
    spyOn(Modal, 'errorDialog');
    spyOn(ErrorRegister, 'clean');

    $rootScope.$apply();

    ErrorRegister.set('foo');
    $rootScope.$apply();

    expect(Modal.errorDialog.calls.length).toBe(1);
    expect(ErrorRegister.clean).toHaveBeenCalled();

    ErrorRegister.set(null);
    $rootScope.$apply();

    expect(Modal.errorDialog.calls.length).toBe(1);
  }));
});


describe('Controller: NotFoundCtrl', function() {
  beforeEach(module('controllers.global'));

  beforeEach(inject(function($controller, $rootScope) {
    $controller('NotFoundCtrl', {$scope: $rootScope.$new()});
  }));

  it('should report the error', inject(function($httpBackend) {
    $httpBackend.expectPOST('/_/not-found').respond({});
    $httpBackend.flush();
  }));
});

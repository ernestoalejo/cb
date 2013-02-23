'use strict';


describe('Controller: GlobalMsgCtrl', function() {
  beforeEach(module('app'));

  var scope, GlobalMsg;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');
    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller(GlobalMsgCtrl, {$scope: scope});
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

'use strict';


describe('GlobalMsg tests', function() {
  beforeEach(module('services.global'));

  var GlobalMsg, $timeout;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');
    $timeout = $injector.get('$timeout');
  }));

  it('should save the last alert', function() {
    GlobalMsg.create('one before', 'success');
    expect(GlobalMsg.getLast().msg).toBe('one before');
    expect(GlobalMsg.getLast().type).toBe('success');

    GlobalMsg.create('the last one', 'info');
    expect(GlobalMsg.getLast().msg).toBe('the last one');
    expect(GlobalMsg.getLast().type).toBe('info');
  });
});


describe('ErrorRegister tests', function() {
  beforeEach(module('services.global'));

  var ErrorRegister, $timeout;
  beforeEach(inject(function($injector) {
    ErrorRegister = $injector.get('ErrorRegister');
    $timeout = $injector.get('$timeout');
  }));

  it('should save & clean the error', function() {
    expect(ErrorRegister.isNull()).toBeTruthy();
    ErrorRegister.set('testing');
    expect(ErrorRegister.isNull()).toBeFalsy();
    ErrorRegister.clean();
    expect(ErrorRegister.isNull()).toBeTruthy();
  });
});

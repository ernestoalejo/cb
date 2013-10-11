'use strict';


describe('GlobalMsg tests', function() {
  beforeEach(module('services.global'));
});


describe('ErrorRegister tests', function() {
  beforeEach(module('services.global'));

  it('should save & clean the error', inject(function(ErrorRegister) {
    expect(ErrorRegister.isNull()).toBeTruthy();
    ErrorRegister.set('testing');
    expect(ErrorRegister.isNull()).toBeFalsy();
    ErrorRegister.clean();
    expect(ErrorRegister.isNull()).toBeTruthy();
  }));
});

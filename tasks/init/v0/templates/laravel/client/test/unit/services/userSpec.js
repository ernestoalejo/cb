'use strict';


describe('Service: User', function() {
  beforeEach(module('services.admin.user'));

  beforeEach(module(function($provide) {
    $provide.constant('admin', true);
    $provide.constant('username', 'foo');
    $provide.constant('id', 'bar');
  }));

  it('should save admin info', inject(function(User) {
    expect(User.isAdmin()).toBeTruthy();
  }));

  it('should save username info', inject(function(User) {
    expect(User.getUsername()).toBe('foo');
  }));

  it('should save user id', inject(function(User) {
    expect(User.getId()).toBe('bar');
  }));
});

'use strict';


describe('GlobalMsg tests', function() {
  beforeEach(module('services.global'));

  var GlobalMsg, $timeout;
  beforeEach(inject(function($injector) {
    GlobalMsg = $injector.get('GlobalMsg');
    $timeout = $injector.get('$timeout');
  }));

  it('should save the message', function() {
    GlobalMsg.set('testing');
    expect(GlobalMsg.get()).toBe('testing');
  });

  it('should respect the CAS semantics', function() {
    GlobalMsg.set('foo');
    GlobalMsg.cas('foo', 'bar');
    expect(GlobalMsg.get()).toBe('bar');

    GlobalMsg.cas('foo', 'baz');
    expect(GlobalMsg.get()).toBe('bar');
  });

  it('should hide temp messages', function() {
    GlobalMsg.setTemp('foo');
    $timeout.flush();
    expect(GlobalMsg.get()).toBe('');
  });

  it('should stop old hiding actions on set', function() {
    GlobalMsg.setTemp('foo');
    GlobalMsg.set('bar');

    var called = 0;
    try {
      $timeout.flush();
    } catch (err) {
      called++;
      expect(err.message).toBe('No deferred tasks to be flushed');
    }
    expect(called).toBe(1);
    expect(GlobalMsg.get()).toBe('bar');
  });

  it('should store the class correctly', function() {
    GlobalMsg.set('foo', 'error');
    expect(GlobalMsg.getClass()).toBe('label-error');
    GlobalMsg.set('foo');
    expect(GlobalMsg.getClass()).toBe('label-success');

    GlobalMsg.setTemp('foo', 'error');
    expect(GlobalMsg.getClass()).toBe('label-error');
    GlobalMsg.setTemp('foo');
    expect(GlobalMsg.getClass()).toBe('label-success');
  });
});

describe('Selector tests', function() {
  beforeEach(module('services.global'));

  var Selector;
  beforeEach(inject(function($injector) {
    Selector = $injector.get('Selector');
  }));

  it('should save the navbar and sidebar correctly', function() {
    expect(Selector.getNavbar()).toBe('');
    Selector.setNavbar('example-navbar')
    expect(Selector.getNavbar()).toBe('example-navbar');
  });

  it('should save the dirty flags correctly', function() {
    expect(Selector.isDirty()).toBeFalsy();

    Selector.setDirty();
    expect(Selector.isDirty()).toBeTruthy();

    Selector.setNavbar('test');
    expect(Selector.isDirty()).toBeFalsy();
    expect(Selector.isNavbarDirty()).toBeFalsy();

    Selector.setDirty();
    expect(Selector.isDirty()).toBeTruthy();
    expect(Selector.isNavbarDirty()).toBeTruthy();
  });
});

'use strict';


describe('Service: Modal', function() {
  beforeEach(module('services.modals'));

  it('should create delete confirms', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.deleteConfirm('foo', 'bar').then(callback, callbackError);

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].templateUrl)
        .toBe('/views/services/modal/delete-confirm.html');
    expect($modal.open.calls[0].args[0].controller).toBe('ModalInstanceCtrl');
    expect($modal.open.calls[0].args[0].resolve.title()).toBe('');
    expect($modal.open.calls[0].args[0].resolve.text()).toBe('foo');
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('bar');

    open.resolve();
    $rootScope.$apply();

    expect(callback).toHaveBeenCalled();
    expect(callbackError).not.toHaveBeenCalled();
  }));

  it('should create default delete confirms', inject(function(Modal, $modal,
      $q) {
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.deleteConfirm('foo');

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('Eliminar');
  }));

  it('should reject delete confirms', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.deleteConfirm('foo', 'bar').then(callback, callbackError);

    open.reject();
    $rootScope.$apply();

    expect(callback).not.toHaveBeenCalled();
    expect(callbackError).toHaveBeenCalled();
  }));

  it('should create confirms', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.confirm('foo', 'bar').then(callback, callbackError);

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].templateUrl)
        .toBe('/views/services/modal/confirm.html');
    expect($modal.open.calls[0].args[0].controller).toBe('ModalInstanceCtrl');
    expect($modal.open.calls[0].args[0].resolve.title()).toBe('');
    expect($modal.open.calls[0].args[0].resolve.text()).toBe('foo');
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('bar');

    open.resolve();
    $rootScope.$apply();

    expect(callback).toHaveBeenCalled();
    expect(callbackError).not.toHaveBeenCalled();
  }));

  it('should create default confirms', inject(function(Modal, $modal,
      $q) {
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.confirm('foo');

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('Aceptar');
  }));

  it('should reject confirms', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.confirm('foo', 'bar').then(callback, callbackError);

    open.reject();
    $rootScope.$apply();

    expect(callback).not.toHaveBeenCalled();
    expect(callbackError).toHaveBeenCalled();
  }));

  it('should create alerts', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.alert('foo', 'bar', 'baz').then(callback, callbackError);

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].templateUrl)
        .toBe('/views/services/modal/alert.html');
    expect($modal.open.calls[0].args[0].controller).toBe('ModalInstanceCtrl');
    expect($modal.open.calls[0].args[0].resolve.title()).toBe('foo');
    expect($modal.open.calls[0].args[0].resolve.text()).toBe('bar');
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('baz');

    open.resolve();
    $rootScope.$apply();

    expect(callback).toHaveBeenCalled();
    expect(callbackError).not.toHaveBeenCalled();
  }));

  it('should create default alerts', inject(function(Modal, $modal,
      $q) {
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.alert('foo', 'bar');

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('Aceptar');
  }));

  it('should reject alerts', inject(function(Modal, $modal, $q,
      $rootScope) {
    var callback = jasmine.createSpy();
    var callbackError = jasmine.createSpy();
    var open = $q.defer();
    spyOn($modal, 'open').andReturn({result: open.promise});

    Modal.alert('foo', 'bar').then(callback, callbackError);

    open.reject();
    $rootScope.$apply();

    expect(callback).not.toHaveBeenCalled();
    expect(callbackError).toHaveBeenCalled();
  }));

  it('should create eror dialogs', inject(function(Modal, $modal) {
    spyOn($modal, 'open');

    Modal.errorDialog();

    expect($modal.open).toHaveBeenCalled();
    expect($modal.open.calls[0].args[0].templateUrl)
        .toBe('/views/services/modal/error-dialog.html');
    expect($modal.open.calls[0].args[0].controller).toBe('ModalInstanceCtrl');
    expect($modal.open.calls[0].args[0].resolve.title()).toBe('');
    expect($modal.open.calls[0].args[0].resolve.text()).toBe('');
    expect($modal.open.calls[0].args[0].resolve.label()).toBe('');
  }));
});


describe('Controller: ModalInstanceCtrl', function() {
  beforeEach(module('services.modals'));

  var scope, modalInstance;
  beforeEach(inject(function($controller, $rootScope) {
    modalInstance = {
      close: jasmine.createSpy(),
      dismiss: jasmine.createSpy()
    };

    scope = $rootScope.$new();
    $controller('ModalInstanceCtrl', {
      $scope: scope,
      $modalInstance: modalInstance,
      title: 'foo',
      text: 'bar',
      label: 'baz'
    });
  }));

  it('should scope some vars', function() {
    expect(scope.title).toBe('foo');
    expect(scope.text).toBe('bar');
    expect(scope.label).toBe('baz');
  });

  it('should close on ok', function() {
    scope.ok();
    expect(modalInstance.close).toHaveBeenCalled();
    expect(modalInstance.dismiss).not.toHaveBeenCalled();
  });

  it('should dismiss on cancel', function() {
    scope.cancel();
    expect(modalInstance.close).not.toHaveBeenCalled();
    expect(modalInstance.dismiss).toHaveBeenCalledWith('cancel');
  });
});

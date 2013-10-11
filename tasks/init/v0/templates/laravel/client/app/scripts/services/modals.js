'use strict';


var m = angular.module('services.modals', ['ui.bootstrap.modal']);


m.factory('Modal', function($q, $modal) {
  return {
    deleteConfirm: function(text, removeLabel) {
      removeLabel = removeLabel || 'Eliminar';

      var modalInstance = $modal.open({
        templateUrl: '/views/services/modal/delete-confirm.html',
        controller: 'ModalInstanceCtrl',
        resolve: {
          title: function() { return ''; },
          text: function() { return text; },
          label: function() { return removeLabel; }
        }
      });

      var result = $q.defer();
      modalInstance.result.then(function() {
        result.resolve();
      }, function() {
        result.reject();
      });

      return result.promise;
    },

    confirm: function(text, okLabel) {
      okLabel = okLabel || 'Aceptar';

      var modalInstance = $modal.open({
        templateUrl: '/views/services/modal/confirm.html',
        controller: 'ModalInstanceCtrl',
        resolve: {
          title: function() { return ''; },
          text: function() { return text; },
          label: function() { return okLabel; }
        }
      });

      var result = $q.defer();
      modalInstance.result.then(function() {
        result.resolve();
      }, function() {
        result.reject();
      });

      return result.promise;
    },

    alert: function(title, text, okLabel) {
      okLabel = okLabel || 'Aceptar';

      var modalInstance = $modal.open({
        templateUrl: '/views/services/modal/alert.html',
        controller: 'ModalInstanceCtrl',
        resolve: {
          title: function() { return title; },
          text: function() { return text; },
          label: function() { return okLabel; }
        }
      });

      var result = $q.defer();
      modalInstance.result.then(function() {
        result.resolve();
      }, function() {
        result.reject();
      });

      return result.promise;
    },

    errorDialog: function() {
      $modal.open({
        templateUrl: '/views/services/modal/error-dialog.html',
        controller: 'ModalInstanceCtrl',
        resolve: {
          title: function() { return ''; },
          text: function() { return ''; },
          label: function() { return ''; }
        }
      });
    }
  };
});


m.controller('ModalInstanceCtrl', function($scope, $modalInstance, title, text,
    label) {
  $scope.title = title;
  $scope.text = text;
  $scope.label = label;

  $scope.ok = function() {
    $modalInstance.close();
  };

  $scope.cancel = function() {
    $modalInstance.dismiss('cancel');
  };
});

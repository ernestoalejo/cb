'use strict';


var m = angular.module('directives.match', []);


m.directive('match', function() {
  return {
    require: 'ngModel',
    link: function(scope, elm, attr, ctrl) {
      var widget = elm.controller('form')[attr.match];
      widget.$parsers.push(function(value) {
        ctrl.$setValidity('match', value === ctrl.$viewValue);
        return value;
      });

      function validator(value) {
        if (value == widget.$viewValue && value) {
          ctrl.$setValidity('match', true);
          return value;
        }
        ctrl.$setValidity('match', false);
      }
      ctrl.$formatters.push(validator);
      ctrl.$parsers.push(validator);
    }
  };
});

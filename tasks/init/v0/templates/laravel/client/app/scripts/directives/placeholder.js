'use strict';


var m = angular.module('directives.placeholder', []);


m.directive('placeholder', function() {
  return {
    link: function(_, elm) {
      $(elm[0]).placeholder();
    }
  };
});

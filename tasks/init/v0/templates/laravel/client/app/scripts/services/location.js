'use strict';


var m = angular.module('services.location', []);


/**
 * Helper to modify the location URL (can be mocked out in tests).
 */
m.factory('Location', function() {
  return {
    set: function(url) {
      location.href = url;
    },

    reload: function() {
      location.reload();
    }
  };
});

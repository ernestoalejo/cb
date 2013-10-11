'use strict';


var m = angular.module('services.admin.user', []);


m.factory('User', function(admin, id, name) {
  return {
    isAdmin: function() {
      return admin;
    },

    getUsername: function() {
      return name;
    },

    getId: function() {
      return id;
    }
  };
});

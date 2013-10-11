'use strict';


var m = angular.module('services.admin.user', []);


m.factory('User', function(admin, id, username) {
  return {
    isAdmin: function() {
      return admin;
    },

    getUsername: function() {
      return username;
    },

    getId: function() {
      return id;
    }
  };
});
